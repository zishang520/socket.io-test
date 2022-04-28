package main

import (
	"github.com/zishang520/engine.io/config"
	"github.com/zishang520/engine.io/engine"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	utils.Log().DEBUG = true

	serverOptions := &config.ServerOptions{}
	serverOptions.SetAllowEIO3(true)
	serverOptions.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})

	http := types.CreateServer(nil).Listen("127.0.0.1:4444", nil)

	engineServer := engine.Attach(http, serverOptions)

	engineServer.On("connection", func(sockets ...interface{}) {
		socket := sockets[0].(engine.Socket)
		socket.On("message", func(...interface{}) {
		})
		socket.On("close", func(...interface{}) {
			utils.Log().Println("client close.")
		})
	})
	utils.Log().Println("%v", engineServer)

	exit := make(chan struct{})
	SignalC := make(chan os.Signal)

	signal.Notify(SignalC, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range SignalC {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
			}
		}
	}()

	<-exit
	os.Exit(0)
}
