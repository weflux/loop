package hook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/samber/lo"
	"github.com/weflux/loop/errcodes"
	"github.com/weflux/loop/protocol/message/v1"
	proxypb "github.com/weflux/loop/protocol/proxy"
	shared "github.com/weflux/loop/protocol/shared"
	"github.com/weflux/loop/proxy"
	"github.com/weflux/loop/utils/clientutil"
	"github.com/weflux/loop/utils/packetutil"
	"github.com/weflux/loop/utils/topicutil"
	"log/slog"
)

var _ mqtt.Hook = new(Proxy)

type Proxy struct {
	mqtt.HookBase
	proxyMap *proxy.ProxyMap
	logger   *slog.Logger
}

func NewProxy(proxyMap *proxy.ProxyMap, logger *slog.Logger) *Proxy {
	return &Proxy{
		HookBase: mqtt.HookBase{},
		proxyMap: proxyMap,
		logger:   logger,
	}
}

func (h *Proxy) ID() string {
	return "proxy-hook"
}

func (h *Proxy) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribe,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribe,
		mqtt.OnUnsubscribed,
		mqtt.OnPublish,
		mqtt.OnConnectAuthenticate,
	}, []byte{b})
}

func (h *Proxy) Init(config any) error {
	return nil
}

func (h *Proxy) Stop() error {
	return nil
}

func (h *Proxy) SetOpts(l *slog.Logger, o *mqtt.HookOptions) {
	h.Log = l
}

func emptyCtx() context.Context {
	return context.Background()
}

func (h *Proxy) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	h.Log.Debug("on connect", "client_id", cl.ID)
	return h.HookBase.OnConnect(cl, pk)
}

func (h *Proxy) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Debug("on disconnect", "client_id", cl.ID)
	h.HookBase.OnDisconnect(cl, err, expire)
}

func (h *Proxy) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	h.Log.Debug("on subscribe channel", "client_id", cl.ID, "topic_name", pk.TopicName)
	if h.proxyMap.SubscribeProxy != nil {
		// TODO
	} else {
		return h.HookBase.OnSubscribe(cl, pk)
	}
	return packets.Packet{}
}

func (h *Proxy) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Debug("on subscribed channel", "client_id", cl.ID, "topic_name", pk.TopicName)
	h.HookBase.OnSubscribed(cl, pk, reasonCodes)
}

func (h *Proxy) OnSelectSubscribers(subs *mqtt.Subscribers, pk packets.Packet) *mqtt.Subscribers {
	return h.HookBase.OnSelectSubscribers(subs, pk)
}

func (h *Proxy) OnUnsubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	h.Log.Debug("on unsubscribe channel", "client_id", cl.ID, "topic_name", pk.TopicName)
	return h.HookBase.OnUnsubscribe(cl, pk)
}

func (h *Proxy) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Debug("on unsubscribed channel", "client_id", cl.ID, "topic_name", pk.TopicName)
	h.HookBase.OnUnsubscribed(cl, pk)
}

