package proxy

import (
	"context"
	proxypb "github.com/weflux/loop/protocol/proxy"
)

type SubscribeProxy interface {
	ProxySubscribe(ctx context.Context, req *proxypb.SubscribeRequest) (*proxypb.SubscribeReply, error)
	Protocol() ProtocolType
}
