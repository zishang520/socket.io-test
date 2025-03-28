package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/quic-go/quic-go/http3"
	"github.com/zishang520/engine.io-go-parser/types"
	"github.com/zishang520/engine.io/v2/utils"
)

// 自定义 Transport，支持多协议
type MultiTransport struct {
	h3Transport  *http3.Transport  // HTTP/3 的 Transport
	stdTransport *http.Transport   // HTTP/2/1.1 的 Transport
	altSvcCache  map[string]string // 缓存服务端支持的协议（例如 Alt-Svc）
	cacheMutex   sync.RWMutex      // 缓存读写锁
}

// 初始化 MultiTransport
func NewMultiTransport() *MultiTransport {

	certPEM, err := os.ReadFile("root.crt")
	if err != nil {
		utils.Log().Fatalf("读取证书失败: %v", err)
	}

	rootCAs := x509.NewCertPool()
	ok := rootCAs.AppendCertsFromPEM(certPEM)
	if !ok {
		utils.Log().Fatal("添加自签名证书失败")
	}
	return &MultiTransport{
		h3Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:   rootCAs,
				ClientCAs: rootCAs,
				// NextProtos: []string{"h3"}, // QUIC 必须的 ALPN 标识
			},
		},
		stdTransport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:   rootCAs,
				ClientCAs: rootCAs,
				// NextProtos: []string{"h2", "http/1.1"}, // 支持 HTTP/2 和 HTTP/1.1
			},
		},
		altSvcCache: make(map[string]string),
	}
}

// 实现 RoundTrip 接口
func (t *MultiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 1. 检查是否有缓存的 Alt-Svc 信息（是否已知服务端支持 HTTP/3）
	t.cacheMutex.RLock()
	altSvc, ok := t.altSvcCache[req.URL.Host]
	t.cacheMutex.RUnlock()

	if ok && altSvc == "h3" {
		// 2. 如果已知支持 HTTP/3，直接使用 QUIC
		return t.h3Transport.RoundTrip(req)
	}

	// 3. 尝试 HTTP/3
	resp, err := t.h3Transport.RoundTrip(req)
	if err == nil {
		// 4. 如果成功，缓存 Alt-Svc 信息
		t.cacheMutex.Lock()
		t.altSvcCache[req.URL.Host] = "h3"
		t.cacheMutex.Unlock()
		fmt.Println("Http3")
		return resp, nil
	}
	fmt.Println("Http2")

	// 5. 如果 HTTP/3 失败，回退到 HTTP/2 或 HTTP/1.1
	return t.stdTransport.RoundTrip(req)
}

// 关闭 Transport 资源
func (t *MultiTransport) Close() error {
	_ = t.h3Transport.Close()             // 关闭 HTTP/3 连接
	t.stdTransport.CloseIdleConnections() // 关闭 HTTP/2/1.1 连接
	return nil
}

// 使用示例
func main() {
	client := &http.Client{
		Transport: NewMultiTransport(),
	}
	resp, err := client.Get("https://127.0.0.1:8000")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, _ := types.NewBytesBufferReader(resp.Body)
	fmt.Println(data.String())
}
