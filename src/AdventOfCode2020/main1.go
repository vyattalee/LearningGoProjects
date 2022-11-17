package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	contents, err := ioutil.ReadFile("input1.txt")
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[int]struct{}, 0)
	for _, line := range bytes.Split(contents, []byte("\n")) {
		n, err := strconv.Atoi(string(line))
		if err != nil {
			log.Println(err)
			continue
		}
		m[n] = struct{}{}
	}

	for x, _ := range m {
		y := 2020 - x
		if _, ok := m[y]; ok {
			log.Printf("%d+%d=2020 / %d*%d=%d", x, y, x, y, x*y)
			return
		}
	}
}
