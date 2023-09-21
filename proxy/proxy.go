package proxy

import (
	"github.com/weflux/loop/option"
)

//goland:noinspection GoNameStartsWithPackageName
type ProxyMap struct {
	AuthenticateProxy AuthenticateProxy
	SubscribeProxy    SubscribeProxy
	RPCProxies        map[string]RPCProxy
}

func NewProxyMap(opts *option.ProxyOption) *ProxyMap {
	var connectProxy AuthenticateProxy

	if opts.Authenticate != nil {
		connectProxy = NewAuthenticateProxy(opts.Authenticate)
	}

	var rpcProxies map[string]RPCProxy
	rpcProxies = map[string]RPCProxy{}
	for k, r := range opts.RPC {
		rpcProxies[k] = NewRPCProxy(r)
	}
	// TODO

	return &ProxyMap{
		AuthenticateProxy: connectProxy,
		RPCProxies:        rpcProxies,
	}
}
