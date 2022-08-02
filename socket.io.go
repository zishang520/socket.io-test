package main

import (
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

func main() {
	utils.Log().DEBUG = true
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
	httpServer.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
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
		utils.Log().Success("/ tets Handshake：%v", client.Handshake())
		client.Broadcast().Emit("hi tets")
		client.On("event", func(clients ...interface{}) {
			utils.Log().Success("/ tets eventeventeventeventevent%v", clients)
		})
		client.On("disconnect", func(...interface{}) {
			utils.Log().Success("/ tets disconnect")
		})
		client.On("chat message", func(msgs ...interface{}) {
			io.Of("/test", nil).Emit("hi", msgs...)
			client.Timeout(2000*time.Millisecond).Emit("chat message", msgs[0], func(err error, args ...interface{}) {
				utils.Log().Error("OUT %v %v", err, args)
			})
			client.To("xxx")
			client.Emit("chat message", msgs...)
			client.Emit("chat message", map[string]interface{}{
				"message": types.NewStringBufferString("xxx"),
				"bin":     types.NewBytesBuffer([]byte{0, 1, 2, 3, 4, 5}),
			})
			utils.Log().Success("/ tets message：%v", msgs[0])
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
	httpServer.Close(nil)
	os.Exit(0)
	// ad := socket.NewAdapter(nil)
	// ad.AddAll("1", types.NewSet("1", "2"))
	// a := socket.NewBroadcastOperator(nil, nil, nil, nil)
	// utils.Log().Info("%v", a.Compress(true))
	// utils.Log().Info("%v", a)
}
