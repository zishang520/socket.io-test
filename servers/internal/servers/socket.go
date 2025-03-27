package servers

import (
	"net/http"

	"github.com/zishang520/engine.io/v2/engine"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/webtransport"
	"github.com/zishang520/socket.io/v2/socket"
)

func Socket(addr string, certFile string, keyFile string) *socket.Server {
	c := socket.DefaultServerOptions()
	c.SetServeClient(true)
	c.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	httpServer := types.NewWebServer(nil)
	httpServer.ListenTLS(addr, certFile, keyFile, nil)

	socketio := socket.NewServer(httpServer, nil)

	// WebTransport start
	// WebTransport uses udp, so you need to enable the new service.
	customServer := types.NewWebServer(nil)
	// A certificate is required and cannot be a self-signed certificate.
	wts := customServer.ListenWebTransportTLS(addr, certFile, keyFile, nil, nil)

	// Here is the core logic of the WebTransport handshake.
	customServer.HandleFunc(socketio.Path()+"/", func(w http.ResponseWriter, r *http.Request) {
		if webtransport.IsWebTransportUpgrade(r) {
			// You need to call socketio.ServeHandler(nil) before this, otherwise you cannot get the Engine instance.
			socketio.Engine().(engine.Server).OnWebTransportSession(types.NewHttpContext(w, r), wts)
		} else {
			customServer.DefaultHandler.ServeHTTP(w, r)
		}
	})
	// WebTransport end

	return socketio
}
