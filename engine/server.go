package engine

import (
	"io"
	"net/http"

	"github.com/zishang520/engine.io/v2/config"
	"github.com/zishang520/engine.io/v2/engine"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/engine.io/v2/webtransport"
)

func EngineServer(addr string, certFile string, keyFile string) engine.Server {
	serverOptions := &config.ServerOptions{}
	serverOptions.SetAllowEIO3(true)
	serverOptions.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	serverOptions.SetTransports(types.NewSet("polling", "websocket", "webtransport"))

	httpServer := types.NewWebServer(nil)
	httpServer.ListenTLS(addr, certFile, keyFile, nil)

	engineServer := engine.New(httpServer, serverOptions)

	server := types.NewWebServer(nil)
	wts := server.ListenWebTransportTLS(addr, certFile, keyFile, nil, nil)
	server.HandleFunc("/engine.io/", func(w http.ResponseWriter, r *http.Request) {
		if webtransport.IsWebTransportUpgrade(r) {
			engineServer.OnWebTransportSession(types.NewHttpContext(w, r), wts)
		} else {
			server.DefaultHandler.ServeHTTP(w, r)
		}
	})

	engineServer.On("connection", func(sockets ...interface{}) {
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
	engineServer.On("connection_error", func(e ...any) {
		utils.Log().Debug("connection_error %v", e[0].(*types.ErrorMessage).Context)
	})

	return engineServer
}
