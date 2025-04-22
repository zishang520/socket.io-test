package clients

import "github.com/zishang520/socket.io/clients/engine/v3"

func Engine(uri string, opts engine.SocketOptionsInterface) engine.Socket {
	return engine.NewSocket(uri, opts)
}
