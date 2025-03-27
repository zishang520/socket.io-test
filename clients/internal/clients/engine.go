package clients

import "github.com/zishang520/engine.io-client-go/engine"

func Engine(uri string, opts engine.SocketOptionsInterface) engine.Socket {
	socket := engine.NewSocket(uri, opts)
	return socket
}
