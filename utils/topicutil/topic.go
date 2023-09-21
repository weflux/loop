package topicutil

import (
	"fmt"
	mqtt "github.com/mochi-mqtt/server/v2"
)

func ClientTopic(cl *mqtt.Client) string {
	return fmt.Sprintf("$CLIENT/%s", cl.ID)
}

func RequestTopic(route string) string {
	return fmt.Sprintf("$RPC/%s", route)
}

func ReplyTopic() string {
	return "$RPC/reply"
}
