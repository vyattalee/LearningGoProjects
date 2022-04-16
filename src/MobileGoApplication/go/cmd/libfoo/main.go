// go/cmd/libfoo/main.go
package main

import "C"

// other imports should be seperate from the special Cgo import
import (
	"rogchap.com/libfoo/foo"
)

//export reverse
func reverse(in *C.char) *C.char {
	return C.CString(foo.Reverse(C.GoString(in)))
}

func main() {}
