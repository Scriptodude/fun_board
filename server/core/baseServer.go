package core

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	prot "server/core/protocols"
	"strconv"
	"strings"
	"time"
)

var (
	fileServer http.Handler
)

type BaseServer struct {
	server    *http.Server
	mux       *http.ServeMux
	clients   chan GameClient
	currentId int
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
		currentId: 0,
	}

	server.mux.Handle("/", LogFileHandler{})
	server.mux.HandleFunc("/connect", server.handleConnection)
	server.server.RegisterOnShutdown(func() { Info.Println("Shutting down the server...") })

	return server
}

func (b *BaseServer) ServeAndListen() error {
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
	b.server.Shutdown(context.Background())
}

func (b *BaseServer) AddRequestListener(
	path string,
	fn func(w http.ResponseWriter, r *http.Request, client GameClient)) {

	b.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		Request.Printf("%s %s\n", r.Method, r.URL.String())

		id := b.getClientId(w, r)

		fn(w, r, GameClient{id})
	})
}

/* Simple method that handles the new client arrival / old client arrival on the
root page */
func (b *BaseServer) handleRoot(w http.ResponseWriter, r *http.Request) {
	Request.Printf("%s %s\n", r.Method, r.URL.String())

	path, err := filepath.Abs("static/game.html")
	if err != nil {
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "text/html")
	w.Write(content)
}

func (b *BaseServer) handleConnection(w http.ResponseWriter, r *http.Request) {
	Request.Printf("%s %s\n", r.Method, r.URL.String())
	clientId := b.getClientId(w, r)

	client := GameClient{clientId}

	// new client, not really hack proof tho...
	if clientId >= b.currentId {
		b.clients <- client
	}
}

/* Gets the client id associated to the request,
if None, gets a new one */
func (b *BaseServer) getClientId(w http.ResponseWriter, r *http.Request) int {
	err := r.ParseForm()

	if err != nil {
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// For some reason, it might be a new client, assign a new Id
	if r.Form.Get("clientId") == "" {
		return b.newClient(w)
	}

	i, err := strconv.Atoi(r.Form.Get("clientId"))

	if err != nil {
		Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return -1
	}

	return i
}

func (b *BaseServer) newClient(w http.ResponseWriter) int {
	b.currentId += 1

	w.Write(prot.GetClientIdMessage(b.currentId))

	return b.currentId
}
