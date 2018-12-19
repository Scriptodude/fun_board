package main

import (
	"os"
	"os/signal"
	"server/checker"
	"server/core"
	"syscall"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	core.InitLoggers()
	base := core.NewDefaultServer()
	game := &checker.CheckerServer{}

	go func() {
		v := <-sigs
		core.Info.Println(v)

		game.Shutdown()
	}()

	game.Init(base)
}
