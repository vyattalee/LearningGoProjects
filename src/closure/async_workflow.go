package main

import "fmt"

func main() {
	a, b := 1, 2
	// b := 2
	// go func() {
	result := doWork1(a, b)
	result = doWork2(result)
	result = doWork3(result)
	// Use the final result
	// }()
	fmt.Println("hi!")
	fmt.Println("vim-go")
}

// func doWork(a, b, function(result) {
//   // use the result here
// })

func doWork1(a int, b int) func() {
	return func() {
		fmt.Println("doWork1: ", a, b)
	}
}
func doWork2(r func()) func() {
	return func() {
		fmt.Println("doWork2: ", r)
	}
}
func doWork3(r func()) func() {
	return func() {
		fmt.Println("doWork3: ", r)
	}
}
