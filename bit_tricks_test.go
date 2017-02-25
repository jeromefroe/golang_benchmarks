package main

import "testing"

func BenchmarkMod(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n = i % 2
	}
	_ = n
}

func BenchmarkAnd(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n = i & 1
	}
	_ = n
}

func BenchmarkDivide(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n = i / 2
	}
	_ = n
}

func BenchmarkShift(b *testing.B) {
	var n int
	for i := 0; i < b.N; i++ {
		n = i >> 1
	}
	_ = n
}
