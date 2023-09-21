package loop

import (
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/weflux/loop/broker"
)

type Options struct {
	Redis  *RedisOption  `json:"redis"`
	MQTT   *MQTTOption   `json:"mqtt"`
	Proxy  *ProxyOption  `json:"proxy"`
	Hooks  *HooksOption  `json:"-"`
	Broker broker.Broker `json:"-"`
}

type RedisOption struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}

type RouteOption struct {
	Endpoint    string `json:"endpoint"`
	ContentType string `json:"content_type"`
}

type ProxyOption struct {
	Connect   *RouteOption            `json:"connect"`
	Subscribe *RouteOption            `json:"subscribe"`
	RPC       map[string]*RouteOption `json:"rpc"`
}

type MQTTOption struct {
	TCP       *MQTTTcpOption       `json:"tcp"`
	WebSocket *MQTTWebSocketOption `json:"websocket"`
}

type MQTTTcpOption struct {
	Addr string `json:"addr"`
}

type MQTTWebSocketOption struct {
	Addr string `json:"addr"`
}

type HooksOption struct {
	ACL          mqtt.Hook
	Authenticate mqtt.Hook
}
