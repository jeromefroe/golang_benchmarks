package main

import "testing"

var rGlobal uint32

func BenchmarkReduceModuloPowerOfTwo(b *testing.B) {
	var (
		u uint32
		n uint32 = 256
	)
	for i := 0; i < b.N; i++ {
		rGlobal = u % n
		u++
	}
}

func BenchmarkReduceModuloNonPowerOfTwo(b *testing.B) {
	var (
		u uint32
		n uint32 = 257
	)
	for i := 0; i < b.N; i++ {
		rGlobal = u % n
		u++
	}
}

func BenchmarkReduceAlternativePowerOfTwo(b *testing.B) {
	var (
		u uint32
		n uint32 = 256
	)
	for i := 0; i < b.N; i++ {
		rGlobal = uint32(uint64(u) * uint64(n) >> 32)
		u++
	}
}

func BenchmarkReduceAlternativeNonPowerOfTwo(b *testing.B) {
	var (
		u uint32
		n uint32 = 257
	)
	for i := 0; i < b.N; i++ {
		rGlobal = uint32(uint64(u) * uint64(n) >> 32)
		u++
	}
}
