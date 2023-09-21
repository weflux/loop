package membroker

import (
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/weflux/loop"
	"github.com/weflux/loop/cluster/broker"
	"github.com/weflux/loop/cluster/eventhandler"
	"log/slog"
)

var _ eventhandler.EventHandler = new(MemHandler)

func NewMemHandler(node *loop.Node, slogger *slog.Logger) *MemHandler {
	return &MemHandler{node: node, logger: slogger}
}

type MemHandler struct {
	node   *loop.Node
	logger *slog.Logger
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

func (h *MemHandler) OnSubscribe(sub *broker.Subscription) error {
	h.logger.Info("on subscribe", "topic", sub.Filter, "client", sub.Client)
	return nil
}

func (h *MemHandler) OnUnsubscribe(sub *broker.Subscription) error {
	h.logger.Info("on unsubscribe", "topic", sub.Filter, "client", sub.Client)
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
