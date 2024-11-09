package core

import (
	"fmt"
	"log"
	"net"
)

type hyper struct {
	routesMap map[HttpMethod]map[string]HandlerFunc
}

func New() *hyper {
	return &hyper{
		routesMap: make(map[HttpMethod]map[string]HandlerFunc),
	}
}

func (h *hyper) ListenAndServe(host, port, startupMessage string) {
	var (
		startupMessageChan = make(chan string, 1)
	)

	go func() {
		for msg := range startupMessageChan {
			log.Println(msg)
		}
	}()

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatal("error in listening ", err)
	}
	defer l.Close()
	log.Printf("Listening on %s:%s\n", host, port)

	startupMessageChan <- startupMessage
	close(startupMessageChan)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("unable to accept connection ", err)
			continue
		}

		go func(c net.Conn) {
			log.Println("Connection is ", c.RemoteAddr().String())
		}(conn)
	}
}
