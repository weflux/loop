package proxy

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/weflux/loop/protocol/message/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"testing"
)

func TestOneOfMarshalJSON(t *testing.T) {
	msg := &message.Message{
		Id:      uuid.NewString(),
		Headers: map[string]string{},
		Body: &message.Message_Error{Error: &message.Error{
			Code:    1001,
			Message: "error",
			Extras:  map[string]string{},
		}},
	}

	bs, _ := protojson.Marshal(msg)
	fmt.Println(string(bs))
}

func TestOneOfUnmarshalJSON(t *testing.T) {
	js := "{\"id\":\"e3d69d3e-3572-4919-815e-beb9c70735f5\", \"encoding\":\"JSON\", \"error\":{\"code\":1001, \"message\":\"error\", \"data\":\"YWJjZGVmZw==\"}}"
	msg := &message.Message{}
	err := protojson.Unmarshal([]byte(js), msg)
	require.NoError(t, err)
	switch payload := msg.Body.(type) {
	case *message.Message_Error:
		require.EqualValues(t, 1001, payload.Error.Code)
	}
}
