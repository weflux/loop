package grpcbroker

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/weflux/loop/cluster/broker"
	brokerpb "github.com/weflux/loop/protocol/broker"
	shared "github.com/weflux/loop/protocol/shared"
	"github.com/weflux/loop/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"net/url"
	"sync"
	"time"
)

var _ broker.Broker = new(GrpcBroker)

type Options struct {
	BrokerID       string
	Addr           string
	ServiceName    string
	ServiceVersion string
}

func NewBroker(opts Options, dis registry.Discovery, logger *slog.Logger) *GrpcBroker {
	b := &GrpcBroker{
		dis:             dis,
		mu:              sync.RWMutex{},
		clients:         map[string]brokerpb.BrokerClient{},
		opts:            opts,
		logger:          logger,
		grpcServiceName: fmt.Sprintf("%s.grpc", opts.ServiceName),
		routing:         broker.NewRoutingTable(),
	}

	return b
}

type GrpcBroker struct {
	mu              sync.RWMutex
	clients         map[string]brokerpb.BrokerClient
	dis             registry.Discovery
	opts            Options
	ClusterTopics   *broker.RoutingTable
	ctx             context.Context
	cancelCtx       context.CancelFunc
	logger          *slog.Logger
	routing         *broker.RoutingTable
	grpcServiceName string
}

func (b *GrpcBroker) ID() string {
	return fmt.Sprintf("%s-%s", b.opts.ServiceName, b.opts.BrokerID)
}

func (b *GrpcBroker) Stop(ctx context.Context) error {
	b.cancelCtx()
	return nil
}

func (b *GrpcBroker) SubscribeClient(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *GrpcBroker) UnsubscribeClient(sub *broker.Subscription) error {
	//TODO implement me
	panic("implement me")
}

func (b *GrpcBroker) DisconnectClient(client *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *GrpcBroker) PublishJoin(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *GrpcBroker) PublishLeave(ch string, info *broker.ClientInfo) error {
	//TODO implement me
	panic("implement me")
}

func (b *GrpcBroker) PublishControl(cmd *broker.Command) error {
	//TODO implement me
	panic("implement me")
}

func containsClient(clients map[string]brokerpb.BrokerClient, brokerId string) bool {
	for id, _ := range clients {
		if brokerId == id {
			return true
		}
	}
	return false
}

func containsInstance(ins []*registry.ServiceInstance, id string) bool {
	for _, in := range ins {
		broerId, ok := getBrokerID(in)
		if !ok {
			continue
		}
		if broerId == id {
			return true
		}
	}
	return false
}

func getBrokerID(in *registry.ServiceInstance) (string, bool) {
	md := in.Metadata
	if md == nil {
		return "", false

	}
	broerID, ok := md["broker_id"]
	return broerID, ok

}

func (b *GrpcBroker) addedClients(ins []*registry.ServiceInstance) []*registry.ServiceInstance {
	var news []*registry.ServiceInstance
	for _, in := range ins {
		broerID, ok := getBrokerID(in)
		if !ok {
			continue
		}
		if containsClient(b.clients, broerID) {
			continue
		} else {
			news = append(news, in)
		}
	}

	return news
}

func (b *GrpcBroker) removedClients(ins []*registry.ServiceInstance) []string {

	var removed []string
	for k, _ := range b.clients {
		if !containsInstance(ins, k) {
			removed = append(removed, k)
		}
	}

	return removed
}

