package core

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

const (
	contentTypeHeader = "Content-Type"
	contentLength     = "Content-Length"
)

type ResponseWriter struct {
	statusCode int
	headers    map[string]string
	body       []byte
}

func newResponse() *ResponseWriter {
	return &ResponseWriter{
		headers: make(map[string]string),
	}
}

// WriteStatus: Write the status code to be returned
func (r *ResponseWriter) WriteStatus(code int) {
	r.statusCode = code
}

// WriteHeader: Write response headers
func (r *ResponseWriter) WriteHeader(key, val string) {
	r.headers[key] = val
}

// WriteJSON: Writes json to the body
func (r *ResponseWriter) WriteJSON(status int, b interface{}) {
	marshalledBytes, err := json.Marshal(b)
	if err != nil {
		panic(errors.Wrap(err, "error when writing json unable to marshal").Error())
	}

	r.headers[contentTypeHeader] = "application/json"
	r.statusCode = status
	r.body = marshalledBytes
}

// WriteHTML: Writes HTML to the body
func (r *ResponseWriter) WriteHTML(status int, html string) {
	r.headers[contentTypeHeader] = "text/html"
	r.statusCode = status
	r.body = []byte(html)
}

// WriteString: Writes string to the body
func (r *ResponseWriter) WriteString(status int, val string) {
	r.body = []byte(val)
	r.statusCode = status
	r.headers[contentTypeHeader] = "text/plain"
}

// Write: Writes raw byte to response
func (r *ResponseWriter) Write(b []byte) {
	r.body = b
}

// ToRaw converts the Response struct to a raw HTTP response string
func (r *ResponseWriter) ToRaw() string {
	statusText := getStatusText(r.statusCode)

	r.headers[contentLength] = strconv.Itoa(len(r.body))
	raw := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.statusCode, statusText)

	// Append headers
	for key, value := range r.headers {
		raw += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Add a blank line to separate headers from body
	raw += "\r\n"

	// Append the body if present
	raw += string(r.body)

	return raw
}

func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Unknown Status"
	}
}
