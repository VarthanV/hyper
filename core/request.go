package core

import "net"

type Request struct {
	Protocol       string
	Path           string
	Method         HttpMethod
	Body           []byte
	RemoteHostAddr *net.Addr
	headers        map[string]string
}
