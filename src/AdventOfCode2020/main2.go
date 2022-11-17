package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	contents, err := ioutil.ReadFile("input2.txt")
	if err != nil {
		log.Fatal(err)
	}
	m := make(map[int]struct{}, 0)
	for _, line := range bytes.Split(contents, []byte("\n")) {
		n, err := strconv.Atoi(string(line))
		if err != nil {
			continue
		}
		m[n] = struct{}{}
	}

	for x, _ := range m {
		s := 2020 - x
		for y, _ := range m {
			z := s - y
			if _, ok := m[z]; ok {
				log.Printf("%d+%d+%d=2020 / %d*%d*%d=%d", x, y, z, x, y, z, x*y*z)
			}
		}
	}
}
