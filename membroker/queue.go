package membroker

import "github.com/weflux/loopin/broker"

type Queue struct {
	Publish           chan *broker.Publication
	Subscribe         chan *broker.Subscription
	SubscribeClient   chan *broker.Subscription
	Unsubscribe       chan *broker.Subscription
	UnsubscribeClient chan *broker.Subscription
	DisconnectClient  chan *broker.ClientInfo
}

func NewQueue() *Queue {
	return &Queue{
		Publish:     make(chan *broker.Publication),
		Subscribe:   make(chan *broker.Subscription),
		Unsubscribe: make(chan *broker.Subscription),
	}
}
