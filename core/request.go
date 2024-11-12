package core

import (
	"encoding/json"
	"net"

	"github.com/pkg/errors"
)

type Request struct {
	Protocol       string
	Path           string
	Method         HttpMethod
	Body           []byte
	RemoteHostAddr net.Addr
	queryParams    map[string]string
	headers        map[string]string
	ContentLength  int64
	pathParams     map[string]string
}

func (r *Request) GetHeader(key string) string {
	return r.headers[key]
}

func (r *Request) Query(key string) string {
	return r.queryParams[key]
}

func (r *Request) Bind(i interface{}) error {
	err := json.Unmarshal(r.Body, i)
	if err != nil {
		return errors.Wrap(err, "unable to bind request")
	}
	return nil
}

func (r *Request) Param(key string) string {
	return r.pathParams[key]
}
