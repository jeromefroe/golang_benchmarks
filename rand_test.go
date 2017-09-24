package main

import (
	"math/rand"
	"testing"
)

func BenchmarkGlobalRandInt63(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Int63()
		}
	})
}

func BenchmarkLocalRandInt63(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		src := rand.NewSource(rand.Int63())
		r := rand.New(src)
		for pb.Next() {
			r.Int63()
		}
	})
}

func BenchmarkGlobalRandFloat64(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Float64()
		}
	})
}

func BenchmarkLocalRandFloat64(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		src := rand.NewSource(rand.Int63())
		r := rand.New(src)
		for pb.Next() {
			r.Float64()
		}
	})
}
