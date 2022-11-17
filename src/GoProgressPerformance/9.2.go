package main

import "fmt"

func lengthOfNonRepeatSubStrOld(s string) int {
	lastOccurred := make(map[rune]int)
	start := 0
	maxLength := 0
	for i, ch := range []rune(s) {
		if lastI, ok := lastOccurred[ch]; ok && lastI >= start {
			start = lastI + 1
		}
		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}
		lastOccurred[ch] = i
	}
	return maxLength
}

var lastOccurred = make([]int, 0xffff) //假定中文字的最大值为65535=0xffff
func lengthOfNonRepeatSubStr(s string) int {
	// lastOccurred := make([]int, 0xffff) //假定中文字的最大值为65535=0xffff
	for i := range lastOccurred {
		lastOccurred[i] = -1
	}
	start := 0
	maxLength := 0
	for i, ch := range []rune(s) {
		if lastI := lastOccurred[ch]; lastI != -1 && lastI >= start {
			start = lastI + 1
		}
		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}
		lastOccurred[ch] = i
	}
	return maxLength
}

func main() {
	fmt.Println("this chapter 9.3")
	str := "黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花"
	fmt.Printf("%s lengthOfNonRepeatSubStr = %d ", str, lengthOfNonRepeatSubStr(str))
}
