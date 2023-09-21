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
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) UnsubscribeClient(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) DisconnectClient(client *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) Unsubscribe(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) Publish(pub *broker.Publication) error {
	//TODO implement me
	panic("implement me")
}

func NewMemBroker(logger *slog.Logger) *MemBroker {
	return &MemBroker{
		logger: logger,
	}
}

func (b *MemBroker) Start(ctx context.Context) error {
	//b.handler = handlers
	return nil
}

type MemHandler struct {
}

func (m *MemHandler) OnPublish(pub *broker.Publication) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemHandler) OnClientJoin(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemHandler) OnClientLeave(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemHandler) OnControl(cmd *broker.Command) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemHandler) OnSubscribe(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemHandler) OnUnsubscribe(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemHandler) OnSubscribeClient(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (m *MemHandler) OnUnsubscribeClient(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

//
//func (b *MemBroker) History(ch string, opts HistoryOptions) ([]*brokerpb.Publication, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (b *MemBroker) RemoveHistory(ch string) error {
//	//TODO implement me
//	panic("implement me")
//}
