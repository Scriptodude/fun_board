package core

import "net/http"

/* Represent the basic server that manages clients and sockets and handlers */
type Server interface {
	ServeAndListen() error
	GetAddress() (string, error)
	GetPort() (string, error)
	AwaitClient() GameClient
	Shutdown()
	AddRequestListener(path string, fn func(w http.ResponseWriter, r *http.Request, client GameClient))
}

/* Represents the resources of the game server which manages:
- The requests
- The validation
- Adds the required handlers

All the main requests should be handled by the GameServer via a call to AddRequestListener
By default only / is handled by the server, which returns the static/game.html if found */
type GameServer interface {
	Init(server Server) // Inits the server's resources, if applicable.
	Shutdown()
}

/* Represents a client that connected to the GameServer.*/
type GameClient struct {
	Id int
}
