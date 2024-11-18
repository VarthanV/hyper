package hyper_test

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/VarthanV/hyper"
	"github.com/stretchr/testify/assert"
)

func TestServerStartup(t *testing.T) {
	server := hyper.New()
	go server.ListenAndServe("127.0.0.1", "8081", "Server started")

	// Wait briefly to ensure server starts
	time.Sleep(1 * time.Second)

	// Check server availability
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if assert.NoError(t, err, "Server failed to start") {
		conn.Close()
	}
}

func TestRequestHandling(t *testing.T) {
	server := hyper.New()
	expectedResponse := "Hello, World!"

	// Define a handler
	server.GET("/hello", func(w hyper.ResponseWriter, r *hyper.Request) {
		w.Write([]byte(expectedResponse))
	})

	go server.ListenAndServe("127.0.0.1", "8082", "Server started")

	// Wait for server to start
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:8082")
	if assert.NoError(t, err) {
		defer conn.Close()

		// Send HTTP GET request
		request := "GET /hello HTTP/1.1\r\nHost: 127.0.0.1\r\n\r\n"
		_, err := conn.Write([]byte(request))
		assert.NoError(t, err)

		// Read response
		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, conn)
		assert.NoError(t, err)

		// Check response body
		assert.Contains(t, buf.String(), expectedResponse)
	}
}
