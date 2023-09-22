package proxy

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"testing"
)

func TestOneOfMarshalJSON(t *testing.T) {
	msg := &envelope.Message{
		Id:      uuid.NewString(),
		Headers: map[string]string{},
		Body: &envelope.Message_Error{Error: &envelope.Error{
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
	msg := &envelope.Message{}
	err := protojson.Unmarshal([]byte(js), msg)
	require.NoError(t, err)
	switch payload := msg.Body.(type) {
	case *envelope.Message_Error:
		require.EqualValues(t, 1001, payload.Error.Code)
	}
}
