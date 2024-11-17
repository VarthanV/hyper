# Hyper Package

`hyper` is a lightweight toy and extensible HTTP server framework for building TCP-based web applications in Go.

## Features

- **Customizable Route Handlers**: Define HTTP methods like `GET`, `POST`, `PUT`, etc., and their handlers.
- **Dynamic Path Matching**: Supports route patterns with regex for flexibility.
- **Query and Header Parsing**: Automatically parses query parameters and HTTP headers.
- **Concurrency**: Efficient connection handling using goroutines and semaphores.
- **Template Support**: Configurable templates directory for rendering responses.

---

## Installation

```bash
go get github.com/VarthanV/hyper
```

---

## Usage

### Basic Server Setup

```go
package main

import (
	"fmt"
	"github.com/VarthanV/hyper"
)

func main() {
	server := hyper.New()

	// Define routes
	server.GET("/hello", func(w hyper.ResponseWriter, r *hyper.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Start the server
	server.ListenAndServe("127.0.0.1", "8080", "Server is running...")
}
```

### Adding Routes

Define routes with HTTP methods using the built-in methods:

- `GET`
- `POST`
- `PUT`
- `PATCH`
- `DELETE`
- `OPTIONS`
- `CONNECT`
- `TRACE`

Example:

```go
server.POST("/submit", func(w hyper.ResponseWriter, r *hyper.Request) {
	fmt.Fprintf(w, "Received POST data: %s", string(r.Body))
})
```

---

## API Reference

### `func New() *hyper`

Creates a new instance of the Hyper server.

### `func (*hyper) ListenAndServe(host, port, startupMessage string)`

Starts the server on the specified host and port. Logs the `startupMessage`.

### `func (*hyper) <METHOD>(path string, handler HandlerFunc)`

Registers a route for the specified HTTP method (`GET`, `POST`, etc.).

- **Parameters:**
  - `path`: Route path (supports regex).
  - `handler`: Function to handle the request.

---

## Request and Response Handling

### Request Structure

```go
type Request struct {
	Method          HttpMethod
	Path            string
	Protocol        string
	headers         map[string]string
	queryParams     map[string]string
	Body            []byte
	RemoteHostAddr  net.Addr
}
```

### Response Writer

```go
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
```

---

## Sample application link
[hyper-blog](https://github.com/VarthanV/hyper-blog)

## License

This project is licensed under the MIT License.