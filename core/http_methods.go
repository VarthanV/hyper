package core

type HttpMethod string

const (
	HttpMethodGet     HttpMethod = "GET"
	HttpMethodPost    HttpMethod = "POST"
	HttpMethodPut     HttpMethod = "PUT"
	HttpMethodDelete  HttpMethod = "DELETE"
	HttpMethodPatch   HttpMethod = "PATCH"
	HttpMethodHead    HttpMethod = "HEAD"
	HttpMethodConnect HttpMethod = "CONNECT"
	HttpMethodOptions HttpMethod = "OPTIONS"
	HttpMethodTrace   HttpMethod = "TRACE"
)
