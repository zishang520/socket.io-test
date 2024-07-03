package main

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zishang520/engine.io/v2/log"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/socket.io/v2/socket"
)

func main() {
	log.DEBUG = true
	c := socket.DefaultServerOptions()
	c.SetAllowEIO3(true)
	// c.SetConnectionStateRecovery(&socket.ConnectionStateRecovery{})
	// c.SetAllowEIO3(true)
	c.SetPingInterval(300 * time.Millisecond)
	c.SetPingTimeout(200 * time.Millisecond)
	c.SetMaxHttpBufferSize(1000000)
	c.SetConnectTimeout(1000 * time.Millisecond)
	c.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	httpServer := types.CreateServer(nil)
	dir, _ := os.Getwd()
	httpServer.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		file, err := http.Dir(dir).Open("index.html")
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		io.Copy(w, file)
	})
	io := socket.NewServer(httpServer, c)
	io.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.Emit("auth", client.Handshake().Auth)

		client.On("message", func(args ...interface{}) {
			client.Emit("message-back", args...)
		})

		client.On("message-with-ack", func(args ...interface{}) {
			ack := args[len(args)-1].(func([]any, error))
			ack(args[:len(args)-1], nil)
		})
	})

	io.Of("/custom", nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.Emit("auth", client.Handshake().Auth)
	})

	httpServer.Listen(":3000", nil)

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
}
