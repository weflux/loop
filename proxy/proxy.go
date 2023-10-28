package proxy

import (
	"github.com/weflux/loopify/option"
)

//goland:noinspection GoNameStartsWithPackageName
type ProxyMap struct {
	ConnectProxy   ConnectProxy
	SubscribeProxy SubscribeProxy
	RPCProxies     map[string]RPCProxy
}

func NewProxyMap(opts *option.ProxyOption) *ProxyMap {
	var connectProxy ConnectProxy
	var subscribeProxy SubscribeProxy

	if opts.Connect != nil {
		connectProxy = NewConnectProxy(opts.Connect)
	}

	if opts.Subscribe != nil {
		subscribeProxy = NewSubscribeProxy(opts.Subscribe)
	}

	var rpcProxies map[string]RPCProxy
	rpcProxies = map[string]RPCProxy{}
	for k, r := range opts.RPC {
		rpcProxies[k] = NewRPCProxy(r)
	}
	// TODO

	return &ProxyMap{
		ConnectProxy:   connectProxy,
		SubscribeProxy: subscribeProxy,
		RPCProxies:     rpcProxies,
	}
}
