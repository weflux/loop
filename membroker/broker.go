package membroker

import (
	"context"
	"github.com/weflux/loop/cluster/broker"
	"github.com/weflux/loop/cluster/eventhandler"
	"log/slog"
)

var _ broker.Broker = new(MemBroker)

type MemBroker struct {
	logger  *slog.Logger
	handler eventhandler.EventHandler
}

func (b *MemBroker) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) ID() string {
	return "0"
}

func (b *MemBroker) SubscribeClient(sub *broker.Subscription) error {
	return b.handler.OnSubscribeClient(sub)
}

func (b *MemBroker) UnsubscribeClient(sub *broker.Subscription) error {
	return b.handler.OnUnsubscribeClient(sub)
}

func (b *MemBroker) DisconnectClient(client *broker.ClientInfo) error {
	//return b.handler.OnD(sub)
	return nil
}

func (b *MemBroker) PublishJoin(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) PublishLeave(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) PublishControl(cmd *broker.Command) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) Subscribe(sub *broker.Subscription) error {
	return b.handler.OnSubscribe(sub)
}

func (b *MemBroker) Unsubscribe(sub *broker.Subscription) error {
	return b.handler.OnUnsubscribe(sub)
}

func (b *MemBroker) Publish(pub *broker.Publication) error {
	return b.handler.OnPublish(pub)
}

func NewMemBroker(logger *slog.Logger) *MemBroker {
	return &MemBroker{
		logger: logger,
	}
}

func (b *MemBroker) SetHandler(handler eventhandler.EventHandler) {
	b.handler = handler
}

func (b *MemBroker) Start(ctx context.Context) error {
	//b.handler = handlers
	return nil
}
