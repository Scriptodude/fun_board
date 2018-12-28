package core

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	i "server/interfaces"
	"strconv"
	"strings"
)

type WebSocketServer struct {
	port     string
	addr     string
	ctx      context.Context
	cancel   context.CancelFunc
	listener net.Listener
	connList []net.Conn
}

func NewWebSocketServer(addr string, port int64) *WebSocketServer {
	ctx, cancel := context.WithCancel(context.Background())

	ws := WebSocketServer{
		port:   strconv.FormatInt(port, 10),
		addr:   addr,
		ctx:    ctx,
		cancel: cancel,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))

	if err != nil {
		panic(err)
	}

	ws.listener = listener
	return &ws
}

func (s *WebSocketServer) ListenAndServe() error {

	// Simple listener on the context to close the server when required
	go func() {
		_ = <-s.ctx.Done()
		s.listener.Close()
	}()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			return err
		}

		go s.manageConnection(conn)
	}
}

func (s *WebSocketServer) manageConnection(conn net.Conn) {
	//for {
	select {
	case <-s.ctx.Done():
		conn.Close()
		return
	default:

		reader := bufio.NewReader(conn)
		request, err := http.ReadRequest(reader)

		if err != nil {
			Error.Println(err)
			conn.Close()
		}

		request.ParseForm()
		Request.Printf("Received Request %s %s with data %+v", request.Method, request.URL.String(), request.Form)

		if request.Form.Get("data") == "exit" {
			conn.Close()
		}

		if request.Method == "GET" {
			serveFile(request, conn)
		}
	}
	//}
}

func (s *WebSocketServer) GetAddress() (string, error) {
	return s.addr, nil
}

func (s *WebSocketServer) GetPort() (string, error) {
	return s.port, nil
}

func (s *WebSocketServer) AwaitClient() *i.GameClient {
	return nil
}

func (s *WebSocketServer) Shutdown() {
	s.listener.Close()
	s.cancel()
}

func (s *WebSocketServer) AddRequestListener(path string, fn func(w http.ResponseWriter, r *http.Request, client *i.GameClient)) {
}

func (s *WebSocketServer) GetContext() context.Context {
	return s.ctx
}

func serveFile(req *http.Request, conn net.Conn) {
	url := strings.Replace(req.URL.String(), "..", "", -1)
	lastDot := strings.LastIndex(req.URL.String(), ".")
	lastSlash := strings.LastIndex(req.URL.String(), "/")
	status := 200

	// Small hack to check if it's a directory
	if lastDot < lastSlash || lastDot == -1 {
		url += "index.html"
		req.Header.Set("Content-Type", "text/html")
	} else {
		val := url[lastDot+1]

		switch val {
		case 'j':
			req.Header.Set("Content-Type", "text/javascript")
		case 'c':
			req.Header.Set("Content-Type", "text/css")
		default:
			req.Header.Set("Content-Type", "text/plain")
		}
	}

	data, err := ioutil.ReadFile("/tmp/fun_board" + url)

	if err != nil {
		Error.Println(err)
		status = 404
		data = []byte{}
	}

	writeHttp(conn, status, req.Header, data)
	conn.Close()
}

func writeHttp(conn net.Conn, status int, header http.Header, body []byte) {

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s", status, http.StatusText(status))
	headers := header
	writer := bufio.NewWriter(conn)

	writer.WriteString(statusLine)
	headers.Write(writer)
	writer.WriteByte('\n')
	writer.Write(body)

	writer.Flush()
}