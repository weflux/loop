package contenttype

import (
	"encoding/json"
	"errors"
	"google.golang.org/protobuf/proto"
)

type PayloadGetter interface {
	GetContentType() string
	GetPayloadText() string
	GetPayloadBytes() []byte
}

type PayloadSetter interface {
	SetPayloadText(p string)
	SetPayloadBytes(p []byte)
}

func UnmarshalPayload(msg PayloadGetter, v interface{}) error {
	ct := FromString(msg.GetContentType())
	switch ct {
	case JSON:
		m, ok := v.(proto.Message)
		if ok {
			return ct.Unmarshal([]byte(msg.GetPayloadText()), m)
		} else {
			return json.Unmarshal([]byte(msg.GetPayloadText()), v)
		}
	case PROTOBUF:
		m, ok := v.(proto.Message)
		if ok {
			return ct.Unmarshal([]byte(msg.GetPayloadText()), m)
		} else {
			return errors.New("protobuf data must unmarshal to proto.Message")
		}
	}

	return errors.New("not support content-type")
}