func (b *GrpcBroker) newClient(in *registry.ServiceInstance) brokerpb.BrokerClient {
	es := in.Endpoints
	if len(es) == 0 {
		return nil
	}
	e := es[0]
	uri, err := url.Parse(e)
	if err != nil {
		return nil
	}

	// insecure
	cc, err := grpc.Dial(uri.Host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	return brokerpb.NewBrokerClient(cc)
}

func (b *GrpcBroker) watchMemberChange() {
	w, err := b.dis.Watch(b.ctx, b.grpcServiceName)
	if err != nil {
		panic(err)
	}

	for {
		ins, err := w.Next()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			time.Sleep(time.Second)
			continue
		}

		for _, in := range b.addedClients(ins) {
			brokerID, _ := getBrokerID(in)
			b.logger.Info("cluster added member", "broker_id", brokerID, "broker_endpoint", in.Endpoints)
			cli := b.newClient(in)
			if cli != nil {
				b.mu.Lock()
				b.clients[brokerID] = cli
				b.mu.Unlock()
			}
		}

		for _, id := range b.removedClients(ins) {
			b.logger.Info("cluster removed member", "broker_id", id)
			b.mu.Lock()
			delete(b.clients, id)
			b.mu.Unlock()
		}
	}
}

func (b *GrpcBroker) Start(ctx context.Context) error {
	b.ctx, b.cancelCtx = context.WithCancel(ctx)

	b.logger.Info("grpc broker starting", "service_name", b.opts.ServiceName)
	ins, err := b.dis.GetService(ctx, b.grpcServiceName)
	if err != nil {
		b.logger.Warn("grpc broker get services failed", "error", err.Error())
	}

	for _, in := range ins {
		brokerID, _ := getBrokerID(in)
		b.logger.Info("cluster added member", "broker_id", brokerID, "broker_endpoint", in.Endpoints)
		cli := b.newClient(in)
		if cli != nil {
			b.clients[brokerID] = cli
		}
	}

	go b.watchMemberChange()
	return nil
}

func (b *GrpcBroker) broadcast(sendFn func(cli brokerpb.BrokerClient) error) error {
	wg := sync.WaitGroup{}
	var errs error
	b.mu.RLock()
	clients := b.clients
	b.mu.RUnlock()
	for _, _cli := range clients {
		wg.Add(1)
		cli := _cli
		go func() {
			defer wg.Done()
			err := sendFn(cli)
			if err != nil {
				errs = multierror.Append(errs, err)
			}
		}()
	}
	wg.Wait()
	return errs
}

func (b *GrpcBroker) Subscribe(sub *broker.Subscription) error {
	sendFn := func(cli brokerpb.BrokerClient) error {
		ctx, cancelCtx := timeoutCtx(b.ctx)
		defer cancelCtx()
		if _, err := cli.Subscribe(ctx, &brokerpb.SubscribeRequest{
			Id:     uuid.NewString(),
			Topic:  sub.Filter,
			Client: sub.Client,
			Qos:    int32(sub.Qos),
			Broker: b.ID(),
		}); err != nil {
			return err
		}
		return nil
	}
	return b.broadcast(sendFn)
}

func (b *GrpcBroker) Unsubscribe(sub *broker.Subscription) error {
	sendFn := func(cli brokerpb.BrokerClient) error {
		ctx, cancelCtx := timeoutCtx(b.ctx)
		defer cancelCtx()
		if _, err := cli.Unsubscribe(ctx, &brokerpb.UnsubscribeRequest{
			Id:     uuid.NewString(),
			Topic:  sub.Filter,
			Client: sub.Client,
			Qos:    int32(sub.Qos),
			Broker: b.ID(),
		}); err != nil {
			return err
		}
		return nil
	}
	return b.broadcast(sendFn)
}

func timeoutCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, 5*time.Second)
}

func (b *GrpcBroker) Publish(pub *broker.Publication) error {
	sendFn := func(cli brokerpb.BrokerClient) error {
		ctx, cancelCtx := timeoutCtx(b.ctx)
		defer cancelCtx()
		if _, err := cli.Publish(ctx, &brokerpb.PublishRequest{
			//Id:     uuid.NewString(),
			Topic:  pub.TopicName,
			Retain: pub.Retain,
			Qos:    int32(pub.Qos),
			Metadata: &shared.Metadata{
				Topic:   pub.TopicName,
				User:    "",
				Session: "",
			},
			Payload: pub.Payload,
		}); err != nil {
			return err
		}
		return nil
	}
	return b.broadcast(sendFn)
}
