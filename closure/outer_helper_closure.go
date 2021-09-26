package main

import "fmt"

func main() {
	fmt.Println("vim-go")
	fmt.Println(outer(5))
	fmt.Println(outer(6))
}

func outer(i int) []int {
	var out []int
	myClosure := func(b int) {
		out = []int{i, b}

	}
	helper(myClosure)
	return out

}

func helper(f func(int)) {
	f(4)
}
