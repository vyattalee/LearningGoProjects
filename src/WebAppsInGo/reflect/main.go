package main

import (
	"fmt"
	"reflect"
)

func main() {

	var x float64 = 3.4
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)
	fmt.Println(x, "'type is ", t)
	fmt.Println("type:", v.Type())
	fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
	fmt.Println("value:", v.Float())

	var y float64 = 5.5
	p := reflect.ValueOf(&y)
	ep := p.Elem()
	fmt.Println("before setfloat y = ", y)
	ep.SetFloat(6.6)
	fmt.Println("after setfloat y = ", y)

	slice := make([]int, 0, 4)
	slice = append(slice, 1, 2, 3)
	TestSlice(slice)
	fmt.Println(slice)

	slice1 := make([]int, 0, 4)
	slice1 = append(slice1, 1, 2, 3)

	slice2 := slice1[:len(slice1)-1]
	slice2 = append(slice2, 11, 12, 13, 14, 15)
	slice2[0] = 10

	fmt.Println(slice1)
	fmt.Println(slice2)

}
func TestSlice(slice []int) {
	slice = append(slice, 4)
}
