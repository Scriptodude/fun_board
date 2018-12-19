package main

import (
	"server/checker"
	"server/core"
)

func main() {
	core.InitLoggers()
	base := core.NewDefaultServer()
	game := &checker.CheckerServer{}

	game.Init(base)
}
