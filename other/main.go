package main

import (
	"fmt"
)

type Implements[T any] interface {
	Prototype(T)

	Proto() T
}

type Extends[T any] struct {
	_proto_ T
}

func (e *Extends[T]) Prototype(_proto_ T) {
	e._proto_ = _proto_
}

func (e *Extends[T]) Proto() T {
	return e._proto_
}

type TestInterface interface {
	Implements[TestInterface]

	Out() string
	F() string
	Foo() string
}

type Test struct {
	Extends[TestInterface]
}

func (t *Test) Out() string {
	return t.Proto().F() + ":" + t.Foo()
}

func (t *Test) F() string {
	return "Test F"
}

func (t *Test) Foo() string {
	return "Test Foo"
}

type AInterface interface {
	TestInterface
}

type A struct {
	Test
}

func MakeA() *A {
	s := &A{}
	s.Prototype(s)

	return s
}

func NewA() *A {
	s := MakeA()
	s.Construct()
	return s
}
func (a *A) Construct() {
}
func (a *A) F() string {
	return "A F" + ":" + a.Test.F()
}

func (a *A) Foo() string {
	return "A Foo"
}

func main() {
	a := NewA()
	fmt.Println(a.Out())
}
