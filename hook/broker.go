package hook

import (
	"bytes"
	"context"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/mochi-mqtt/server/v2/system"
	"github.com/weflux/loopin/broker"
	"github.com/weflux/loopin/proxy"
	"log/slog"
)

func NewBroker(
	broker broker.Broker,
	proxyMap *proxy.ProxyMap,
	logger *slog.Logger,
) *Broker {
	b := &Broker{
		HookBase: mqtt.HookBase{},
		broker:   broker,
		logger:   logger,
		proxyMap: proxyMap,
	}

	b.ctx, b.cancelCtx = context.WithCancel(context.Background())

	return b
}

type Broker struct {
	mqtt.HookBase
	broker    broker.Broker
	logger    *slog.Logger
	ctx       context.Context
	cancelCtx context.CancelFunc
	proxyMap  *proxy.ProxyMap
}

// ID returns the ID of the hook.
func (h *Broker) ID() string {
	return "broker-hook"
}

// Provides indicates which methods a hook provides. The default is none - this method
// should be overridden by the embedding hook.
func (h *Broker) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnStarted,
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribe,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribe,
		mqtt.OnUnsubscribed,
		mqtt.OnPublish,
	}, []byte{b})
}

// Init performs any pre-start initializations for the hook, such as connecting to databases
// or opening files.
func (h *Broker) Init(config any) error {
	//return h.broker.Start(h.ctx)
	return nil
}

// SetOpts is called by the server to propagate internal values and generally should
// not be called manually.
func (h *Broker) SetOpts(l *slog.Logger, opts *mqtt.HookOptions) {
	h.Log = l
	h.Opts = opts
}

// Stop is called to gracefully shut down the hook.
func (h *Broker) Stop() error {
	h.cancelCtx()
	return h.broker.Stop(context.Background())
}

// OnStarted is called when the server starts.
func (h *Broker) OnStarted() {
	//h.broker.PublishJoin()
	if err := h.broker.Start(context.Background()); err != nil {
		panic(err)
	}
}

// OnStopped is called when the server stops.
func (h *Broker) OnStopped() {}

// OnSysInfoTick is called when the server publishes system info.
func (h *Broker) OnSysInfoTick(*system.Info) {}

// OnSubscribe is called when a client subscribes to one or more filters.
func (h *Broker) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

// OnSubscribed is called when a client subscribes to one or more filters.
func (h *Broker) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	subs := []*broker.Subscription{}
	for _, sub := range pk.Filters {
		subs = append(subs, &broker.Subscription{
			Filter: sub.Filter,
			Qos:    sub.Qos,
			Client: cl.ID,
		})
	}
	if err := h.broker.Subscribe(subs); err != nil {
		h.logger.Error("on subscribed error", "error", err)
	}
}

// OnSelectSubscribers is called when selecting subscribers to receive a message.
func (h *Broker) OnSelectSubscribers(subs *mqtt.Subscribers, pk packets.Packet) *mqtt.Subscribers {
	return subs
}

// OnUnsubscribe is called when a client unsubscribes from one or more filters.
func (h *Broker) OnUnsubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

// OnUnsubscribed is called when a client unsubscribes from one or more filters.
func (h *Broker) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {}

// OnPublish is called when a client publishes a message.
func (h *Broker) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	if pk.Ignore {
		return pk, nil
	}

	if !cl.Net.Inline {
		if h.proxyMap != nil {
			if _, ok := h.proxyMap.RPCProxies[pk.TopicName]; ok {
				return pk, nil
			}
		}
		if err := h.broker.Publish(&broker.Publication{
			TopicName: pk.TopicName,
			Retain:    pk.FixedHeader.Retain,
			Qos:       pk.FixedHeader.Qos,
			Payload:   pk.Payload,
		}); err != nil {
			return packets.Packet{}, err
		}
		return packets.Packet{}, nil
	}

	return pk, nil
}

// OnPublished is called when a client has published a message to subscribers.
func (h *Broker) OnPublished(cl *mqtt.Client, pk packets.Packet) {}
