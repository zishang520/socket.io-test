package main

import (
	// "bytes"
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
	// "io"
	// "strings"
)

func main() {
	a := types.NewStringBufferString("aaaaaaaa,")
	str, err := a.ReadString('x')
	utils.Log().Success("io.Reader %v %v", str, err)
}
