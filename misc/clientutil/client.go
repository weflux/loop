package clientutil

import (
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/weflux/loop/misc/contenttype"
)

func GetContentType(cl *mqtt.Client) contenttype.ContentType {
	for _, kv := range cl.Properties.Props.User {
		if kv.Key == "content-type" {
			return contenttype.ContentType(kv.Val)
		}
	}
	return contenttype.JSON
}
