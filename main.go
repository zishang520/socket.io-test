package main

import (
	"fmt"
)

type Test interface {
	Foo()
	Foo1()
}

type test struct {
	_self Test
}

func (t *test) self() Test {
	if t._self != nil {
		return t._self
	}
	return t
}
func (t *test) Super(c Test) {
	t._self = c
}

func (t *test) Foo() {
	fmt.Println("test foo")
}

func (t *test) Foo1() {
}

func (t *test) Test() {
	t.Foo()
	t.self().Foo()
	t.self().Foo1()
}

type A struct {
	*test
}

func (t *A) Foo() {
	fmt.Println("A foo")
}

func (t *A) Foo1() {
	fmt.Println("A foo1")
}

func main() {
	x := &A{&test{}}
	x.Super(x)
	x.Test()
}
