package hook

import (
	"bytes"
	"context"
	mqtt "github.com/mochi-mqtt/server/v2"
	proxypb "github.com/weflux/loopify/protocol/proxy"
	"github.com/weflux/loopify/proxy"
	"github.com/weflux/loopify/utils/clientutil"
)

type ACL struct {
	mqtt.HookBase
	proxy *proxy.ProxyMap
}

func NewACL(proxy *proxy.ProxyMap) *ACL {
	return &ACL{
		proxy: proxy,
	}
}

func (h *ACL) ID() string {
	return "acl-hook"
}

func (h *ACL) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnACLCheck,
	}, []byte{b})
}

func (h *ACL) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	p := h.proxy.SubscribeProxy
	if p != nil {
		rep, err := p.ProxySubscribe(context.Background(), &proxypb.SubscribeRequest{
			Topic:  topic,
			Client: cl.ID,
			User:   clientutil.User(cl),
		})
		if err != nil {
			h.Log.Error("topic subscribe check error", "error", err)
			return false
		}
		if rep.Code != 0 {
			h.Log.Error("topic subscribe check failed", "code", rep.Code)
			return false
		}
		return true
	}
	h.Log.Info("no subscribe proxy config, use default `true`")
	return true
}
