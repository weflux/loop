package clientutil

import (
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/weflux/loopify/contenttype"
)

func ContentType(cl *mqtt.Client) contenttype.ContentType {
	for _, kv := range cl.Properties.Props.User {
		if kv.Key == "content-type" {
			return contenttype.ContentType(kv.Val)
		}
	}
	return contenttype.JSON
}

func User(cl *mqtt.Client) string {
	return string(cl.Properties.Username)
}
