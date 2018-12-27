package checker

import (
	"server/core"
	i "server/interfaces"
	"time"
)

var (
	blackPlayer *i.GameClient
	whitePlayer *i.GameClient
	coreServer  i.Server
	spectators  []*i.GameClient
)

type CheckerServer struct{}

func (CheckerServer) Init(server i.Server) {
	coreServer = server
	core.Info.Println("The checker server was initialized.")

	whitePlayer = server.AwaitClient()
	blackPlayer = server.AwaitClient()

	core.Info.Printf(`There are two players connected,
	 				 starting the game with %+v and %+v`, whitePlayer, blackPlayer)

	go func() {
		if len(spectators) < 10 {
			spectators = append(spectators, server.AwaitClient())
		}
	}()
}

func (CheckerServer) Shutdown() {
	core.Info.Println("Shutting down the checker server...")
	// For now we shutdown the core server as well, eventually we might want
	// to change the game type without restarting the server..
	coreServer.Shutdown()
}

func (CheckerServer) Play() {
	for {
		select {
		case <-time.After(time.Second * 10):
		case <-coreServer.GetContext().Done():
			core.Info.Println("Leaving the game")
			return
		}
		core.Info.Printf("Playing, we have %d spectators", len(spectators))
	}
}
