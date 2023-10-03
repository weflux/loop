package grpcbroker

import (
	mqtt "github.com/mochi-mqtt/server/v2"
	matching "github.com/weflux/fast-topic-matching"
	"sync"
)

type RoutingTable struct {
	matching.Matcher
	subscriptions map[string]map[string]*matching.Subscription
	mu            sync.RWMutex
}

type LocalClient struct {
	Client *mqtt.Client
}

type RemoteClient struct {
	Client string
	Broker string
}

func NewRoutingTable() *RoutingTable {
	return &RoutingTable{
		Matcher:       matching.NewCSTrieMatcher(),
		subscriptions: make(map[string]map[string]*matching.Subscription),
		mu:            sync.RWMutex{},
	}
}

func (t *RoutingTable) Unsubscribe(topic string, client string, broker string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if s, ok := t.subscriptions[client]; ok {
		if sub, ok := s[topic]; ok {
			t.Matcher.Unsubscribe(sub)
			delete(t.subscriptions[client], topic)
		}
		if len(t.subscriptions) == 0 {
			delete(t.subscriptions, client)
		}
	}

	return nil
}

func (t *RoutingTable) Subscribe(topic string, client string, broker string) error {
	sub, err := t.Matcher.Subscribe(topic, &RemoteClient{
		Client: client,
		Broker: broker,
	})
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if s, ok := t.subscriptions[client]; ok {
		s[topic] = sub
	} else {
		s = map[string]*matching.Subscription{}
		s[topic] = sub
		t.subscriptions[client] = s
	}

	return nil
}
