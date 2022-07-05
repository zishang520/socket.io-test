package main

import (
	"fmt"
)

type A struct {
	s bool
}

type B struct {
	c bool
}

type C struct {
	A
	B
	s bool
}

func main() {
	a := interface{}(&C{})
	fmt.Printf("%v", a.(A))
}
