package membroker

import (
	"context"
	"github.com/weflux/loopin/broker"
	"log/slog"
)

var _ broker.Broker = new(MemBroker)

type MemBroker struct {
	logger *slog.Logger
	queue  *Queue
}

func (b *MemBroker) Stop(ctx context.Context) error {
	return nil
}

func (b *MemBroker) ID() string {
	return "0"
}

func (b *MemBroker) SubscribeClient(sub *broker.Subscription) error {
	b.queue.SubscribeClient <- sub
	return nil
}

func (b *MemBroker) UnsubscribeClient(sub *broker.Subscription) error {
	b.queue.UnsubscribeClient <- sub
	return nil
}

func (b *MemBroker) DisconnectClient(client *broker.ClientInfo) error {
	b.queue.DisconnectClient <- client
	return nil
}

func (b *MemBroker) PublishJoin(ch string, info *broker.ClientInfo) error {
	return nil
}

func (b *MemBroker) PublishLeave(ch string, info *broker.ClientInfo) error {
	return nil
}

func (b *MemBroker) PublishControl(cmd *broker.Command) error {
	return nil
}

func (b *MemBroker) Subscribe(subs []*broker.Subscription) error {
	//b.queue.Subscribe <-
	for _, sub := range subs {
		b.queue.Subscribe <- sub
	}
	return nil
}

func (b *MemBroker) Unsubscribe(subs []*broker.Subscription) error {
	for _, sub := range subs {
		b.queue.Unsubscribe <- sub
	}
	return nil
}

func (b *MemBroker) Publish(pub *broker.Publication) error {
	b.queue.Publish <- pub
	return nil
}

func NewMemBroker(queue *Queue, logger *slog.Logger) *MemBroker {
	return &MemBroker{
		logger: logger,
		queue:  queue,
	}
}

func (b *MemBroker) Start(ctx context.Context) error {
	return nil
}
