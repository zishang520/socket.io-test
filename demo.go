package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/v2/socket"
)

func main() {
	io := socket.NewServer(nil, nil)
	io.On("connection", func(clients ...any) {
		socket := clients[0].(*socket.Socket)

		utils.Log().Info(`socket %s connected`, socket.Id())

		// send an event to the client
		socket.Emit("foo", "bar")

		socket.On("foobar", func(...any) {
			// an event was received from the client
		})

		// upon disconnection
		socket.On("disconnect", func(reason ...any) {
			utils.Log().Info(`socket %s disconnected due to %s`, socket.Id(), reason[0])
		})
	})
	io.Listen(3000, nil)

	exit := make(chan struct{})
	SignalC := make(chan os.Signal)

	signal.Notify(SignalC, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	io.Close(nil)
	os.Exit(0)
	// ad := socket.NewAdapter(nil)
	// ad.AddAll("1", types.NewSet("1", "2"))
	// a := socket.NewBroadcastOperator(nil, nil, nil, nil)
	// utils.Log().Info("%v", a.Compress(true))
	// utils.Log().Info("%v", a)
}
