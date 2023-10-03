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
	"github.com/weflux/loop/hook/badgerv4"
	"github.com/weflux/loop/option"
	"github.com/weflux/loop/proxy"
	"log"
	"log/slog"
)

type Node struct {
	*mqtt.Server
	broker           broker.Broker
	logger           *slog.Logger
	serverSideClient *ServerSideClient
}

func (h *Node) Start(ctx context.Context) error {
	h.logger.Info("node server starting")
	if err := h.broker.Start(ctx); err != nil {
		panic(err)
	}
	h.serverSideClient = newServerSideClient(h, h.broker, h.logger)
	return h.Serve()
}

func (h *Node) Stop(_ context.Context) error {
	h.logger.Info("node server stopping")
	return h.Close()
}

func (h *Node) SubscribeClient(cl *mqtt.Client, filter string, qos byte) error {
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

func (h *Node) UnsubscribeClient(cl *mqtt.Client, filter string) error {
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

func NewNode(
	broker broker.Broker,
	opts option.Options,
	slogger *slog.Logger,
	hooks ...mqtt.Hook,
) *Node {
	s := mqtt.New(&mqtt.Options{
		InlineClient: true,
		Logger:       slogger,
	})

	if opts.Store == nil || opts.Store.Badger != nil {
		//bdopt := badger2.DefaultOptions(".badger")
		//bdopt.Mem
		if err := s.AddHook(new(badgerv4.Hook), nil); err != nil {
			panic(err)
		}
	} else if opts.Store != nil && opts.Store.Redis != nil {
		c := opts.Store.Redis
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

	for _, h := range hooks {
		if err := s.AddHook(h, map[string]interface{}{}); err != nil {
			log.Fatal("add mqtt server hook failed", err)
		}
	}

	var proxyMap *proxy.ProxyMap
	if c := opts.Proxy; c != nil {
		proxyMap = proxy.NewProxyMap(c)
		proxyHook := hook.NewProxy(proxyMap, slogger)
		if err := s.AddHook(proxyHook, map[string]interface{}{}); err != nil {
			log.Fatal("add mqtt server hook failed", err)
		}
	}

	var b = broker
	if b == nil {
		panic("Options.Broker must not be nil")
	}

	brokerHook := hook.NewBroker(b, proxyMap, slogger)
	if err := s.AddHook(brokerHook, map[string]interface{}{}); err != nil {
		log.Fatal("add mqtt server hook failed", err)
	}

	c := opts.MQTT
	if c == nil {
		panic("config mqtt must not null")
	}
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

	node := &Node{
		Server: s,
		logger: slogger,
		broker: b,
	}

	return node
}

func (h *Node) ServerSideClient() *ServerSideClient {
	//return newServerSideClient(h, h.broker, h.logger)
	return h.serverSideClient
}
