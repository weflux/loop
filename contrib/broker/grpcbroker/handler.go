package grpcbroker

import (
	"github.com/weflux/loop"
	"github.com/weflux/loop/cluster/broker"
	"github.com/weflux/loop/cluster/eventhandler"
)

var _ eventhandler.EventHandler = new(Handler)

type Handler struct {
	hub    *loop.Hub
	broker *GrpcBroker
}

func NewHandler(hub *loop.Hub, b *GrpcBroker) *Handler {
	return &Handler{
		hub:    hub,
		broker: b,
	}
}

func (h *Handler) OnPublish(pub *broker.Publication) error {
	return h.hub.Publish(pub.TopicName, pub.Payload, pub.Retain, pub.Qos)
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

func (h *Handler) OnSubscribe(sub *broker.Subscription) error {
	if h.broker.ID() != sub.Broker {
		return h.broker.routing.Subscribe(sub.Filter, sub.Client, sub.Broker)
	}
	return nil
}

func (h *Handler) OnUnsubscribe(sub *broker.Subscription) error {
	if h.broker.ID() != sub.Broker {
		return h.broker.routing.Unsubscribe(sub.Filter, sub.Client, sub.Broker)
	}
	return nil
}

func (h *Handler) OnSubscribeClient(sub *broker.Subscription) error {
	if cl, ok := h.hub.Clients.Get(sub.Client); ok {
		return h.hub.SubscribeClient(cl, sub.Filter, sub.Qos)
	}

	return nil
}

func (h *Handler) OnUnsubscribeClient(sub *broker.Subscription) error {
	if cl, ok := h.hub.Clients.Get(sub.Client); ok {
		return h.hub.UnsubscribeClient(cl, sub.Filter)
	}
	return nil
}
