package main

import (
	"testing"
)

func BenchmarkIndexRangeArray(b *testing.B) {
	var buf [16]byte
	for n := 0; n < b.N; n++ {
		for i := range buf {
			buf[i] = buf[i] + 1
		}
	}
}

func BenchmarkIndexValueRangeArray(b *testing.B) {
	var buf [16]byte
	for n := 0; n < b.N; n++ {
		for i, v := range buf {
			buf[i] = v + 1
		}
	}
}

func BenchmarkIndexValueRangeArrayPtr(b *testing.B) {
	var buf [16]byte
	for n := 0; n < b.N; n++ {
		for i, v := range &buf {
			buf[i] = v + 1
		}
	}
}

func BenchmarkIndexSlice(b *testing.B) {
	buf := make([]byte, 16)
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(buf); i++ {
			buf[i] = buf[i] + 1
		}
	}
}

func BenchmarkIndexValueSlice(b *testing.B) {
	buf := make([]byte, 16)
	for n := 0; n < b.N; n++ {
		for i, v := range buf {
			buf[i] = v + 1
		}
	}
}
