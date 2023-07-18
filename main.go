package main

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
)

func main() {
	dir, _ := os.Getwd()
	// create a new webtransport.Server, listening on (UDP) port 443
	s := webtransport.Server{
		H3: http3.Server{Addr: ":443"},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// Create a new HTTP endpoint /webtransport.
	http.HandleFunc("/webtransport", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.Upgrade(w, r)
		if err != nil {
			log.Printf("upgrading failed: %s", err)
			w.WriteHeader(500)
			return
		}
		log.Printf("conn: %v", conn)
		// Handle the connection. Here goes the application logic.
	})
	s.ListenAndServeTLS(path.Join(dir, "server.crt"), path.Join(dir, "server.key"))
}
