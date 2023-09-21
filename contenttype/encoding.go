package contenttype

import (
	sharedpb "github.com/weflux/loop/protocol/shared"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"strings"
)

type ContentType string

const (
	JSON     ContentType = "json"
	PROTOBUF ContentType = "protobuf"
)

func (ct ContentType) ToProto() sharedpb.ContentType {
	c := FromString(string(ct))
	switch c {
	case JSON:
		return sharedpb.ContentType_JSON
	case PROTOBUF:
		return sharedpb.ContentType_PROTOBUF
	default:
		return sharedpb.ContentType_UNDEFINED
	}
}

func (ct ContentType) Marshal(msg proto.Message) ([]byte, error) {
	if ct == JSON {
		return protojson.Marshal(msg)
	}

	return proto.Marshal(msg)
}

func (ct ContentType) Unmarshal(data []byte, msg proto.Message) error {
	if ct == JSON {
		return protojson.Unmarshal(data, msg)
	}
	return proto.Unmarshal(data, msg)
}

func FromString(e string) ContentType {
	return ContentType(strings.ToLower(e))
}
