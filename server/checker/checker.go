package checker

import (
	"server/core"
)

var (
	blackPlayer *core.GameClient
	whitePlayer *core.GameClient
)

type CheckerServer struct{}

func (CheckerServer) Init(server core.Server) error {
	core.Info.Println("The checker server was initialized.")

	return nil
}

func (CheckerServer) Shutdown() error {
	core.Info.Println("The checker server was shutdown")

	return nil
}
