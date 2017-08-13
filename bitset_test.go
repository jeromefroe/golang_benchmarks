package main

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"github.com/willf/bitset"
)

func BenchmarkBitsetRoaringConsecutive1000(b *testing.B) {
	benchmarkBitsetRoaringConsecutive(b, 1000)
}

func BenchmarkBitsetWillfConsecutive1000(b *testing.B) {
	benchmarkBitsetWillfConsecutive(b, 1000)
}

func BenchmarkBitsetRoaringConsecutive10000(b *testing.B) {
	benchmarkBitsetRoaringConsecutive(b, 10000)
}

func BenchmarkBitsetWillfConsecutive10000(b *testing.B) {
	benchmarkBitsetWillfConsecutive(b, 10000)
}

func BenchmarkBitsetRoaringConsecutive100000(b *testing.B) {
	benchmarkBitsetRoaringConsecutive(b, 100000)
}

func BenchmarkBitsetWillfConsecutive100000(b *testing.B) {
	benchmarkBitsetWillfConsecutive(b, 100000)
}

func BenchmarkBitsetRoaringConsecutive1000000(b *testing.B) {
	benchmarkBitsetRoaringConsecutive(b, 1000000)
}

func BenchmarkBitsetWillfConsecutive1000000(b *testing.B) {
	benchmarkBitsetWillfConsecutive(b, 1000000)
}

func benchmarkBitsetRoaringConsecutive(b *testing.B, end uint32) {
	for i := 0; i < b.N; i++ {
		rb := roaring.NewBitmap()
		for i := uint32(0); i < end; i++ {
			rb.Add(i)
		}
	}
}

func benchmarkBitsetWillfConsecutive(b *testing.B, end uint) {
	for i := 0; i < b.N; i++ {
		b := bitset.New(end)
		for i := uint(0); i < end; i++ {
			b.Set(i)
		}
	}
}

// func benchmarkBitsetRoaringRandom(nums []uint32) {
// 	rb := roaring.NewBitmap()
// 	for _, num := range nums {
// 		rb.Add(num)
// 	}
// }

// func benchmarkBitsetWillfRandom(nums []uint) {
// 	b := bitset.New(len(nums))
// 	for _, num := range nums {
// 		b.Set(num)
// 	}
// }
