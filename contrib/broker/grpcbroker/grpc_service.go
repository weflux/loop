package grpcbroker

import (
	"context"
	"github.com/weflux/loop/cluster/broker"
	"github.com/weflux/loop/cluster/eventhandler"
	brokerpb "github.com/weflux/loop/protocol/broker"
)

func NewBrokerService(handler eventhandler.EventHandler) *BrokerService {
	return &BrokerService{
		handler: handler,
	}
}

type BrokerService struct {
	brokerpb.UnimplementedBrokerServer
	handler eventhandler.EventHandler
}

func (s *BrokerService) Publish(ctx context.Context, req *brokerpb.PublishRequest) (*brokerpb.PublishReply, error) {
	if err := s.handler.OnPublish(&broker.Publication{
		TopicName: req.Topic,
		Retain:    req.Retain,
		Qos:       byte(req.Qos),
		Payload:   req.Payload,
	}); err != nil {
		return nil, err
	}

	return &brokerpb.PublishReply{}, nil

}

func (s *BrokerService) Subscribe(ctx context.Context, req *brokerpb.SubscribeRequest) (*brokerpb.SubscribeReply, error) {
	if err := s.handler.OnSubscribe(&broker.Subscription{
		Filter: req.Topic,
		Qos:    byte(req.Qos),
		Client: req.Client,
	}); err != nil {
		return nil, err
	}

	return &brokerpb.SubscribeReply{}, nil
}

func (s *BrokerService) Unsubscribe(ctx context.Context, req *brokerpb.UnsubscribeRequest) (*brokerpb.UnsubscribeReply, error) {

	if err := s.handler.OnUnsubscribe(&broker.Subscription{
		Filter: req.Topic,
		Client: req.Client,
	}); err != nil {
		return nil, err
	}

	return &brokerpb.UnsubscribeReply{}, nil

}
func (s *BrokerService) SubscribeClient(ctx context.Context, sub *brokerpb.SubscribeRequest) (*brokerpb.SubscribeReply, error) {
	if err := s.handler.OnSubscribeClient(&broker.Subscription{
		Filter: sub.Topic,
		Qos:    byte(sub.Qos),
		Client: sub.Client,
		Broker: sub.Broker,
	}); err != nil {
		return nil, err
	}
	return &brokerpb.SubscribeReply{}, nil
}

func (s *BrokerService) UnsubscribeClient(ctx context.Context, sub *brokerpb.UnsubscribeRequest) (*brokerpb.UnsubscribeReply, error) {
	if err := s.handler.OnUnsubscribeClient(&broker.Subscription{
		Filter: sub.Topic,
		Qos:    byte(sub.Qos),
		Client: sub.Client,
		Broker: sub.Broker,
	}); err != nil {
		return nil, err
	}
	return &brokerpb.UnsubscribeReply{}, nil
}
func (s *BrokerService) DisconnectClient(context.Context, *brokerpb.DisconnectRequest) (*brokerpb.DisconnectReply, error) {
	//if err := s.handler.OnSubscribe()
	return &brokerpb.DisconnectReply{}, nil
}
