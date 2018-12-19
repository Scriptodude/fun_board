package core

import (
	"net/http"
	"time"
)

type BaseServer struct {
	server *http.Server
	mux    *http.ServeMux
}

func NewDefaultServer() *BaseServer {
	mux := http.NewServeMux()

	return &BaseServer{
		server: &http.Server{
			Addr:         ":8080",
			Handler:      mux,
			ReadTimeout:  60 * time.Second,
			WriteTimeout: 60 * time.Second,
		},
		mux: mux,
	}
}

func (b *BaseServer) ServeAndListen(address string, port int) error {
	return nil
}

func (b *BaseServer) GetAddress() (string, error) {
	return "", nil
}

func (b *BaseServer) GetPort() (string, error) {
	return "", nil
}

func (b *BaseServer) AwaitClient() (GameClient, error) {
	return nil, nil
}

func (b *BaseServer) Shutdown() {

}

func (b *BaseServer) AddRequestListener(
	path string,
	fn func(w http.ResponseWriter, r *http.Request, client GameClient)) {

}
