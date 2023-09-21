package hook

import (
	"bytes"
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/weflux/loop/proxy"
	"github.com/weflux/loop/utils/packetutil"
)

type Auth struct {
	mqtt.HookBase
	proxyMap *proxy.ProxyMap
}

func NewAuth(proxyMap *proxy.ProxyMap) *Auth {
	return &Auth{
		HookBase: mqtt.HookBase{},
		proxyMap: proxyMap,
	}
}

func (h *Auth) ID() string {
	return "auth-hook"
}

func (h *Auth) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnSessionEstablished,
	}, []byte{b})
}

func (h *Auth) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	h.Log.Debug(fmt.Sprintf("[id=%s] try connect authenticate: %s %s", cl.ID, string(pk.Connect.Username), string(pk.Connect.Password)))
	if h.proxyMap.AuthenticateProxy != nil {
		req := packetutil.ToConnectRequest(cl, &pk)
		rep, err := h.proxyMap.AuthenticateProxy.ProxyAuthenticate(emptyCtx(), req)
		if err != nil {
			h.Log.Error("authenticate failed", "error", err)
			return false
		}
		if rep.Error == nil {
			return true
		}

		h.Log.Error("msg", "authenticate failed", "errno", rep.Error.Code, "errmsg", rep.Error.Message)
		return false
	} else {
		return h.HookBase.OnConnectAuthenticate(cl, pk)
	}
}

func (h *Auth) OnSessionEstablished(cl *mqtt.Client, pk packets.Packet) {
}
