package main

import (
	"math/rand"
	"testing"
)

func BenchmarkGlobalRand(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Int63()
		}
	})
}

func BenchmarkLocalRand(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		src := rand.NewSource(rand.Int63())
		r := rand.New(src)
		for pb.Next() {
			r.Int63()
		}
	})
}
