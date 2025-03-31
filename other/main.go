package main

import (
	"github.com/zishang520/engine.io-client-go/engine"
	"github.com/zishang520/engine.io-client-go/transports"
	"github.com/zishang520/engine.io/v2/types"
)

func main() {
	opts := engine.DefaultSocketOptions()
	opts.SetPath("/engine.io")
	opts.SetQuery(map[string][]string{
		"token": {"abc123"},
	})
	opts.SetTransports(types.NewSet(transports.WebSocket))

	socket := engine.NewSocket("ws://localhost", opts)

	// ... Handle connection events
}
