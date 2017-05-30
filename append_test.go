package main

import "testing"

const (
	sliceSize = 1024
)

var (
	srcSlice = make([]int64, sliceSize)
	dstSlice = make([]int64, 0, sliceSize)
)

func BenchmarkAppendLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dstSlice = dstSlice[:0]
		for j := 0; j < len(srcSlice); j++ {
			dstSlice = append(dstSlice, srcSlice[j])
		}
	}
}

func BenchmarkAppendVariadic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dstSlice = dstSlice[:0]
		dstSlice = append(dstSlice, srcSlice...)
	}
}
