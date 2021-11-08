package main

import (
	"fmt"
)

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func main() {
	myAdder := adder()

	// 从1加到10
	for i := 1; i <= 101; i++ {
		myAdder(i)
	}

	fmt.Println(myAdder(0))
	// 再加上45
	//fmt.Println(myAdder(45))
}
