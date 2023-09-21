package loop

import (
	"context"
	rv8 "github.com/go-redis/redis/v8"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/storage/redis"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/weflux/loop/broker"
	"github.com/weflux/loop/hook"
	"github.com/weflux/loop/proxy"
	"log"
	"log/slog"
)

type Hub struct {
	*mqtt.Server
	broker broker.Broker
	logger *slog.Logger
}

func (h *Hub) Start(_ context.Context) error {
	h.logger.Info("hub server starting")
	return h.Serve()
}

func (h *Hub) Stop(_ context.Context) error {
	h.logger.Info("hub server stopping")
	return h.Close()
}

func (h *Hub) SubscribeClient(cl *mqtt.Client, filter string, qos byte) error {
	// 需要注意是否正确，待测试
	pk := packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type: packets.Subscribe,
		},
		Filters: []packets.Subscription{
			{
				Filter: filter,
				Qos:    qos,
			},
		},
		PacketID: 1,
	}

	if err := h.Server.InjectPacket(cl, pk); err != nil {
		return err
	}

	return nil
}

func (h *Hub) UnsubscribeClient(cl *mqtt.Client, filter string) error {
	// 需要注意是否正确，待测试
	pk := packets.Packet{
		FixedHeader: packets.FixedHeader{
			Type: packets.Unsubscribe,
		},
		Filters: []packets.Subscription{
			{
				Filter: filter,
			},
		},
		PacketID: 1,
	}

	if err := h.Server.InjectPacket(cl, pk); err != nil {
		return err
	}
	return nil
}

func NewHub(
	conf Options,
	slogger *slog.Logger,
) *Hub {
	s := mqtt.New(&mqtt.Options{
		InlineClient: true,
		Logger:       slogger,
	})

	if c := conf.Redis; c != nil {
		if err := s.AddHook(new(redis.Hook), &redis.Options{
			HPrefix: c.Prefix,
			Options: &rv8.Options{
				Addr:     c.Addr,
				Password: c.Password,
				DB:       c.DB,
			},
		}); err != nil {
			log.Fatal("add mqtt server hook failed", err)
		}
	}

	if conf.Hooks != nil && conf.Hooks.ACL != nil {
		if err := s.AddHook(conf.Hooks.ACL, map[string]interface{}{}); err != nil {
			log.Fatal("add mqtt server hook failed", err)
		}
	}

	if conf.Hooks != nil && conf.Hooks.Authenticate != nil {
		if err := s.AddHook(conf.Hooks.Authenticate, map[string]interface{}{}); err != nil {
			log.Fatal("add mqtt server hook failed", err)
		}
	}

	if c := conf.Proxy; c != nil {
		pm := proxy.NewProxyMap(c)
		ph := hook.NewProxy(pm, slogger)
		if err := s.AddHook(ph, map[string]interface{}{}); err != nil {
			log.Fatal("add mqtt server hook failed", err)
		}
	}

	var b broker.Broker
	if b := conf.Broker; b != nil {
		b = conf.Broker
	} else {
		b = broker.NewMemBroker(slogger)
	}

	bh := hook.NewBroker(b, slogger)
	if err := s.AddHook(bh, map[string]interface{}{}); err != nil {
		log.Fatal("add mqtt server hook failed", err)
	}

	c := conf.MQTT
	if c.TCP == nil && c.WebSocket == nil {
		panic("config mqtt.tcp/mqtt.websocket must not null")
	}
	if c.TCP != nil {
		tcp := listeners.NewTCP("edgehub-tcp", c.TCP.Addr, nil)
		if err := s.AddListener(tcp); err != nil {
			panic(err)
		}
	}

	if c.WebSocket != nil {
		ws := listeners.NewWebsocket("edgehub-ws", c.WebSocket.Addr, nil)
		if err := s.AddListener(ws); err != nil {
			panic(err)
		}
	}

	hub := &Hub{
		Server: s,
		logger: slogger,
		broker: b,
	}

	return hub
}

func (h *Hub) ServerSideClient() *ServerSideClient {
	return newServerSideClient(h, h.broker, h.logger)
}
