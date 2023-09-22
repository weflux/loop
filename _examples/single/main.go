package main

import (
	"context"
	"github.com/weflux/loop"
	"github.com/weflux/loop/hook"
	"github.com/weflux/loop/membroker"
	"github.com/weflux/loop/option"
	"go.mrchanchal.com/zaphandler"
	"go.uber.org/zap"
	"log/slog"
)

func main() {
	opts := option.Options{
		Store: &option.StoreOption{
			Mem: &option.MemOption{},
		},
		MQTT: &option.MQTTOption{
			TCP:       &option.MQTTTcpOption{Addr: "0.0.0.0:31000"},
			WebSocket: &option.MQTTWebSocketOption{Addr: "0.0.0.0:32000"},
		},
		Proxy: &option.ProxyOption{
			Connect: &option.RouteOption{
				Endpoint:    "http://127.0.0.1:17000/proxy/connect?format=json",
				ContentType: "json",
			},
			Subscribe: &option.RouteOption{
				Endpoint:    "http://127.0.0.1:17000/proxy/subscribe?format=json",
				ContentType: "json",
			},
			RPC: map[string]*option.RouteOption{
				"$RPC/http_json": {
					Endpoint:    "http://127.0.0.1:17000/proxy/rpc?format=json",
					ContentType: "json",
				},
			},
		},
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	zlogger, _ := zap.NewDevelopment()
	slogger := slog.New(zaphandler.New(zlogger))
	queue := membroker.NewQueue()
	broker := membroker.NewMemBroker(queue, slogger)
	node := loop.NewNode(broker, opts, slogger, hook.NewACL())
	_ = membroker.NewMemHandler(node, queue, slogger)
	_ = node.Start(ctx)
	<-ctx.Done()
}
