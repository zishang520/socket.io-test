package main

import (
	"fmt"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/socket.io/parser"
)

func main() {
	_parser := parser.NewParser()
	data := _parser.Encoder().Encode(&parser.Packet{
		Type: parser.EVENT,
		Nsp:  "/name",
		Data: []any{
			types.NewBytesBuffer([]byte{0, 1, 2, 3, 4, 5}),
		},
		/*"aaaa",
		types.NewStringBufferString("xxx"),
		types.NewBytesBuffer([]byte{0, 1, 2, 3, 4, 5}),
		types.NewBytesBuffer([]byte{0, 1, 2, 3, 4, 5}),
		types.NewBytesBuffer([]byte{0, 1, 2, 3, 4, 5}),*/

	})
	_parser.Decoder().On("decoded", func(args ...any) {
		fmt.Printf("decoded ----> %v\n", args[0].(*parser.Packet))
	})
	// fmt.Printf("d -> %v\n", _parser.Decoder().Add(data))
	fmt.Printf("%v\n", data)
	for _, v := range data {
		fmt.Printf("d -> %v\n", _parser.Decoder().Add(v))
	}
}
