package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"app/engine"
	"app/socket"

	"github.com/zishang520/engine.io/v2/log"
)

func main() {
	log.DEBUG = true

	engine := engine.EngineServer("127.0.0.1:8000", "server.crt", "server.key")
	socket := socket.SocketServer("127.0.0.1:3000", "server.crt", "server.key")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()
	<-ctx.Done()

	engine.Close()
	socket.Close(nil)
}
