package membroker

import (
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/weflux/loop"
	"github.com/weflux/loop/broker"
	"github.com/weflux/loop/broker/eventhandler"
	"log/slog"
)

var _ eventhandler.EventHandler = new(MemHandler)

func NewMemHandler(node *loop.Node, queue *Queue, slogger *slog.Logger) *MemHandler {
	h := &MemHandler{node: node, queue: queue, logger: slogger}
	go h.eventLoop()
	return h
}

type MemHandler struct {
	node   *loop.Node
	logger *slog.Logger
	queue  *Queue
}

func (h *MemHandler) eventLoop() {
	for {
		select {
		case pub := <-h.queue.Publish:
			if err := h.OnPublish(pub); err != nil {
				h.logger.Error("on publish error", "error", err)
			}
		case sub := <-h.queue.Subscribe:
			if err := h.OnSubscribe([]*broker.Subscription{sub}); err != nil {
				h.logger.Error("on subscribe error", "error", err)
			}
		}
	}
}

func (h *MemHandler) OnDisconnect(client string, broker string) error {
	if cl, ok := h.node.Clients.Get(client); ok {
		return h.node.DisconnectClient(cl, packets.CodeDisconnect)
	}
	return nil
}

func (h *MemHandler) OnPublish(pub *broker.Publication) error {
	return h.node.Publish(pub.TopicName, pub.Payload, pub.Retain, pub.Qos)
}

func (h *MemHandler) OnClientJoin(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (h *MemHandler) OnClientLeave(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (h *MemHandler) OnControl(cmd *broker.Command) error {
	//TODO implement me
	panic("implement me")
}

func (h *MemHandler) OnSubscribe(subs []*broker.Subscription) error {
	for _, sub := range subs {
		h.logger.Info("on subscribe", "topic", sub.Filter, "client", sub.Client)
	}
	return nil
}

func (h *MemHandler) OnUnsubscribe(subs []*broker.Subscription) error {
	for _, sub := range subs {
		h.logger.Info("on unsubscribe", "topic", sub.Filter, "client", sub.Client)
	}
	return nil
}

func (h *MemHandler) OnSubscribeClient(sub *broker.Subscription) error {
	if cl, ok := h.node.Clients.Get(sub.Client); ok {
		return h.node.SubscribeClient(cl, sub.Filter, sub.Qos)
	}

	return nil
}

func (h *MemHandler) OnUnsubscribeClient(sub *broker.Subscription) error {
	if cl, ok := h.node.Clients.Get(sub.Client); ok {
		return h.node.UnsubscribeClient(cl, sub.Filter)
	}

	return nil
}
