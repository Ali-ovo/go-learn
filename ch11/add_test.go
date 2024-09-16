package main

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	re := add(1, 2)

	if re != 3 {
		t.Error("1 + 2 should be 3")
	}

}

func TestAdd2(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode jump")
	}

	re := add(1, 5)
	if re != 6 {
		t.Errorf("1+5 should be 6")
	}
}

func TestAdd3(t *testing.T) {
	var dataset = []struct {
		a   int
		b   int
		out int
	}{
		{1, 1, 2},
		{
			-9, 8, -1,
		},
		{
			6, 6, 12,
		},
	}

	for _, v := range dataset {
		re := add(v.a, v.b)
		if re != v.out {
			t.Errorf("%d + %d should be %d", v.a, v.b, v.out)
		}
	}
}

func BenchmarkAdd(bb *testing.B) {
	a, b, c := 123, 456, 789

	for i := 0; i < bb.N; i++ {
		if actual := add(a, b); actual != c {
			fmt.Printf("%d + %d should be %d", a, b, c)
		}
	}
}
