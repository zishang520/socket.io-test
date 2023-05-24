module test

go 1.18

replace github.com/zishang520/socket.io => ../socket.io/

replace github.com/zishang520/engine.io => ../engine.io/

require github.com/zishang520/engine.io v1.4.4

require github.com/zishang520/socket.io v1.0.4

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/gookit/color v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
