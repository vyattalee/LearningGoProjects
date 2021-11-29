package main

//go build -v -ldflags="-X 'main.Version=v1.0.0' -X
//'build/build.User=$(id -u -n)' -X
//'build/build.Time=$(date)'"

import (
	"HowToCodeinGo/LDFlagsExample/build"
	"fmt"
)

var Version = "development"

func main() {
	fmt.Println("Version:\t", Version)
	fmt.Println("build.Time:\t", build.Time)
	fmt.Println("build.User:\t", build.User)
}
