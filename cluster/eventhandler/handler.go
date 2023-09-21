package eventhandler

import "github.com/weflux/loop/cluster/broker"

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

	OnSubscribe(sub *broker.Subscription) error

	OnUnsubscribe(sub *broker.Subscription) error

	OnSubscribeClient(sub *broker.Subscription) error

	OnUnsubscribeClient(sub *broker.Subscription) error
}
