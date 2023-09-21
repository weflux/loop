package broker

type Publication struct {
	TopicName string
	Payload   []byte
	Retain    bool
	Qos       byte
}

type Subscription struct {
	Filter string
	Qos    byte
	Client string
	Broker string
}

type ClientInfo struct {
	ClientID string
	Session  string
	User     string
	BrokerID string
}

type Command struct {
}

type PublishOptions struct {
}

type HistoryOptions struct {
}
