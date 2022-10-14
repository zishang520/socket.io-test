package main

import (
	"fmt"
	// "time"
	// "github.com/zishang520/engine.io/utils"
)

type Test struct{ B int }

func aaaa(s *string) {
	*s = "1"
}
func index(a int) (b int, s string) {
	aaaa(&s)
	b = a
	return a, s
}
func main() {
	a, s := index(0)
	fmt.Println(a, s)
	s = "9"
	fmt.Println(index(1))
	fmt.Println(a, s)
	// httpServer := types.CreateServer(nil).Listen("127.0.0.1:3000", nil)

	// httpServer.HandleFunc("/engine.io", func(w http.ResponseWriter, r *http.Request) {
	// 	ctx := types.NewHttpContext(w, r)
	// 	ctx.On("close", func(...any) {
	// 		fmt.Println("connection closed")
	// 	})
	// 	// ctx.Write(nil)
	// 	utils.SetTimeOut(func() {
	// 		if ctx != nil {
	// 			if h, ok := ctx.Response().(http.Hijacker); ok {
	// 				if netConn, _, err := h.Hijack(); err == nil {
	// 					if netConn.Close() == nil && !ctx.IsDone() {
	// 						ctx.Flush()
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}, 2000*time.Millisecond)
	// 	<-ctx.Done()
	// })

	// exit := make(chan struct{})
	// SignalC := make(chan os.Signal)

	// signal.Notify(SignalC, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// go func() {
	// 	for s := range SignalC {
	// 		switch s {
	// 		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
	// 			exit <- struct{}{}
	// 		}
	// 	}
	// }()

	// <-exit
	// httpServer.Close(nil)
	// os.Exit(0)
}
