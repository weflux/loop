package hook

import (
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
)

type ACL struct {
	mqtt.HookBase
}

func NewACL() *ACL {
	return &ACL{}
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
	//return h.HookBase.OnACLCheck(cl, topic, write)
	// TODO
	//h.Log.Info("OnACLCheck: ", json.MustMarshalToString(cl), topic, write)
	return true
}
