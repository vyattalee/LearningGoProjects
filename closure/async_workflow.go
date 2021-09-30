package main

import "fmt"

func main() {
	a, b := 1, 2
	// b := 2
	go func() {
		result := doWork1(a, b)
		result = doWork2(result)
		result = doWork3(result)
		// Use the final result
	}()
	fmt.Println("hi!")
	fmt.Println("vim-go")
}

// func doWork(a, b, function(result) {
//   // use the result here
// })

func doWork1() func(a int, b int) {
	fmt.Println("doWork1: ", a, b)
}
func doWork2() func(r func(int, int)) {
	fmt.Println("doWork2: ", r)
}
func doWork3() func(r func(int, int)) {
	fmt.Println("doWork3: ", r)
}
