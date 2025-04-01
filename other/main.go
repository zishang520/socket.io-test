package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/zishang520/engine.io-client-go/engine"
	"github.com/zishang520/engine.io-client-go/transports"
	"github.com/zishang520/engine.io/v2/log"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
)

func main() {
	log.DEBUG = true
	opts := engine.DefaultSocketOptions()
	opts.SetTransports(types.NewSet(transports.Polling /*transports.WebSocket, transports.WebTransport*/))

	e := engine.NewSocket("http://127.0.0.1:4444", opts)
	e.On("open", func(args ...any) {
		e.Send(types.NewStringBufferString("88888"), nil, nil)
		utils.Log().Debug("close %v", args)
	})

	e.On("close", func(args ...any) {
		utils.Log().Debug("close %v", args)
	})

	e.On("packet", func(args ...any) {
		utils.Log().Warning("packet: %+v", args)
	})

	e.On("ping", func(...any) {
		utils.Log().Warning("ping")
	})

	e.On("pong", func(...any) {
		utils.Log().Warning("pong")
	})

	e.On("message", func(args ...any) {
		e.Send(types.NewStringBufferString("6666666"), nil, nil)
		utils.Log().Warning("message %v", args)
	})

	e.On("heartbeat", func(...any) {
		utils.Log().Debug("heartbeat")
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()
	<-ctx.Done()

	e.Close()
}
