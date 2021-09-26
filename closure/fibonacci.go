package main

import (
	"fmt"
)

func fibonacci() func() int {
	b0 := 0
	b1 := 1
	return func() int {
		tmp := b0 + b1
		b0 = b1
		b1 = tmp
		return b1
	}

}

func main() {
	myFibonacci := fibonacci()
	for i := 1; i <= 15; i++ {
		fmt.Println(myFibonacci())
	}
}
