package main

import "fmt"

func main() {
	defer setup()()
	f := setup()
	defer f()
}

func setup() func() {
	fmt.Println("pretend to set things up")

	return func() {
		fmt.Println("pretend to tear things down")
	}
}