func (h *Proxy) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	h.Log.Debug("on publish to channel", "client_id", cl.ID, "topic_name", pk.TopicName)

	topic := pk.TopicName
	if p, ok := h.proxyMap.RPCProxies[topic]; ok {
		ct := clientutil.GetContentType(cl)
		msg := &message.Message{}
		if err := ct.Unmarshal(pk.Payload, msg); err != nil {
			return packets.Packet{}, err
		}

		request := msg.GetRequest()
		if request == nil {
			return packets.Packet{}, errors.New("request is null for rpc channel")
		}

		req := &proxypb.RPCRequest{
			Id: msg.Id,
			Metadata: &shared.Metadata{
				Client: cl.ID,
				User:   string(cl.Properties.Username),
			},
			Command:      request.GetCommand(),
			ContentType:  request.ContentType,
			PayloadBytes: request.PayloadBytes,
			PayloadText:  request.PayloadText,
		}

		var pack packets.Packet
		err, ok := lo.TryWithErrorValue(func() error {
			var err error
			rep, err := p.ProxyRPC(context.Background(), req)
			if err != nil {
				h.logger.Error("rpc proxy error", err)
				return err
			}
			rep.Id = req.Id
			if rep.Metadata == nil {
				rep.Metadata = req.Metadata
			}
			pack = replyPacket(cl, rep, req)
			return nil
		})
		if !ok && err != nil {
			h.logger.Error("rpc proxy error ", err)
			pack = errorPacket(cl, &proxypb.Error{
				Code:    errcodes.RuntimeError,
				Message: err.(error).Error(),
			}, req)
		}

		if cl.ID == "inline" {
			pack.Ignore = false
			return pack, nil
		} else {
			if err := cl.WritePacket(pack); err != nil {
				return pack, err
			}
		}

		return pack, nil
	}
	return h.HookBase.OnPublish(cl, pk)
}

func errorPacket(cl *mqtt.Client, errRep *proxypb.Error, req *proxypb.RPCRequest) packets.Packet {
	ct := clientutil.GetContentType(cl)
	var e *message.Error
	if errRep != nil {
		e = &message.Error{
			Code:    errRep.Code,
			Message: errRep.Message,
			Extras:  map[string]string{},
		}
	}
	msg := &message.Message{
		Id: req.Id,
		Body: &message.Message_Reply{Reply: &message.Reply{
			Command: req.Command,
			Error:   e,
		}},
	}
	data, err := ct.Marshal(msg)
	if err != nil {
		return packets.Packet{}
	}
	pk := packets.Packet{
		FixedHeader: packets.FixedHeader{Type: packets.Publish},
		TopicName:   topicutil.ReplyTopic(),
		Ignore:      true,
		Payload:     data,
	}
	return pk
}

func replyPacket(cl *mqtt.Client, rep *proxypb.RPCReply, req *proxypb.RPCRequest) packets.Packet {

	ct := clientutil.GetContentType(cl)
	var e *message.Error
	if rep.Error != nil {
		e = &message.Error{
			Code:    rep.Error.Code,
			Message: rep.Error.Message,
			Extras:  map[string]string{},
		}
	}
	msg := &message.Message{
		Id:      rep.Id,
		Headers: map[string]string{},
		Body: &message.Message_Reply{
			Reply: &message.Reply{
				Error:        e,
				Command:      req.Command,
				ContentType:  rep.ContentType,
				PayloadText:  rep.PayloadText,
				PayloadBytes: rep.PayloadBytes,
			},
		},
	}
	data, err := ct.Marshal(msg)
	if err != nil {
		return packets.Packet{}
	}
	pk := packets.Packet{
		FixedHeader: packets.FixedHeader{Type: packets.Publish},
		TopicName:   topicutil.ReplyTopic(),
		Ignore:      true,
		Payload:     data,
	}
	return pk

}

func (h *Proxy) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Debug("on published to channel", "client_id", cl.ID, "topic_name", pk.TopicName)
	h.HookBase.OnPublished(cl, pk)
}

func (h *Proxy) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	h.Log.Debug(fmt.Sprintf("[id=%s] try connect authenticate: %s %s", cl.ID, string(pk.Connect.Username), string(pk.Connect.Password)))
	if h.proxyMap.AuthenticateProxy != nil {
		req := packetutil.ToConnectRequest(cl, &pk)
		rep, err := h.proxyMap.AuthenticateProxy.ProxyAuthenticate(emptyCtx(), req)
		if err != nil {
			h.Log.Error("authenticate failed", "error", err)
			return false
		}
		if rep.Error == nil {
			return true
		}

		h.Log.Error("msg", "authenticate failed", "errno", rep.Error.Code, "errmsg", rep.Error.Message)
		return false
	} else {
		return h.HookBase.OnConnectAuthenticate(cl, pk)
	}
}
