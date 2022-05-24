package main

import (
	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
	"github.com/zishang520/socket.io/socket"
)

func main() {
	utils.Log().DEBUG = true
	a := socket.NewBroadcastOperator(nil, nil, nil, nil)
	utils.Log().Info("%v", a.To("xxxx", "xxxxxxxxxx"))
	utils.Log().Info("%v", a.Compress(true))
	utils.Log().Info("%v", a)
}
