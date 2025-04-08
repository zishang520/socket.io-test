package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "app/clients/internal"
	"app/clients/internal/clients"

	"github.com/zishang520/engine.io-client-go/engine"
	"github.com/zishang520/engine.io-client-go/transports"
	"github.com/zishang520/engine.io/v2/log"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/socket.io-client-go/socket"
)

func e() {
	certPEM, err := os.ReadFile("root.crt")
	if err != nil {
		utils.Log().Fatalf("读取证书失败: %v", err)
	}

	rootCAs := x509.NewCertPool()
	ok := rootCAs.AppendCertsFromPEM(certPEM)
	if !ok {
		utils.Log().Fatal("添加自签名证书失败")
	}

	opts := engine.DefaultSocketOptions()
	opts.SetTLSClientConfig(&tls.Config{
		RootCAs:   rootCAs,
		ClientCAs: rootCAs,
	})
	opts.SetTransports(types.NewSet( /*transports.Polling, */ /*transports.WebSocket, */ transports.WebTransport))

	e := clients.Engine("https://127.0.0.1:8000", opts)
	e.On("open", func(args ...any) {
		utils.SetTimeout(func() {
			e.Send(types.NewStringBufferString("88888"), nil, nil)
			e.Send(types.NewStringBufferString("88888"), nil, nil)
		}, 1*time.Second)
		utils.Log().Debug("close %v", args)
	})

	e.On("close", func(args ...any) {
		utils.Log().Debug("close %v", args)
	})

	e.On("packet", func(args ...any) {
		utils.Log().Warning("packet: %+v", args)
	})

	e.On("ping", func(...any) {
		utils.Log().Warning("ping")
	})

	e.On("pong", func(...any) {
		utils.Log().Warning("pong")
	})

	e.On("message", func(args ...any) {
		e.Send(types.NewStringBufferString("6666666"), nil, nil)
		utils.Log().Warning("message %v", args)
	})

	e.On("heartbeat", func(...any) {
		utils.Log().Debug("heartbeat")
	})
}

func s() {
	certPEM, err := os.ReadFile("root.crt")
	if err != nil {
		utils.Log().Fatalf("读取证书失败: %v", err)
	}

	rootCAs := x509.NewCertPool()
	ok := rootCAs.AppendCertsFromPEM(certPEM)
	if !ok {
		utils.Log().Fatal("添加自签名证书失败")
	}

	opts := socket.DefaultOptions()
	opts.SetTLSClientConfig(&tls.Config{
		RootCAs:   rootCAs,
		ClientCAs: rootCAs,
	})
	opts.SetTransports(types.NewSet( /*transports.Polling, */ /*transports.WebSocket, */ transports.WebTransport))

	e, err := clients.Socket("https://127.0.0.1:8000", opts)
	utils.Log().Error("socket %v", e)
	if err != nil {
		utils.Log().Fatal("exit %v", err)
		return
	}
	e.On("open", func(args ...any) {
		utils.SetTimeout(func() {
			e.Emit("test", types.NewStringBufferString("88888"))
		}, 1*time.Second)
		utils.Log().Debug("close %v", args)
	})

	e.On("close", func(args ...any) {
		utils.Log().Debug("close %v", args)
	})

	e.On("packet", func(args ...any) {
		utils.Log().Warning("packet: %+v", args)
	})

	e.On("ping", func(...any) {
		utils.Log().Warning("ping")
	})

	e.On("pong", func(...any) {
		utils.Log().Warning("pong")
	})

	e.On("message", func(args ...any) {
		e.Send(types.NewStringBufferString("6666666"), nil, nil)
		utils.Log().Warning("message %v", args)
	})

	e.On("heartbeat", func(...any) {
		utils.Log().Debug("heartbeat")
	})
}

func main() {
	log.DEBUG = true

	s()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()
	<-ctx.Done()
}
