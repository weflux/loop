package proxy

import (
	proxypb "github.com/weflux/loop/protocol/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCBase struct {
	Addr string
	cc   *grpc.ClientConn
}

func (r *GRPCBase) GetClient() proxypb.ProxyClient {
	return proxypb.NewProxyClient(r.GetConn())
}

func (r *GRPCBase) GetConn() *grpc.ClientConn {
	if r.cc == nil {
		var err error
		r.cc, err = grpc.Dial(r.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
	}
	return r.cc
}
