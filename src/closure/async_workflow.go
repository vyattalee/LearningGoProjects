package main

import (
	"fmt"
	"sync"
)

func main() {

	var wg sync.WaitGroup
	a, b := 1, 2
	// b := 2
	wg.Add(1)
	go func() {
		// wg.Add(1)
		defer wg.Done()
		result := doWork1(a, b)
		result()
		result = doWork2(result)
		result()
		result = doWork3(result)
		result()
		// Use the final result
	}()

	wg.Wait()
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
