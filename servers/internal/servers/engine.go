package servers

import (
	"net/http"

	"github.com/zishang520/engine.io/v2/config"
	"github.com/zishang520/engine.io/v2/engine"
	"github.com/zishang520/engine.io/v2/transports"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/webtransport"
)

func Engine(addr string, certFile string, keyFile string) engine.Server {
	serverOptions := &config.ServerOptions{}
	serverOptions.SetAllowEIO3(true)
	serverOptions.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	serverOptions.SetTransports(types.NewSet(transports.POLLING, transports.WEBSOCKET, transports.WEBTRANSPORT))

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

	return engineServer
}
