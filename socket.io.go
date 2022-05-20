package main

import (
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/parser"
)

func main() {
	utils.Log().DEBUG = true
	var t float64 = 1.9
	utils.Log().Info("%v", int(t))

	encode := parser.NewParser().Encoder()
	decode := parser.NewParser().Decoder()

	pack := &parser.Packet{
		Type: parser.EVENT,
		Nsp:  "/name/",
		Data: []interface{}{
			types.NewBytesBufferString("123456789"), // bytes
		},
		Id: 13,
	}
	packets := encode.Encode(pack)

	decode.On("decoded", func(args ...any) {
		packet := args[0].(*parser.Packet)
		packet.Type = parser.EVENT
		packets := encode.Encode(packet)
		utils.Log().Info("[write] decode %v", packets)
	})
	decode.Add(`1`)

	for _, packet := range packets {
		decode.Add(packet)
	}

}
