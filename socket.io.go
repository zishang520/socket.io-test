package main

import (
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	// "time"

	"github.com/zishang520/engine.io/log"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
)

func main() {
	log.DEBUG = true
	go func() {
		utils.Log().Success("%v", http.ListenAndServe("localhost:6060", nil))
	}()
	c := socket.DefaultServerOptions()
	c.SetAllowEIO3(true)
	c.SetCors(&types.Cors{
		Origin:      "http://127.0.0.1:8000",
		Credentials: true,
	})
	utils.Log().Success("AllowEIO3：%v", c.AllowEIO3())
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
	io.Of(
		regexp.MustCompile(`/\w+`),
		nil,
	).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		utils.Log().Success("/ test Handshake：%v", client.Handshake())
		client.Broadcast().Emit("hi test")
		client.On("event", func(clients ...interface{}) {
			utils.Log().Success("/ test eventeventeventeventevent%v", clients)
		})
		client.On("disconnect", func(...interface{}) {
			utils.Log().Success("/ test disconnect")
		})
		client.On("chat message", func(msgs ...interface{}) {
			io.Of("/test", nil).Emit("hi", msgs...)
			client.Emit("chat message", msgs...)
			client.Emit("chat message", map[string]interface{}{
				"message": types.NewStringBufferString("xxx"),
				"bin":     types.NewBytesBuffer([]byte{0, 1, 2, 3, 4, 5}),
			})
		})
	})
	io.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		utils.Log().Success("Handshake：%v", client.Handshake())
		client.Broadcast().Emit("hi")
		client.On("event", func(clients ...interface{}) {
			utils.Log().Success("eventeventeventeventevent%v", clients)
		})
		client.On("disconnect", func(...interface{}) {
			utils.Log().Success("disconnect")
		})
		client.On("chat message", func(msgs ...interface{}) {
			io.Of("/test", nil).Emit("hi", msgs...)
			client.Emit("chat message", msgs...)
			utils.Log().Success("message：%v", msgs[0])
			utils.Log().Success("FetchSockets %v", io.Adapter().FetchSockets(&socket.BroadcastOptions{
				Rooms: types.NewSet[socket.Room]("/"),
			}))
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
