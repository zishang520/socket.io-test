module test

go 1.18

replace github.com/zishang520/socket.io => ../socket.io/

replace github.com/zishang520/engine.io => ../engine.io/

require github.com/zishang520/engine.io v1.1.7

require github.com/zishang520/socket.io v1.0.4

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/google/brotli/go/cbrotli v0.0.0-20220512075048-9801a2c5d6c6 // indirect
	github.com/gookit/color v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44 // indirect
)
