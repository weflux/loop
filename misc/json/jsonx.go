package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

func init() {
	extra.RegisterFuzzyDecoders()
}

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(v interface{}) ([]byte, error) {
	return JSON.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return JSON.Unmarshal(data, v)
}

func MustMarshalToString(v interface{}) string {
	if s, err := JSON.Marshal(v); err != nil {
		panic(err)
	} else {
		return string(s)
	}
}

func MustMarshal(v interface{}) []byte {
	if s, err := JSON.Marshal(v); err != nil {
		panic(err)
	} else {
		return s
	}
}
