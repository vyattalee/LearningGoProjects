package main

import "testing"

func TestTriangle(t *testing.T) {
	tests := []struct{ a, b, c int }{
		{3, 4, 5},
		{5, 12, 13},
		{8, 15, 17},
		{12, 35, 37},
		{3000, 4000, 5000},
		//{3000, 4000, 5001},
	}
	for _, tt := range tests {
		if actual := calTriangle(tt.a, tt.b); actual != tt.c {
			t.Errorf("calTriangle(%d, %d); got %d; expected %d", tt.a, tt.b, tt.c, calTriangle(tt.a, tt.b))
		}
	}
}

func BenchmarkTriangle(b *testing.B) {
	a1, b1, c1 := 30000, 40000, 50000
	for i := 0; i < b.N; i++ {
		actual := calTriangle(a1, b1)
		if actual != c1 {
			b.Errorf("calTriangle(%d, %d); got %d; expected %d", a1, b1, c1, calTriangle(a1, b1))
		}
	}
}
