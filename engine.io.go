package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"syscall"
	// "time"

	"github.com/zishang520/engine.io/config"
	"github.com/zishang520/engine.io/engine"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
)

type Test struct{ B int }

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
func main() {
	// utils.Log().DEBUG = true

	serverOptions := &config.ServerOptions{}
	serverOptions.SetAllowEIO3(true)
	serverOptions.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	xxx := types.NewSet("1", "2", "3", "5")
	utils.Log().Debug("%v", xxx)
	cache := xxx.All()
	utils.Log().Debug("%v", cache)
	delete(cache, "1")
	utils.Log().Debug("%v", cache)
	utils.Log().Debug("%v", xxx)

	httpServer := types.CreateServer(nil).Listen("127.0.0.1:4444", nil)

	// utils.SetTimeOut(func() {
	// 	httpServer.Close(nil)
	// }, 10000*time.Millisecond)

	// utils.SetTimeOut(func() {
	// 	httpServer.Close(nil)
	// }, 12000*time.Millisecond)

	dir, _ := os.Getwd()
	httpServer.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		file, err := http.Dir(dir).Open(path.Clean("/" + r.URL.Path))
		if err != nil {
			http.Error(w, "file not found:"+path.Clean("/"+r.URL.Path), http.StatusNotFound)
			return
		}
		io.Copy(w, file)
	})
	engineServer := engine.New(serverOptions)

	engineServer.Attach(httpServer, nil)

	httpServer.On("close", func(...any) {
		engineServer.Close()
	})

	engineServer.On("connection", func(sockets ...interface{}) {
		socket := sockets[0].(engine.Socket)
		socket.On("message", func(args ...interface{}) {
			socket.Send(types.NewStringBufferString("xxx"), nil, nil)
			utils.Log().Debug("%v", socket.Protocol())
			utils.Log().Debug("%v", socket.Id())
			utils.Log().Debug("%v", socket.Request().Headers())
			utils.Log().Debug("%v", socket.Request().Query())
			utils.Log().Debug("'%v'", socket.Request().Request().Body)
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
				exit <- struct{}{}
			}
		}
	}()

	<-exit
	httpServer.Close(nil)
	os.Exit(0)
}
