package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/quic-go/quic-go/http3"
)

type AltSvcInfo struct {
	Protocol string
	Port     string
	Expiry   time.Time
}

type MultiTransport struct {
	h3Transport  *http3.Transport
	stdTransport *http.Transport
	altSvcCache  sync.Map
	re           *regexp.Regexp
}

func NewMultiTransport() *MultiTransport {
	certPEM, err := os.ReadFile("root.crt")
	if err != nil {
		panic(fmt.Sprintf("读取证书失败: %v", err))
	}

	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(certPEM) {
		panic("添加自签名证书失败")
	}

	tlsConfig := &tls.Config{RootCAs: rootCAs}

	return &MultiTransport{
		h3Transport:  &http3.Transport{TLSClientConfig: tlsConfig},
		stdTransport: &http.Transport{TLSClientConfig: tlsConfig},
		re:           regexp.MustCompile(`(\w+)="([^"]+)"(?:;\s*ma=(\d+))?`),
	}
}

func (t *MultiTransport) parseAltSvcHeader(header string) []AltSvcInfo {
	var altSvcList []AltSvcInfo
	for _, match := range t.re.FindAllStringSubmatch(header, -1) {
		maxAge := 86400
		if match[3] != "" {
			if ma, err := strconv.Atoi(match[3]); err == nil {
				maxAge = ma
			}
		}
		altSvcList = append(altSvcList, AltSvcInfo{
			Protocol: match[1],
			Port:     match[2],
			Expiry:   time.Now().Add(time.Duration(maxAge) * time.Second),
		})
	}
	return altSvcList
}

func (t *MultiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if altSvc, ok := t.altSvcCache.Load(req.URL.Host); ok {
		if as := altSvc.(AltSvcInfo); as.Protocol == "h3" && time.Now().Before(as.Expiry) {
			req.URL.Scheme = "https"
			req.URL.Host = req.URL.Hostname() + ":" + as.Port
			if resp, err := t.h3Transport.RoundTrip(req); err == nil {
				return resp, nil
			}
		}
	}

	resp, err := t.stdTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if altSvcHeader := resp.Header.Get("Alt-Svc"); altSvcHeader != "" {
		for _, svc := range t.parseAltSvcHeader(altSvcHeader) {
			if svc.Protocol == "h3" {
				fmt.Printf("检测到 Alt-Svc: h3=%s, 有效期: %s\n", svc.Port, svc.Expiry)
				t.altSvcCache.Store(req.URL.Host, svc)
				break
			}
		}
	}

	return resp, nil
}

func (t *MultiTransport) Close() error {
	t.h3Transport.Close()
	t.stdTransport.CloseIdleConnections()
	return nil
}

func main() {
	client := &http.Client{
		Transport: NewMultiTransport(),
		Timeout:   10 * time.Second,
	}

	resp, err := client.Get("https://127.0.0.1:8000")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("响应:", string(body))
}
