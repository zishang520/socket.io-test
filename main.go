package main

import (
	"fmt"

	"github.com/zishang520/engine.io/types"
	"github.com/zishang520/engine.io/utils"
)

type X string

type A struct {
	Uid string
	AA  X
}

type B struct {
	Uid string
	AA  string
}

type Slice[F any, T any] []F

func (s Slice[F, T]) Map(_call func(value F) T) []T {
	r := make([]T, len(s))
	for i, v := range s {
		r[i] = _call(v)
	}
	return r
}

func (s Slice[F, T]) Range(_call func(value F, key int) bool) {
	for k, v := range s {
		if !_call(v, k) {
			break
		}
	}
}

func main() {
	s := types.NewSet("1", "2").Keys()
	fmt.Println(utils.MsgPack().Encode(s))
}
