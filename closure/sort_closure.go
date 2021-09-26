package main

import (
	"fmt"
	"sort"
)

func main() {

	numbers := []int{1, 11, -5, 7, 2, 0, 12}
	sort.Ints(numbers)
	fmt.Println("Sorted:", numbers)
	index := sort.SearchInts(numbers, 7)
	fmt.Println("7 is at index:", index)

}
