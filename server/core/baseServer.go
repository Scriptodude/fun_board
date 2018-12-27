package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	prot "server/core/protocols"
	i "server/interfaces"
	"strings"
	"time"
)

var (
	fileServer http.Handler
)

type BaseServer struct {
	server     *http.Server
	mux        *http.ServeMux
	newClients chan *i.GameClient
	currentId  int
	clients    []*i.GameClient
	ctx        context.Context
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
		mux:        mux,
		newClients: make(chan *i.GameClient),
		currentId:  1,
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

func (b *BaseServer) AwaitClient() *i.GameClient {
	return <-b.newClients
}

func (b *BaseServer) Shutdown() {
	for _, client := range b.clients {
		Info.Printf("Closing client %d", client.Id)
		close(client.Messages)
	}
	close(b.newClients)
	b.server.Shutdown(b.ctx)
}

func (b *BaseServer) AddRequestListener(
	path string,
	fn func(w http.ResponseWriter, r *http.Request, client *i.GameClient)) {

	b.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		Request.Printf("%s %s\n", r.Method, r.URL.String())

		client, _ := b.getClient(w, r)

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
	client, isNew := b.getClient(w, r)

	// Find a way to manage concurrency
	if isNew {
		b.newClients <- client
		b.clients = append(b.clients, client)
	}

	b.longPoll(w, client)
}

/* Gets the client id associated to the request,
if None, gets a new one */
func (b *BaseServer) getClient(w http.ResponseWriter, r *http.Request) (*i.GameClient, bool) {
	reader := i.GameClient{}

	err := json.NewDecoder(r.Body).Decode(&reader)
	if err != nil || reader.Id == 0 {
		// Shit went wrong, just create a new client
		client := i.GameClient{Messages: make(chan string), Writer: w}
		go prot.NewClient(&client)
		Info.Printf("New Client : %d\n", client.Id)

		return &client, true
	}

	client := b.clients[reader.Id]
	Info.Printf("Returning Existing Client : %+v\n", client)
	go prot.ExistingClient(client)
	return client, false
}

func (b *BaseServer) longPoll(w http.ResponseWriter, client *i.GameClient) {
	notifier, ok := w.(http.CloseNotifier)
	if !ok {
		panic("Expected http.ResponseWriter to be an http.CloseNotifier")
	}

	Info.Printf("Starting long poll for %d\n", client.Id)
	ctx, cancel := context.WithCancel(context.Background())
	b.ctx = ctx

	select {
	case msg := <-client.Messages:
		Info.Printf("Received message %s for client %d", msg, client.Id)

		if msg == "done" {
			cancel()
			return
		}

		fmt.Fprint(w, msg)

	case <-time.After(time.Minute * 10):
		Info.Printf("Client %d will be disconnected\n", client.Id)
		cancel()
		return

	case <-notifier.CloseNotify():
		Info.Printf("Client %d has disconnected\n", client.Id)
		cancel()
		return
	}
}
