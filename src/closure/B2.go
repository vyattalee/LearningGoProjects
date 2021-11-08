package main

import (
	"fmt"
)

func B() []func() {
	b := make([]func(), 3, 3)
	for i := 0; i < 3; i++ {
		b[i] = func(j int) func() {
			return func() {
				fmt.Println(j)
			}
		}(i)
	}
	return b
}

func main() {
	c := B()
	c[0]()
	c[1]()
	c[2]()
}
