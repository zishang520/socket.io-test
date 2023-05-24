module test

go 1.18

replace github.com/zishang520/socket.io => ../socket.io/

replace github.com/zishang520/engine.io => ../engine.io/

replace github.com/zishang520/engine.io-server-go-fasthttp => ../engine.io-server-go-fasthttp/

require github.com/zishang520/engine.io v1.4.4

require (
	github.com/zishang520/engine.io-server-go-fasthttp v0.0.0-00010101000000-000000000000
	github.com/zishang520/socket.io v1.0.4
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/fasthttp/websocket v1.5.3 // indirect
	github.com/gookit/color v1.5.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.47.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)
