package core

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	prot "server/core/protocols"
	"strings"
	"time"
)

var (
	fileServer http.Handler
)

type BaseServer struct {
	server       *http.Server
	mux          *http.ServeMux
	clients      chan GameClient
	currentId    int
	hasNewClient bool
}

type LogFileHandler struct{}

func (l LogFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Request.Printf("%s %s\n", r.Method, r.URL.String())

	if fileServer == nil {
		fileServer = http.StripPrefix("/", http.FileServer(http.Dir("/tmp/fun_board")))
	}

	fileServer.ServeHTTP(w, r)
}

func NewDefaultServer() *BaseServer {
	mux := http.NewServeMux()

	server := &BaseServer{
		server: &http.Server{
			Addr:         ":8080",
			Handler:      mux,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
			ErrorLog:     Error,
		},
		mux:       mux,
		clients:   make(chan GameClient, 2),
		currentId: 1,
	}

	server.mux.Handle("/", LogFileHandler{})
	server.mux.HandleFunc("/connect", server.handleConnection)
	server.server.RegisterOnShutdown(func() { Info.Println("Shutting down the server...") })

	return server
}

func (b *BaseServer) ListenAndServe() error {
	return b.server.ListenAndServe()
}

func (b *BaseServer) GetAddress() (string, error) {
	if b.server == nil || b.server.Addr == "" {
		return "", errors.New("The server is not yet created.")
	}

	return b.server.Addr, nil
}

func (b *BaseServer) GetPort() (string, error) {
	if b.server == nil || b.server.Addr == "" {
		return "", errors.New("The server is not yet created.")
	}

	// Let us assume the Addr is in IPv4 for now
	split := strings.Split(b.server.Addr, ":")

	if len(split) < 1 {
		return "", errors.New("There was no port in the server address, maybe it isn't open yet ?")
	}

	return split[1], nil
}

func (b *BaseServer) AwaitClient() GameClient {
	return <-b.clients
}

func (b *BaseServer) Shutdown() {
	for i := 0; i <= b.currentId; {
		b.clients <- GameClient{}
	}
	b.server.Shutdown(context.Background())
}

func (b *BaseServer) AddRequestListener(
	path string,
	fn func(w http.ResponseWriter, r *http.Request, client GameClient)) {

	b.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		Request.Printf("%s %s\n", r.Method, r.URL.String())

		client := b.getClient(w, r)

		fn(w, r, client)
	})
}

func (b *BaseServer) handleConnection(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {

		w.Header().Add("Content-Type", "")
		// Since we don't want to support anything other than POST
		http.Redirect(w, r, "/", 405)
		return
	}

	Request.Printf("%s %s\n", r.Method, r.URL.String())
	client := b.getClient(w, r)

	// Find a way to manage concurrency
	if b.hasNewClient {
		b.hasNewClient = false
		b.clients <- client
	}
}

/* Gets the client id associated to the request,
if None, gets a new one */
func (b *BaseServer) getClient(w http.ResponseWriter, r *http.Request) GameClient {
	client := GameClient{}

	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil || client.Id == 0 {
		// Shit went wrong, just create a new client
		client.Id = b.newClient(w)
		return client
	}

	Info.Printf("Returning Existing Client : %+v\n", client)
	prot.GetClientIdMessage(w, client.Id)
	return client
}

func (b *BaseServer) newClient(w http.ResponseWriter) int {
	Info.Printf("New client; their ID is %d", b.currentId)
	prot.GetClientIdMessage(w, b.currentId)

	b.currentId += 1
	b.hasNewClient = true

	return b.currentId - 1
}
