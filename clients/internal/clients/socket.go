package clients

import "github.com/zishang520/socket.io/clients/socket/v3"

func Socket(uri string, opts socket.OptionsInterface) (*socket.Socket, error) {
	return socket.Io(uri, opts)
}
