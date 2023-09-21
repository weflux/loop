package proxy

import (
	"context"
	"github.com/weflux/loop"
	"github.com/weflux/loop/contenttype"
	proxypb "github.com/weflux/loop/protocol/proxy"
	"net/url"
)

type RPCProxy interface {
	ProxyRPC(ctx context.Context, req *proxypb.RPCRequest) (*proxypb.RPCReply, error)
}

func NewRPCProxy(c *loop.RouteOption) RPCProxy {
	endpoint := c.Endpoint
	uri, err := url.Parse(endpoint)
	if err != nil {
		panic(err)
	}
	if uri.Scheme == "http" || uri.Scheme == "https" {
		return &httpRPCProxy{
			HttpBase: HttpBase{URL: endpoint, ContentType: contenttype.ContentType(c.ContentType)},
		}
	}

	return &grpcRPCProxy{
		GRPCBase: GRPCBase{Addr: c.Endpoint},
	}
}

type httpRPCProxy struct {
	HttpBase
}

func (h *httpRPCProxy) ProxyRPC(ctx context.Context, req *proxypb.RPCRequest) (*proxypb.RPCReply, error) {
	rep := &proxypb.RPCReply{}
	if err := h.ProxyHTTP(req, rep); err != nil {
		return nil, err
	}

	return rep, nil
}

type grpcRPCProxy struct {
	GRPCBase
}

func (p *grpcRPCProxy) ProxyRPC(ctx context.Context, req *proxypb.RPCRequest) (*proxypb.RPCReply, error) {
	client := proxypb.NewProxyClient(p.GetConn())
	return client.RPC(ctx, req)
}
