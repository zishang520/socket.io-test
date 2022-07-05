package main

import (
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	utils.Log().DEBUG = true
	c := socket.DefaultServerOptions()
	c.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	httpServer := types.CreateServer(nil)
	io := socket.NewServer(httpServer, c)
	io.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.On("event", func(clients ...interface{}) {
			utils.Log().Success("eventeventeventeventevent%v", clients)
		})
		client.On("disconnect", func(...interface{}) {
			utils.Log().Success("disconnect")
		})
		client.On("chat message", func(msgs ...interface{}) {
			utils.Log().Success("消息：%v", msgs)
		})
	})
	httpServer.Listen("127.0.0.1:9999", nil)
	exit := make(chan struct{})
	SignalC := make(chan os.Signal)

	signal.Notify(SignalC, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				exit <- struct{}{}
			}
		}
	}()

	<-exit
	httpServer.Close(nil)
	os.Exit(0)
	// ad := socket.NewAdapter(nil)
	// ad.AddAll("1", types.NewSet("1", "2"))
	// a := socket.NewBroadcastOperator(nil, nil, nil, nil)
	// utils.Log().Info("%v", a.Compress(true))
	// utils.Log().Info("%v", a)
}
