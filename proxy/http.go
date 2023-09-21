package proxy

import (
	"bytes"
	"errors"
	"github.com/weflux/loop/contenttype"
	"google.golang.org/protobuf/proto"
	"io"
	"net/http"
)

type HttpBase struct {
	URL         string
	ContentType contenttype.ContentType
}

func (h *HttpBase) ProxyHTTP(req proto.Message, rep proto.Message) error {
	var ct string
	if h.ContentType == contenttype.JSON {
		ct = "application/json"
	}
	data, err := h.ContentType.Marshal(req)
	if err != nil {
		return err
	}
	r, err := http.Post(h.URL, ct, bytes.NewReader(data))
	if err != nil {
		return err
	}
	if r.StatusCode != 200 {
		return errors.New("proxy error: " + r.Status)
	}
	data, err = io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := h.ContentType.Unmarshal(data, rep); err != nil {
		return err
	}

	return nil
}
