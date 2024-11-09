package core

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
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

	printStartupMessage := func() {
		for msg := range startupMessageChan {
			log.Println(msg)
		}
	}

	go printStartupMessage()

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

		h.handleConnection(conn)

	}
}

func (h *hyper) handleConnection(c net.Conn) {
	var (
		maxConnectionCanHandleConcurrently = runtime.NumCPU() * 10 // 10 times the magnitude of num cpu
		sem                                = make(chan struct{}, maxConnectionCanHandleConcurrently)
	)

	handleConn := func(c net.Conn) {

		// Put in semaphore that we are handling the connection
		sem <- struct{}{}
		fmt.Println("Connection from ", c.RemoteAddr().String())

		// Release the sem
		<-sem
		_, err := h.parseRequest(c)
		if err != nil {
			log.Println("error in parsing request")
			c.Close()
		}
	}

	go handleConn(c)

}

// parseRequest: Reads the raw http request and parses it into “Request“ struct
func (h *hyper) parseRequest(conn net.Conn) (*Request, error) {
	request := &Request{}
	reader := bufio.NewReader(conn)

	// Parse the equest line (e.g., "GET /path HTTP/1.1")
	requestLine, err := reader.ReadString(delimNewLine)
	if err != nil {
		log.Println("error in reading request line ", err)
		return nil, err
	}
	log.Println("Request line ", requestLine)

	// Split into parts

	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		log.Println(ErrInvalidRequestLine.Error())
		return nil, ErrInvalidRequestLine
	}

	request.Method = HttpMethod(parts[0])
	request.Path = parts[1]
	request.Protocol = parts[2]

	log.Printf("%+v", request)

	return nil, nil
}
