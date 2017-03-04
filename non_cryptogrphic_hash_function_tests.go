package main

import (
	"hash/adler32"
	"hash/crc32"
	"hash/crc64"
	"hash/fnv"
	"math/rand"
	"testing"

	"github.com/OneOfOne/xxhash"
	"github.com/dchest/siphash"
	"github.com/dgryski/go-farm"
	"github.com/dgryski/go-highway"
	"github.com/dgryski/go-spooky"
	"github.com/spaolacci/murmur3"
	"github.com/zhenjl/cityhash"
)

const testString = "kahfkgjfq2348r742gydi71382rvkjyaci71138yakdkchvk73773"

var testBytes = []byte(testString)

func BenchmarkHash32Fnv(b *testing.B) {
	h := fnv.New32()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum32()
	}
}

func BenchmarkHash32Fnva(b *testing.B) {
	h := fnv.New32a()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum32()
	}
}

func BenchmarkHash64Fnv(b *testing.B) {
	h := fnv.New64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum64()
	}
}

func BenchmarkHash64Fnva(b *testing.B) {
	h := fnv.New64a()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum64()
	}
}

func BenchmarkHash32Crc(b *testing.B) {
	h := crc32.NewIEEE()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum32()
	}
}

func BenchmarkHash64Crc(b *testing.B) {
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum64()
	}
}

func BenchmarkHash32Adler(b *testing.B) {
	h := adler32.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum32()
	}
}

func BenchmarkHash32Xxhash(b *testing.B) {
	h := xxhash.New32()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum32()
	}
}

func BenchmarkHash64Xxhash(b *testing.B) {
	h := xxhash.New64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum64()
	}
}

func BenchmarkHash32Murmur3(b *testing.B) {
	h := murmur3.New32()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum32()
	}
}

func BenchmarkHash128Murmur3(b *testing.B) {
	h := murmur3.New128()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Write(testBytes)
		h.Sum128()
	}
}

func BenchmarkHash64CityHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cityhash.CityHash64(testBytes, uint32(len(testBytes)))
	}
}

func BenchmarkHash128CityHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cityhash.CityHash128(testBytes, uint32(len(testBytes)))
	}
}

func BenchmarkHash32FarmHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		farm.Hash32(testBytes)
	}
}

func BenchmarkHash64FarmHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		farm.Hash64(testBytes)
	}
}

func BenchmarkHash128FarmHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		farm.Hash128(testBytes)
	}
}

func BenchmarkHash64SipHash(b *testing.B) {
	k0 := uint64(rand.Int63())
	k1 := uint64(rand.Int63())
	for i := 0; i < b.N; i++ {
		siphash.Hash(k0, k1, testBytes)
	}
}

func BenchmarkHash128SipHash(b *testing.B) {
	k0 := uint64(rand.Int63())
	k1 := uint64(rand.Int63())
	for i := 0; i < b.N; i++ {
		siphash.Hash128(k0, k1, testBytes)
	}
}

func BenchmarkHash64HighwayHash(b *testing.B) {
	keys := highway.Lanes{}
	for i := 0; i < b.N; i++ {
		highway.Hash(keys, testBytes)
	}
}

func BenchmarkHash32SpookyHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spooky.Hash32(testBytes)
	}
}

func BenchmarkHash64SpookyHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		spooky.Hash64(testBytes)
	}
}

func BenchmarkHash128SpookyHash(b *testing.B) {
	k0 := uint64(rand.Int63())
	k1 := uint64(rand.Int63())
	for i := 0; i < b.N; i++ {
		spooky.Hash128(testBytes, &k0, &k1)
	}
}
