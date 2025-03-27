package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"app/servers/internal/servers"

	"github.com/zishang520/engine.io/v2/engine"
	"github.com/zishang520/engine.io/v2/log"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/socket.io/v2/socket"
)

func main() {
	log.DEBUG = true

	e := servers.Engine("127.0.0.1:8000", "server.crt", "server.key")

	e.On("connection", func(sockets ...interface{}) {
		socket := sockets[0].(engine.Socket)
		socket.On("message", func(args ...interface{}) {
			socket.Send(args[0].(io.Reader), nil, nil)
		})
		socket.On("heartbeat", func(...any) {
			utils.Log().Debug("heartbeat %v", socket.Request().Query())
		})
		socket.On("close", func(e ...any) {
			utils.Log().Debug("close %v", e)
		})
	})
	e.On("connection_error", func(e ...any) {
		utils.Log().Debug("connection_error %v", e[0].(*types.ErrorMessage).Context)
	})

	s := servers.Socket("127.0.0.1:3000", "server.crt", "server.key")
	s.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)

		client.On("message", func(args ...interface{}) {
			client.Emit("message-back", args...)
		})
		client.Emit("auth", client.Handshake().Auth)

		client.On("message-with-ack", func(args ...interface{}) {
			ack := args[len(args)-1].(socket.Ack)
			ack(args[:len(args)-1], nil)
		})
		client.OnAny(func(args ...any) {
			client.Emit(args[0].(string), args[1:]...)
		})
	})

	s.Of(regexp.MustCompile(`/\w+`), nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.Emit("auth", client.Handshake().Auth)
		client.OnAny(func(args ...any) {
			client.Emit(args[0].(string), args[1:]...)
		})
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()
	<-ctx.Done()

	e.Close()
	s.Close(nil)
}
