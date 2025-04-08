package clients

import "github.com/zishang520/socket.io-client-go/socket"

func Socket(uri string, opts socket.OptionsInterface) (*socket.Socket, error) {
	return socket.Io(uri, opts)
}
