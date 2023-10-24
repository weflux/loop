package proxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/weflux/loopify/contenttype"
	proxypb "github.com/weflux/loopify/protocol/proxy"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func NewDispatcher(routes map[string]Handler) *Dispatcher {
	return &Dispatcher{routes: routes}
}

type Result interface {
	MarshalBinary() ([]byte, error)
	MarshalString() (string, error)
}

type JSONResult struct {
	Data any
}

func (r *JSONResult) MarshalBinary() ([]byte, error) {
	return json.Marshal(r.Data)
}

func (r *JSONResult) MarshalString() (string, error) {
	s, err := json.Marshal(r.Data)
	if err != nil {
		return "", nil
	}
	return string(s), nil
}

type ProtoResult struct {
	Data proto.Message
}

func (r *ProtoResult) MarshalBinary() ([]byte, error) {
	return proto.Marshal(r.Data)
}

func (r *ProtoResult) MarshalString() (string, error) {
	s, err := protojson.Marshal(r.Data)
	if err != nil {
		return "", err
	}
	//return base64.StdEncoding.EncodeToString(s), nil
	return string(s), nil
}

type Handler func(ctx context.Context, req *proxypb.RPCRequest) (Result, error)

type Dispatcher struct {
	routes map[string]Handler
}

func (d *Dispatcher) Dispatch(ctx context.Context, req *proxypb.RPCRequest) (*proxypb.RPCReply, error) {
	if h, ok := d.routes[req.Command]; ok {
		ct := contenttype.FromString(req.ContentType)
		result, err := h(ctx, req)
		if err != nil {
			return nil, err
		}

		rep := &proxypb.RPCReply{
			Id:       req.Id,
			Metadata: req.Metadata,
		}
		if err := marshalPayload(ct, result, rep); err != nil {
			return nil, err
		}
		return rep, nil
	}
	return &proxypb.RPCReply{
		Id: req.Id,
		Error: &proxypb.Error{
			Code:    404,
			Message: fmt.Sprintf("handler not found for method %s", req.Command),
		},
	}, nil
}

func marshalPayload(ct contenttype.ContentType, res Result, rep *proxypb.RPCReply) error {
	rep.ContentType = string(ct)
	switch ct {
	case contenttype.JSON:
		data, err := res.MarshalString()
		if err != nil {
			return err
		}
		rep.PayloadText = data
	case contenttype.PROTOBUF:
		data, err := res.MarshalBinary()
		if err != nil {
			return err
		}
		rep.PayloadBytes = data
	default:
		return errors.New("not support content-type")
	}

	return nil
}
