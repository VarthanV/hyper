package core

import (
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	contentTypeHeader = "Content-Type"
)

type Response struct {
	statusCode int
	headers    map[string]string
	body       []byte
}

// WriteStatus: Write the status code to be returned
func (r *Response) WriteStatus(code int) {
	r.statusCode = code
}

// WriteHeader: Write response headers
func (r *Response) WriteHeader(key, val string) {
	r.headers[key] = val
}

// WriteJSON: Writes json to the body
func (r *Response) WriteJSON(status int, b interface{}) {
	marshalledBytes, err := json.Marshal(b)
	if err != nil {
		panic(errors.Wrap(err, "error when writing json unable to marshal").Error())
	}

	r.headers[contentTypeHeader] = "application/json"
	r.statusCode = status
	r.body = marshalledBytes
}

// WriteHTML: Writes HTML to the body
func (r *Response) WriteHTML(status int, html string) {
	r.headers[contentTypeHeader] = "text/html"
	r.statusCode = status
	r.body = []byte(html)
}

// WriteString: Writes string to the body
func (r *Response) WriteString(status int, val string) {
	r.body = []byte(val)
	r.statusCode = status
	r.headers[contentTypeHeader] = "text/plain"
}
