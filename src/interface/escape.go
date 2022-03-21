package main

import "fmt"

type Addifier interface{ Add(a, b int32) int32 }

type Adder struct {
	name string
	id   int32
}

//go:noinline
func (adder Adder) Add(a, b int32) int32 { return a + b }

func main() {
	adder := Adder{name: "myAdder"}
	fmt.Println(adder.Add(10, 32))           // doesn't escape
	fmt.Println(Addifier(adder).Add(10, 32)) // escapes
}
