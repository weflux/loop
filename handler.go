package loop

import "github.com/weflux/loop/broker"

var _ broker.EventHandler = new(Handler)

type Handler struct {
	hub     *Hub
	routing *broker.RoutingTable
	broker  broker.Broker
}

func NewHandler(hub *Hub, rtable *broker.RoutingTable, broker broker.Broker) *Handler {
	return &Handler{
		hub:     hub,
		routing: rtable,
		broker:  broker,
	}
}

func (h *Handler) OnPublish(pub *broker.Publication) error {
	return h.hub.Publish(pub.TopicName, pub.Payload, pub.Retain, byte(pub.Qos))
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
		return h.routing.Subscribe(sub.Filter, sub.Client, sub.Broker)
	}
	return nil
}

func (h *Handler) OnUnsubscribe(sub *broker.Subscription) error {
	if h.broker.ID() != sub.Broker {
		return h.routing.Unsubscribe(sub.Filter, sub.Client, sub.Broker)
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
