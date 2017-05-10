package main

import "testing"

var btGlobal int

func BenchmarkBitTricksModPowerOfTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btGlobal = i % 256
	}
}

func BenchmarkBitTricksModNonPowerOfTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btGlobal = i % 257
	}
}

func BenchmarkBitTricksAnd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btGlobal = i & 256
	}
}

func BenchmarkBitTricksDividePowerOfTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btGlobal = i / 256
	}
}

func BenchmarkBitTricksDivideNonPowerOfTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btGlobal = i / 257
	}
}

func BenchmarkBitTricksShift(b *testing.B) {
	for i := 0; i < b.N; i++ {
		btGlobal = i >> 8
	}
}
