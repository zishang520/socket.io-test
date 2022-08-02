package main

import (
	"github.com/zishang520/engine.io/utils"
)

type a struct {
}

func (x *a) aa() string {
	if x == nil {
		return "nil"
	}
	return "*a"
}

func main() {
	var x *a = &a{}
	utils.Log().Success("%v", x)
	utils.Log().Success("%v", x.aa())
}
