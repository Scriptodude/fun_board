package checker

import (
	"server/core"
)

var (
	blackPlayer core.GameClient
	whitePlayer core.GameClient
	coreServer  core.Server
)

type CheckerServer struct{}

func (CheckerServer) Init(server core.Server) {
	coreServer = server
	core.Info.Println("The checker server was initialized.")

	err := server.ServeAndListen()

	if err != nil {
		core.Error.Println(err)
	}
}

func (CheckerServer) Shutdown() {
	core.Info.Println("Shutting down the checker server...")

	// For now we shutdown the core server as well, eventually we might want
	// to change the game type without restarting the server..
	coreServer.Shutdown()
}
