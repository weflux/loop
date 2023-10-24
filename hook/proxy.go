package hook

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/weflux/loopify/errcodes"
	loopifyv1 "github.com/weflux/loopify/protocol/loopify/v1"
	proxypb "github.com/weflux/loopify/protocol/proxy"
	sharedpb "github.com/weflux/loopify/protocol/shared"
	"github.com/weflux/loopify/proxy"
	"github.com/weflux/loopify/utils/clientutil"
	"github.com/weflux/loopify/utils/packetutil"
	"github.com/weflux/loopify/utils/topicutil"
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
	topics := []string{}
	for _, f := range pk.Filters {
		topics = append(topics, fmt.Sprintf("%s [%d]", f.Filter, f.Qos))
	}
	h.Log.Debug("on subscribe channel", "client_id", cl.ID, "topics", topics)

	if h.proxyMap.SubscribeProxy != nil {
		// TODO
	} else {
		return h.HookBase.OnSubscribe(cl, pk)
	}
	return pk
}

func (h *Proxy) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	topics := []string{}
	for _, f := range pk.Filters {
		topics = append(topics, fmt.Sprintf("%s [%d]", f.Filter, f.Qos))
	}
	h.Log.Debug("on subscribed channel", "client_id", cl.ID, "topics", topics)

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
	defer func() {
		if r := recover(); r != nil {
			h.logger.Error("on publish error", "error", r)
		}
	}()

	h.Log.Debug("on publish to channel", "client_id", cl.ID, "topic_name", pk.TopicName)

	topic := pk.TopicName
	if p, ok := h.proxyMap.RPCProxies[topic]; ok {
		ct := clientutil.GetContentType(cl)
		msg := &loopifyv1.Message{}

		var pack packets.Packet
		var err error

		if err := ct.Unmarshal(pk.Payload, msg); err != nil {
			return packets.Packet{}, err
		}

		request := msg.GetRequest()
		if request == nil {
			return packets.Packet{}, errors.New("request is null for rpc channel")
		}

		req := &proxypb.RPCRequest{
			Id: msg.Id,
			Metadata: &sharedpb.Metadata{
				Client: cl.ID,
				User:   string(cl.Properties.Username),
			},
			Command:      request.GetCommand(),
			ContentType:  request.ContentType,
			PayloadBytes: request.PayloadBytes,
			PayloadText:  request.PayloadText,
		}

		rep, err := p.ProxyRPC(context.Background(), req)
		if err != nil {
			h.logger.Error("rpc proxy error ", err)
			pack = errorPacket(cl, &proxypb.Error{
				Code:    errcodes.RuntimeError,
				Message: err.(error).Error(),
			}, req)
		} else {
			rep.Id = req.Id
			if rep.Metadata == nil {
				rep.Metadata = req.Metadata
			}
			pack = replyPacket(cl, rep, req)
		}
		if cl.Net.Inline {
			pack.Ignore = false
			return pack, nil
		} else {
			if err := cl.WritePacket(pack); err != nil {
				return emptyPacket(), err
			}
		}

		return emptyPacket(), nil
	}
	return h.HookBase.OnPublish(cl, pk)
}

func errorPacket(cl *mqtt.Client, errRep *proxypb.Error, req *proxypb.RPCRequest) packets.Packet {
	ct := clientutil.GetContentType(cl)
	var e *loopifyv1.Error
	if errRep != nil {
		e = &loopifyv1.Error{
			Code:     errRep.Code,
			Message:  errRep.Message,
			Metadata: map[string]string{},
		}
	}
	id := ""
	if req != nil {
		id = req.Id
	}

	command := ""
	if req != nil {
		command = req.Command
	}
	msg := &loopifyv1.Message{
		Id: id,
		Body: &loopifyv1.Message_Reply{Reply: &loopifyv1.Reply{
			Command: command,
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

func emptyPacket() packets.Packet {
	return packets.Packet{
		TopicName: "$RPC/reply",
		Ignore:    true,
	}
}

func replyPacket(cl *mqtt.Client, rep *proxypb.RPCReply, req *proxypb.RPCRequest) packets.Packet {

	ct := clientutil.GetContentType(cl)
	var e *loopifyv1.Error
	if rep.Error != nil {
		e = &loopifyv1.Error{
			Code:     rep.Error.Code,
			Message:  rep.Error.Message,
			Metadata: map[string]string{},
		}
	}
	msg := &loopifyv1.Message{
		Id:      rep.Id,
		Headers: map[string]string{},
		Body: &loopifyv1.Message_Reply{
			Reply: &loopifyv1.Reply{
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
	if h.proxyMap.ConnectProxy != nil {
		req := packetutil.ToConnectRequest(cl, &pk)
		rep, err := h.proxyMap.ConnectProxy.ProxyConnect(emptyCtx(), req)
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
