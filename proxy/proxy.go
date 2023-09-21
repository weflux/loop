package proxy

import (
	"github.com/weflux/loop"
)

//goland:noinspection GoNameStartsWithPackageName
type ProxyMap struct {
	ConnectProxy   ConnectProxy
	SubscribeProxy SubscribeProxy
	RPCProxies     map[string]RPCProxy
}

func NewProxyMap(opts *loop.ProxyOption) *ProxyMap {
	var connectProxy ConnectProxy

	if opts.Connect != nil {
		connectProxy = NewConnectProxy(opts.Connect)
	}

	var rpcProxies map[string]RPCProxy
	rpcProxies = map[string]RPCProxy{}
	for k, r := range opts.RPC {
		rpcProxies[k] = NewRPCProxy(r)
	}
	// TODO

	return &ProxyMap{
		ConnectProxy: connectProxy,
		RPCProxies:   rpcProxies,
	}
}
