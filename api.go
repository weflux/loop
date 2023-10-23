package loopin

import (
	"context"
	"errors"
	"github.com/google/uuid"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/weflux/loopin/broker"
	"github.com/weflux/loopin/contenttype"
	apiv1 "github.com/weflux/loopin/protocol/api/v1"
	"github.com/weflux/loopin/utils/topicutil"
	"log/slog"
	"sync"
	"time"
)

type ServerSideClient struct {
	node         *Node
	replyAwaiter map[string]chan *envelope.Reply
	mu           sync.Mutex
	broker       broker.Broker
	logger       *slog.Logger
}

func newServerSideClient(node *Node, broker broker.Broker, logger *slog.Logger) *ServerSideClient {

	ssc := &ServerSideClient{
		replyAwaiter: map[string]chan *envelope.Reply{},
		mu:           sync.Mutex{},
		logger:       logger,
		broker:       broker,
		node:         node,
	}
	if err := ssc.Start(context.Background()); err != nil {
		panic(err)
	}

	return ssc
}

func (c *ServerSideClient) Init(node *Node) {
	c.node = node
}

func (c *ServerSideClient) Start(ctx context.Context) error {
	if c.node == nil {
		return errors.New("api client not initialized")
	}
	if err := c.node.Subscribe(topicutil.ReplyTopic(), 1, func(cl *mqtt.Client, sub packets.Subscription, pk packets.Packet) {
		c.logger.Debug("received reply message", "topic", pk.TopicName)
		msg := &envelope.Message{}
		if err := contenttype.JSON.Unmarshal(pk.Payload, msg); err != nil {
			c.logger.Error("unmarshal message error ", err)
			return
		}

		reply := msg.GetReply()
		if reply == nil {
			c.logger.Warn("message.reply is empty")
			return
		}

		c.logger.Debug("received reply", "command", reply.Command)
		if ch, ok := c.replyAwaiter[msg.Id]; ok {
			select {
			case ch <- reply:
				c.logger.Debug("send reply to awaiter ", msg.Id)
			default:
				c.logger.Debug("send reply to awaiter blocked", "msg_id", msg.Id)
			}
		}
	}); err != nil {
		return err
	}

	return nil
}

func (c *ServerSideClient) Stop(ctx context.Context) error {
	return c.node.Unsubscribe(topicutil.ReplyTopic(), 1)
}

func timeoutCtx(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}

// Survey Survey 调用
func (c *ServerSideClient) Survey(ctx context.Context, req *apiv1.SurveyRequest) (*apiv1.SurveyReply, error) {
	id := uuid.NewString()
	bs, err := contenttype.JSON.Marshal(&envelope.Message{
		Id:      id,
		Headers: make(map[string]string),
		Body: &envelope.Message_Request{
			Request: &envelope.Request{
				Command:      req.Command,
				ContentType:  req.ContentType,
				PayloadBytes: req.PayloadBytes,
				PayloadText:  req.PayloadText,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	d := time.Duration(req.Timeout) * time.Millisecond
	if d == 0 {
		d = time.Millisecond * 10_000
	}

	cctx, cancelCtx := timeoutCtx(ctx, d)
	defer cancelCtx()
	// 将 pub/sub 模式转化为 request/reply 模式
	c.mu.Lock()
	replyCh := make(chan *envelope.Reply, req.ExpectReplies)
	c.replyAwaiter[id] = replyCh
	c.mu.Unlock()
	if err := c.node.Server.Publish(req.Topic, bs, false, 0); err != nil {
		return nil, err
	}
	defer func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.replyAwaiter, id)
	}()

	rep := &apiv1.SurveyReply{
		Id:      req.Id,
		Command: req.Command,
		Results: []*apiv1.SurveyReply_Result{},
	}
	for {
		select {
		case <-cctx.Done():
			if len(rep.Results) == 0 {
				rep.Error = &apiv1.Error{
					Code:    502,
					Message: "rpc timeout",
				}
			}
			return rep, nil
		case reply := <-replyCh:
			c.logger.Debug("awaiter channel received reply", "command", reply.Command)
			var e *apiv1.Error
			if reply.Error != nil {
				e = &apiv1.Error{
					Code:    reply.Error.Code,
					Message: reply.Error.Message,
				}
			}

			rep.Results = append(rep.Results, &apiv1.SurveyReply_Result{
				Error:        e,
				Metadata:     reply.Metadata,
				ContentType:  reply.ContentType,
				PayloadText:  reply.PayloadText,
				PayloadBytes: reply.PayloadBytes,
			})
			if req.ExpectReplies == 0 {
				continue
			}
			if int(req.ExpectReplies) <= len(rep.Results) {
				return rep, nil
			}

		}
	}
}

// Publish 发布消息
func (c *ServerSideClient) Publish(ctx context.Context, pub *apiv1.PublishRequest) (*apiv1.PublishReply, error) {
	if c.node == nil {
		return nil, errors.New("api client not initialized")
	}

	//encoder := getEncoder(pub.ContentType)
	msg := &envelope.Message{
		Id:      uuid.NewString(),
		Headers: map[string]string{},
		Body: &envelope.Message_Publication{
			Publication: &envelope.Publication{
				Type:         pub.Type,
				ContentType:  pub.ContentType,
				PayloadText:  pub.PayloadText,
				PayloadBytes: pub.PayloadBytes,
			},
		},
	}

	if data, err := contenttype.JSON.Marshal(msg); err != nil {
		return &apiv1.PublishReply{Error: &apiv1.Error{
			Code:    500,
			Message: err.Error(),
		}}, nil
	} else {
		if err := c.node.Publish(pub.Topic, data, pub.Retain, byte(pub.Qos)); err != nil {
			return &apiv1.PublishReply{Error: &apiv1.Error{
				Code:    500,
				Message: err.Error(),
			}}, nil
		}
	}
	return &apiv1.PublishReply{}, nil
}

// Subscribe 服务端订阅
func (c *ServerSideClient) Subscribe(ctx context.Context, req *apiv1.SubscribeRequest) (*apiv1.SubscribeReply, error) {
	if c.node == nil {
		return nil, errors.New("api client not initialized")
	}
	if err := c.broker.Subscribe([]*broker.Subscription{
		{
			Filter: req.Filter,
			Qos:    byte(req.Qos),
			Client: req.Client,
			Broker: req.Broker,
		},
	}); err != nil {
		return nil, err
	}

	return &apiv1.SubscribeReply{}, nil
}

// Unsubscribe 服务端取消订阅
func (c *ServerSideClient) Unsubscribe(ctx context.Context, req *apiv1.UnsubscribeRequest) (*apiv1.UnsubscribeReply, error) {
	if c.node == nil {
		return nil, errors.New("api client not initialized")
	}
	if err := c.broker.Unsubscribe([]*broker.Subscription{
		{
			Filter: req.Filter,
			Client: req.Client,
			Broker: req.Broker,
		},
	}); err != nil {
		return nil, err
	}

	return &apiv1.UnsubscribeReply{}, nil
}

// Disconnect 断开客户端连接
func (c *ServerSideClient) Disconnect(ctx context.Context, req *apiv1.DisconnectRequest) (*apiv1.DisconnectReply, error) {

	if c.node == nil {
		return nil, errors.New("api client not initialized")
	}

	cl, ok := c.node.Server.Clients.Get(req.Client)
	if !ok {
		return nil, nil
	}
	if err := c.node.Server.DisconnectClient(cl, packets.CodeDisconnect); err != nil {
		return nil, err
	}

	return &apiv1.DisconnectReply{}, nil
}
