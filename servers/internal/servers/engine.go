package servers

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/zishang520/socket.io/servers/engine/v3"
	"github.com/zishang520/socket.io/servers/engine/v3/config"
	"github.com/zishang520/socket.io/servers/engine/v3/transports"
	"github.com/zishang520/socket.io/v3/pkg/types"
	"github.com/zishang520/socket.io/v3/pkg/webtransport"
)

func Engine(addr string, certFile string, keyFile string) engine.Server {
	serverOptions := &config.ServerOptions{}
	serverOptions.SetAllowEIO3(true)
	serverOptions.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	// serverOptions.SetPingInterval(120 * time.Second)
	// serverOptions.SetPingTimeout(100 * time.Second)
	serverOptions.SetMaxHttpBufferSize(1000000)
	serverOptions.SetTransports(types.NewSet( /*engine.Polling, engine.WebSocket,*/ engine.WebTransport))

	httpServer := types.NewWebServer(nil)
	// httpServer.ListenHTTP3TLS(addr, certFile, keyFile, nil, nil)

	engineServer := engine.New(httpServer, serverOptions)

	wts := httpServer.ListenWebTransportTLS(addr, certFile, keyFile, nil, nil)
	httpServer.HandleFunc("/engine.io/", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade http3
		wts.H3.SetQUICHeaders(w.Header())
		if webtransport.IsWebTransportUpgrade(r) {
			engineServer.OnWebTransportSession(types.NewHttpContext(w, r), wts)
		} else if !websocket.IsWebSocketUpgrade(r) {
			engineServer.HandleRequest(types.NewHttpContext(w, r))
		} else if engineServer.Transports().Has(transports.WEBSOCKET) {
			engineServer.HandleUpgrade(types.NewHttpContext(w, r))
		} else {
			httpServer.DefaultHandler.ServeHTTP(w, r)
		}
	})

	httpServer.ListenTLS(addr, certFile, keyFile, nil)

	return engineServer
}
