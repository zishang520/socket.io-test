package socket

import (
	"net/http"
	"regexp"

	"github.com/zishang520/engine.io/v2/engine"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/webtransport"
	"github.com/zishang520/socket.io/v2/socket"
)

func SocketServer(addr string, certFile string, keyFile string) *socket.Server {
	c := socket.DefaultServerOptions()
	c.SetServeClient(true)
	c.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: true,
	})
	httpServer := types.NewWebServer(nil)
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

	socketio.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)

		client.On("message", func(args ...interface{}) {
			client.Emit("message-back", args...)
		})
		client.Emit("auth", client.Handshake().Auth)

		client.On("message-with-ack", func(args ...interface{}) {
			ack := args[len(args)-1].(socket.Ack)
			ack(args[:len(args)-1], nil)
		})
		client.OnAny(func(args ...any) {
			client.Emit(args[0].(string), args[1:]...)
		})
	})

	socketio.Of(regexp.MustCompile(`/\w+`), nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.Emit("auth", client.Handshake().Auth)
		client.OnAny(func(args ...any) {
			client.Emit(args[0].(string), args[1:]...)
		})
	})

	httpServer.ListenTLS(":3000", "server.crt", "server.key", nil)

	return socketio
}
