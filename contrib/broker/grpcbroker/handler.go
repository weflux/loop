package grpcbroker

import (
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/weflux/loopify"
	"github.com/weflux/loopify/broker"
	"github.com/weflux/loopify/broker/eventhandler"
)

var _ eventhandler.EventHandler = new(Handler)

type Handler struct {
	node   *loopify.Node
	broker *GrpcBroker
}

func (h *Handler) OnDisconnect(client string, broker string) error {
	if cl, ok := h.node.Clients.Get(client); ok {
		return h.node.DisconnectClient(cl, packets.CodeDisconnect)
	}

	return nil
}

func NewHandler(node *loopify.Node, b *GrpcBroker) *Handler {
	return &Handler{
		node:   node,
		broker: b,
	}
}

func (h *Handler) OnPublish(pub *broker.Publication) error {
	return h.node.Publish(pub.TopicName, pub.Payload, pub.Retain, pub.Qos)
}

func (h *Handler) OnClientJoin(ch string, info *broker.ClientInfo) error {
	return nil
}

func (h *Handler) OnClientLeave(ch string, info *broker.ClientInfo) error {
	return nil
}

func (h *Handler) OnControl(cmd *broker.Command) error {
	return nil
}

func (h *Handler) OnSubscribe(subs []*broker.Subscription) error {
	for _, sub := range subs {
		if h.broker.ID() != sub.Broker {
			return h.broker.routing.Subscribe(sub.Filter, sub.Client, sub.Broker)
		}
	}
	return nil
}

func (h *Handler) OnUnsubscribe(subs []*broker.Subscription) error {
	for _, sub := range subs {
		if h.broker.ID() != sub.Broker {
			return h.broker.routing.Unsubscribe(sub.Filter, sub.Client, sub.Broker)
		}
	}
	return nil
}

func (h *Handler) OnSubscribeClient(sub *broker.Subscription) error {
	if cl, ok := h.node.Clients.Get(sub.Client); ok {
		return h.node.SubscribeClient(cl, sub.Filter, sub.Qos)
	}

	return nil
}

func (h *Handler) OnUnsubscribeClient(sub *broker.Subscription) error {
	if cl, ok := h.node.Clients.Get(sub.Client); ok {
		return h.node.UnsubscribeClient(cl, sub.Filter)
	}
	return nil
}
