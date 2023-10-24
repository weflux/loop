package proxy

import (
	"context"
	"github.com/weflux/loopify/contenttype"
	"github.com/weflux/loopify/option"
	proxypb "github.com/weflux/loopify/protocol/proxy"
	"net/url"
)

type ConnectProxy interface {
	ProxyConnect(ctx context.Context, req *proxypb.ConnectRequest) (*proxypb.ConnectReply, error)
}

func NewConnectProxy(c *option.RouteOption) ConnectProxy {

	endpoint := c.Endpoint
	uri, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	if uri.Scheme == "http" || uri.Scheme == "https" {
		return &httpConnectProxy{
			HttpBase: HttpBase{URL: endpoint, ContentType: contenttype.ContentType(c.ContentType)},
		}
	}

	return &grpcConnectProxy{
		GRPCBase: GRPCBase{Addr: c.Endpoint},
	}
}

type httpConnectProxy struct {
	HttpBase
}

func (h *httpConnectProxy) ProxyConnect(ctx context.Context, req *proxypb.ConnectRequest) (*proxypb.ConnectReply, error) {
	rep := &proxypb.ConnectReply{}
	if err := h.ProxyHTTP(req, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type grpcConnectProxy struct {
	GRPCBase
}

func (p *grpcConnectProxy) ProxyConnect(ctx context.Context, req *proxypb.ConnectRequest) (*proxypb.ConnectReply, error) {
	client := p.GetClient()
	return client.Connect(ctx, req)
}
