package main

import (
	"testing"
)

func BenchmarkIndexRangeArray(b *testing.B) {
	var buf [16]byte
	for i := 0; i < b.N; i++ {
		for j := range buf {
			buf[j] = buf[j] + 1
		}
	}
}

func BenchmarkIndexValueRangeArray(b *testing.B) {
	var buf [16]byte
	for i := 0; i < b.N; i++ {
		for j, v := range buf {
			buf[j] = v + 1
		}
	}
}

func BenchmarkIndexValueRangeArrayPtr(b *testing.B) {
	var buf [16]byte
	for i := 0; i < b.N; i++ {
		for j, v := range &buf {
			buf[j] = v + 1
		}
	}
}
