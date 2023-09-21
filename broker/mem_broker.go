package broker

import (
	"context"
	"log/slog"
)

var _ Broker = new(MemBroker)

type MemBroker struct {
	logger  *slog.Logger
	handler EventHandler
}

func (b *MemBroker) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) ID() string {
	return "0"
}

func (b *MemBroker) SubscribeClient(sub *Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) UnsubscribeClient(sub *Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) DisconnectClient(client *ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) PublishJoin(ch string, info *ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) PublishLeave(ch string, info *ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) PublishControl(cmd *Command) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) Subscribe(sub *Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) Unsubscribe(sub *Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *MemBroker) Publish(pub *Publication) error {
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
