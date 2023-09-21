package broker

import (
	"context"
)

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
