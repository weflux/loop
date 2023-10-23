package eventhandler

import "github.com/weflux/loopin/broker"

// EventHandler can handle messages received from PUB/SUB system.
type EventHandler interface {
	// OnPublish to handle received broker.Publications.
	OnPublish(pub *broker.Publication) error
	// OnClientJoin to handle received Join messages.
	OnClientJoin(ch string, info *broker.ClientInfo) error
	// OnClientLeave to handle received Leave messages.
	OnClientLeave(ch string, info *broker.ClientInfo) error
	// OnControl to handle received control data.
	OnControl(cmd *broker.Command) error

	OnSubscribe(subs []*broker.Subscription) error

	OnUnsubscribe(subs []*broker.Subscription) error

	OnSubscribeClient(sub *broker.Subscription) error

	OnUnsubscribeClient(sub *broker.Subscription) error

	OnDisconnect(client string, broker string) error
}
