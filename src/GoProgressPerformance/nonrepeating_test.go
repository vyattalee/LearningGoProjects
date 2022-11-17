package main

import "testing"

func BenchmarkSubstr(b *testing.B) {
	s := "黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花"
	for i := 0; i < 13; i++ {
		s = s + s
	}
	b.Logf("len(s) = %d", len(s))
	ans := 8
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		actual := lengthOfNonRepeatSubStr(s)
		if actual != ans {
			b.Errorf("got %d for input %s; "+"expected %d", actual, s, ans)
		}
	}
}
