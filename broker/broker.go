package broker

import (
	"context"
)

// EventHandler can handle messages received from PUB/SUB system.
type EventHandler interface {
	// OnPublish to handle received Publications.
	OnPublish(pub *Publication) error
	// OnClientJoin to handle received Join messages.
	OnClientJoin(ch string, info *ClientInfo) error
	// OnClientLeave to handle received Leave messages.
	OnClientLeave(ch string, info *ClientInfo) error
	// OnControl to handle received control data.
	OnControl(cmd *Command) error

	OnSubscribe(sub *Subscription) error

	OnUnsubscribe(sub *Subscription) error

	OnSubscribeClient(sub *Subscription) error

	OnUnsubscribeClient(sub *Subscription) error
}

type Broker interface {
	Start(ctx context.Context) error

	Stop(ctx context.Context) error

	ID() string

	SubscribeClient(sub *Subscription) error

	UnsubscribeClient(sub *Subscription) error

	DisconnectClient(client *ClientInfo) error

	// Subscribe 广播订阅消息
	Subscribe(sub *Subscription) error

	// Unsubscribe 广播取消订阅消息
	Unsubscribe(sub *Subscription) error

	Publish(pub *Publication) error

	PublishJoin(ch string, info *ClientInfo) error

	PublishLeave(ch string, info *ClientInfo) error

	PublishControl(cmd *Command) error

	//History(ch string, opts HistoryOptions) ([]*brokerpb.Publication, error)

	//RemoveHistory(ch string) error
}
