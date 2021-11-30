package main

import (
	"HowToCodeinGo/BuildTagExample/os"
	"fmt"
	"strings"
)

func Join(parts ...string) string {
	return strings.Join(parts, os.PathSeparator)
}
func main() {
	s := Join("a", "b", "c")
	fmt.Println(s)
}
