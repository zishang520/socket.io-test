package main

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path"
	"syscall"

	// "time"

	"github.com/gorilla/websocket"
	_types "github.com/zishang520/engine.io-go-parser/types"
	"github.com/zishang520/engine.io/config"
	"github.com/zishang520/engine.io/engine"
	"github.com/zishang520/engine.io/log"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
)

type Test struct{ B int }

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World")
}
func main() {
	log.DEBUG = true

	serverOptions := &config.ServerOptions{}
	serverOptions.SetAllowEIO3(true)
	serverOptions.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	serverOptions.SetTransports(types.NewSet("polling", "websocket", "webtransport"))

	xxx := types.NewSet("1", "2", "3", "5")
	utils.Log().Debug("%v", xxx)
	cache := xxx.All()
	utils.Log().Debug("%v", cache)
	delete(cache, "1")
	utils.Log().Debug("%v", cache)
	utils.Log().Debug("%v", xxx)

	go http.ListenAndServe("127.0.0.1:6060", nil)

	dir, _ := os.Getwd()
	httpServer := types.NewWebServer(nil)
	// utils.SetTimeOut(func() {
	// 	httpServer.Close(nil)
	// }, 10000*time.Millisecond)

	// utils.SetTimeOut(func() {
	// 	httpServer.Close(nil)
	// }, 12000*time.Millisecond)

	httpServer.HandleFunc("/public/", func(w http.ResponseWriter, r *http.Request) {
		file, err := http.Dir(dir).Open(path.Clean("/" + r.URL.Path))
		if err != nil {
			http.Error(w, "file not found:"+path.Clean("/"+r.URL.Path), http.StatusNotFound)
			return
		}
		io.Copy(w, file)
	})
	httpServer.ListenTLS(":443", "server.crt", "server.key", nil)
	wts := httpServer.ListenWebTransportTLS(":443", "server.crt", "server.key", nil, nil)

	engineServer := engine.New(serverOptions)

	httpServer.HandleFunc("/engine.io/", func(w http.ResponseWriter, r *http.Request) {
		if !websocket.IsWebSocketUpgrade(r) {
			engineServer.OnWebTransportSession(types.NewHttpContext(w, r), wts)
			// engineServer.HandleRequest(types.NewHttpContext(w, r))
		} else if engineServer.Opts().Transports().Has("websocket") {
			engineServer.HandleUpgrade(types.NewHttpContext(w, r))
		} else {
			// httpServer.DefaultHandler.ServeHTTP(w, r)
			engineServer.OnWebTransportSession(types.NewHttpContext(w, r), wts)
		}
	})

	engineServer.On("connection", func(sockets ...interface{}) {
		socket := sockets[0].(engine.Socket)
		socket.On("message", func(args ...interface{}) {
			// socket.Send(_types.NewBytesBufferString("66666666"), nil, nil)
			socket.Send(_types.NewStringBufferString("66666666666"), nil, nil)
			// utils.Log().Debug("%v", socket.Protocol())
			// utils.Log().Debug("%v", socket.Id())
			// utils.Log().Debug("%v", socket.Request().Headers())
			// utils.Log().Debug("%v", socket.Request().Query())
			// utils.Log().Debug("'%v'", socket.Request().Request().Body)
		})
		socket.On("heartbeat", func(...any) {
			utils.Log().Debug("heartbeat %v", socket.Request().Query())
		})
		socket.On("close", func(e ...any) {
			utils.Log().Debug("close %v", e)
		})
	})
	engineServer.On("connection_error", func(e ...any) {
		utils.Log().Debug("connection_error %v", e[0].(*types.ErrorMessage).Context)
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
