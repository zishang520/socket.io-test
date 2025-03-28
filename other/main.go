package main

import "fmt"

func def(v int) {
	fmt.Println(v)
}

func main() {
	defer def(1)
	defer def(2)
	defer def(3)
	defer def(4)
	def(5)
}
