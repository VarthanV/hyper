package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

const (
	contentTypeHeader = "Content-Type"
	contentLength     = "Content-Length"
)

type ResponseWriter interface {
	WriteStatus(code int)
	Write([]byte) (int, error)
	WriteHeader(key, val string)
	WriteJSON(status int, b interface{})
	WriteHTML(status int, html string)
	WriteString(status int, val string)
	ToRaw() string
	StatusCode() int
}

type responseWriter struct {
	statusCode    int
	headers       map[string]string
	body          bytes.Buffer
	templatesPath string
}

func newResponseWriter(templatesPath string) ResponseWriter {
	return &responseWriter{
		headers:       make(map[string]string),
		templatesPath: templatesPath,
	}
}

// WriteStatus: Write the status code to be returned
func (r *responseWriter) WriteStatus(code int) {
	r.statusCode = code
}

// WriteHeader: Write response headers
func (r *responseWriter) WriteHeader(key, val string) {
	r.headers[key] = val
}

// WriteJSON: Writes json to the body
func (r *responseWriter) WriteJSON(status int, b interface{}) {
	marshalledBytes, err := json.Marshal(b)
	if err != nil {
		panic(errors.Wrap(err, "error when writing json unable to marshal").Error())
	}

	r.headers[contentTypeHeader] = "application/json"
	r.statusCode = status
	r.Write(marshalledBytes)
}

// WriteHTML: Writes HTML to the body
func (r *responseWriter) WriteHTML(status int, html string) {
	r.headers[contentTypeHeader] = "text/html"
	r.statusCode = status
	r.Write([]byte(html))
}

// WriteString: Writes string to the body
func (r *responseWriter) WriteString(status int, val string) {
	r.statusCode = status
	r.headers[contentTypeHeader] = "text/plain"
	r.Write([]byte(val))

}

// Write: Writes raw byte to response
func (r *responseWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

// ToRaw converts the Response struct to a raw HTTP response string
func (r *responseWriter) ToRaw() string {
	statusText := getStatusText(r.statusCode)

	r.headers[contentLength] = strconv.Itoa(r.body.Len())
	raw := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.statusCode, statusText)

	// Append headers
	for key, value := range r.headers {
		raw += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	// Add a blank line to separate headers from body
	raw += "\r\n"

	// Append the body if present
	raw += r.body.String()

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

// StatusCode implements ResponseWriter.
func (r *responseWriter) StatusCode() int {
	return r.statusCode
}
