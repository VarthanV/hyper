package core

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/VarthanV/hyper/pkg/runtimeutils"
)

type hyper struct {
	routes     map[HttpMethod]map[string]routeStruct
	staticPath string
}

func New() *hyper {
	return &hyper{
		routes:     make(map[HttpMethod]map[string]routeStruct),
		staticPath: "",
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

	// Print the routes:handler mapping
	for method, routes := range h.routes {
		for _, f := range routes {
			fmt.Printf("%s %s ----------------------> %s\n", method, f.UserGivenPath, runtimeutils.GetFunctionName(f.Handler))
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

func (h *hyper) mapHandlers(path string, method HttpMethod, handler HandlerFunc) {
	if _, ok := h.routes[method]; !ok {
		h.routes[method] = make(map[string]routeStruct)
	}

	compliledPath := path
	if !strings.HasPrefix(path, "^") {
		compliledPath = fmt.Sprintf(`^%s$`, path)
	}

	h.routes[method][compliledPath] = routeStruct{UserGivenPath: path, Handler: handler}
}

func (h *hyper) POST(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodPost, handler)
}

func (h *hyper) PUT(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodPut, handler)
}

func (h *hyper) PATCH(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodPatch, handler)
}

func (h *hyper) GET(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodGet, handler)
}

func (h *hyper) DELETE(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodPatch, handler)
}

func (h *hyper) OPTIONS(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodOptions, handler)
}

func (h *hyper) CONNECT(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodConnect, handler)
}

func (h *hyper) TRACE(path string, handler HandlerFunc) {
	h.mapHandlers(path, HttpMethodTrace, handler)
}

func (h *hyper) handleConnection(c net.Conn) {
	var (
		maxConnectionCanHandleConcurrently = runtime.NumCPU() * 10 // 10 times the magnitude of num cpu
		sem                                = make(chan struct{}, maxConnectionCanHandleConcurrently)
	)

	handleConn := func(c net.Conn) {
		defer c.Close()
		// Put in semaphore that we are handling the connection
		sem <- struct{}{}
		fmt.Println("Connection from ", c.RemoteAddr().String())

		// Release the sem
		<-sem
		res := newResponse()

		request, err := h.parseRequest(c)
		if err != nil {
			log.Println("error in parsing request ", err)
			return
		}

		log.Printf("%+v", request)

		if request != nil {
			handlerMap, ok := h.routes[request.Method]
			if !ok {
				log.Println("handler not found for method")
				return
			}

			isHandlerMatched := false

			for path, r := range handlerMap {
				compiledPath := regexp.MustCompile(path)
				isHandlerMatched = compiledPath.MatchString(request.Path)
				if isHandlerMatched {
					// Check if it matches the /:x param
					re := regexp.MustCompile(`^/([^/]+)$`)
					match := re.FindStringSubmatch(path)
					if match != nil {
						log.Println("match found with id ", 1)
					}

					r.Handler(request, res)
				}
			}
		}

		_, err = c.Write([]byte(res.ToRaw()))
		if err != nil {
			log.Println("unable to write to conn ", err)
		}

		log.Printf("%s %s  %d", request.Method, request.Path, res.statusCode)
	}

	go handleConn(c)

}

// parseRequest: Reads the raw http request and parses it into “Request“ struct
func (h *hyper) parseRequest(conn net.Conn) (*Request, error) {
	request := &Request{
		RemoteHostAddr: conn.RemoteAddr(),
	}
	reader := bufio.NewReader(conn)

	// Parse the equest line (e.g., "GET /path HTTP/1.1")
	requestLine, err := reader.ReadString(delimNewLine)
	if err != nil {
		return nil, err
	}
	// Split into parts

	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		log.Println(ErrInvalidRequestLine.Error())
		return nil, ErrInvalidRequestLine
	}

	request.Method = HttpMethod(parts[0])
	request.Path = parts[1]
	request.Protocol = parts[2]

	request.headers = make(map[string]string)

	queryParams := make(map[string]string)
	if idx := strings.Index(request.Path, "?"); idx != -1 {
		queryString := request.Path[idx+1:]
		request.Path = request.Path[:idx]

		queryPairs := strings.Split(queryString, "&")
		for _, pair := range queryPairs {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				key := kv[0]
				value := kv[1]
				queryParams[key] = value
			}
		}
	}

	request.queryParams = queryParams

	// Populate headers
	for {
		line, err := reader.ReadString(delimNewLine)
		if err != nil {
			fmt.Println("Error reading header:", err)
			return nil, err
		}
		line = strings.TrimSpace(line)
		// Break if we reach an empty line (end of headers)
		if line == "" {
			break
		}

		colonIndex := strings.Index(line, ":")
		if colonIndex == -1 {
			continue
		}
		key := strings.TrimSpace(line[:colonIndex])
		value := strings.TrimSpace(line[colonIndex+1:])
		request.headers[key] = value
	}

	//  If there’s a Content-Length header, read the body
	if l, ok := request.headers["Content-Length"]; ok {

		lint, err := strconv.Atoi(l)
		if err != nil {
			return nil, err
		}

		body := make([]byte, lint)
		_, err = reader.Read(body)
		if err != nil {
			return nil, err
		}
		request.Body = body
	}

	return request, nil
}

/*
Sample raw request
	POST /submit-form HTTP/1.1
	Host: www.example.com
	User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36
	Content-Type: application/x-www-form-urlencoded
	Content-Length: 27
	Connection: keep-alive

	username=johndoe&password=1234

*/
