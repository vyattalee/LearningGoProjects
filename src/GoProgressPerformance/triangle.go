package main

import (
	"fmt"
	"math"
)

func triangle() {
	var a, b int = 3, 4
	c := calTriangle(a, b)
	fmt.Println(c)
}

func calTriangle(a, b int) int {
	var c int
	c = int(math.Sqrt(float64(a*a + b*b)))
	return c
}
