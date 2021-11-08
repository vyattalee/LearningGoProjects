package main

import (
	"fmt"
)

func B() []func() {
	b := make([]func(), 3, 3)
	for i := 0; i < 3; i++ {
		j := i
		b[i] = func() {
			fmt.Println(j)
		}
	}
	return b
}

func main() {
	c := B()
	c[0]()
	c[1]()
	c[2]()
}
