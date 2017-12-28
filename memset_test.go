package main

import (
	"testing"
)

type mCase struct {
	name string
	len  int
}

var mCases = []mCase{
	{"1K", 1024},
	{"16K", 16 * 1024},
	{"128K", 128 * 1024},
}

func BenchmarkSliceClearZero(b *testing.B) {
	for _, mCase := range mCases {
		data := make([]byte, mCase.len)
		b.Run(mCase.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range data {
					data[j] = 0
				}
			}
		})
	}
}

func BenchmarkSliceClearNonZero(b *testing.B) {
	for _, mCase := range mCases {
		data := make([]byte, mCase.len)
		b.Run(mCase.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for j := range data {
					data[j] = 1
				}
			}
		})
	}
}
