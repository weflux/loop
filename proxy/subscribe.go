package proxy

import (
	"context"
	"github.com/weflux/loopify/contenttype"
	"github.com/weflux/loopify/option"
	proxypb "github.com/weflux/loopify/protocol/proxy"
	"net/url"
)

type SubscribeProxy interface {
	ProxySubscribe(ctx context.Context, req *proxypb.SubscribeRequest) (*proxypb.SubscribeReply, error)
}

func NewSubscribeProxy(c *option.RouteOption) SubscribeProxy {

	endpoint := c.Endpoint
	uri, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	if uri.Scheme == "http" || uri.Scheme == "https" {
		return &httpSubscribeProxy{
			HttpBase: HttpBase{URL: endpoint, ContentType: contenttype.ContentType(c.ContentType)},
		}
	}

	return &grpcSubscribeProxy{
		GRPCBase: GRPCBase{Addr: c.Endpoint},
	}
}

type httpSubscribeProxy struct {
	HttpBase
}

func (h *httpSubscribeProxy) ProxySubscribe(ctx context.Context, req *proxypb.SubscribeRequest) (*proxypb.SubscribeReply, error) {
	rep := &proxypb.SubscribeReply{}
	if err := h.ProxyHTTP(req, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type grpcSubscribeProxy struct {
	GRPCBase
}

func (p *grpcSubscribeProxy) ProxySubscribe(ctx context.Context, req *proxypb.SubscribeRequest) (*proxypb.SubscribeReply, error) {
	client := p.GetClient()
	return client.Subscribe(ctx, req)
}
