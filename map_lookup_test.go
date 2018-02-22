package main

import (
	"math/rand"
	"testing"
)

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	benchSetSize = 1024
)

var (
	benchPlaceholder bool
)

func BenchmarkMapUint64(b *testing.B) {
	set := make(map[uint64]struct{}, benchSetSize)
	keys := make([]uint64, 0, benchSetSize)
	for i := 0; i < benchSetSize; i++ {
		key := rand.Uint64()
		set[key] = struct{}{}
		keys = append(keys, key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range keys {
			_, benchPlaceholder = set[k]
		}
	}
}

func BenchmarkMapString1(b *testing.B) {
	benchmarkString(b, 1)
}

func BenchmarkMapString10(b *testing.B) {
	benchmarkString(b, 10)
}

func BenchmarkMapString100(b *testing.B) {
	benchmarkString(b, 100)
}

func BenchmarkMapString1000(b *testing.B) {
	benchmarkString(b, 1000)
}

func BenchmarkMapString10000(b *testing.B) {
	benchmarkString(b, 10000)
}

func genString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func benchmarkString(b *testing.B, n int) {
	set := make(map[string]struct{}, benchSetSize)
	keys := make([]string, 0, benchSetSize)
	for i := 0; i < benchSetSize; i++ {
		key := genString(n)
		set[key] = struct{}{}
		keys = append(keys, key)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, k := range keys {
			_, benchPlaceholder = set[k]
		}
	}
}
