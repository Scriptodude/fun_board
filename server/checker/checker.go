package checker

import (
	"server/core"
	i "server/interfaces"
)

var (
	blackPlayer *i.GameClient
	whitePlayer *i.GameClient
	coreServer  i.Server
)

type CheckerServer struct{}

func (CheckerServer) Init(server i.Server) {
	coreServer = server
	core.Info.Println("The checker server was initialized.")

	whitePlayer = server.AwaitClient()
	blackPlayer = server.AwaitClient()

	core.Info.Printf(`There are two players connected,
	 				 starting the game with %+v and %+v`, whitePlayer, blackPlayer)
}

func (CheckerServer) Shutdown() {
	core.Info.Println("Shutting down the checker server...")

	// For now we shutdown the core server as well, eventually we might want
	// to change the game type without restarting the server..
	coreServer.Shutdown()
}
