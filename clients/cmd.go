package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"os/signal"
	"syscall"

	_ "app/clients/internal"
	"app/clients/internal/clients"

	"github.com/zishang520/engine.io-client-go/engine"
	"github.com/zishang520/engine.io-client-go/transports"
	"github.com/zishang520/engine.io/v2/log"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
)

func main() {
	log.DEBUG = true

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
	opts.SetTransports(types.NewSet(transports.Polling /*, transports.WebSocket, transports.WebTransport*/))

	e := clients.Engine("https://127.0.0.1:8000", opts)
	e.On("open", func(args ...any) {
		utils.Log().Debug("close %v", args)
		e.Send(types.NewStringBufferString("ping"), nil, nil)
	})

	e.On("close", func(args ...any) {
		utils.Log().Debug("close %v", args)
	})

	e.On("packet", func(args ...any) {
		utils.Log().Info("packet: %+v", args)
	})

	e.On("ping", func(...any) {
		utils.Log().Info("ping")
	})

	e.On("pong", func(...any) {
		utils.Log().Info("pong")
	})

	e.On("message", func(args ...any) {
		utils.Log().Info("close %v", args)
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()
	<-ctx.Done()

	e.Close()
}
