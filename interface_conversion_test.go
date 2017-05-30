package main

import "testing"

var (
	dogBark string
)

type dog interface {
	bark() string
}

type poodle struct{}

func (p *poodle) bark() string { return "woof" }

func playWithDog(d dog) string { return d.bark() }

func BenchmarkInterfaceNoConversion(b *testing.B) {
	var d dog = &poodle{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dogBark = playWithDog(d)
	}
}

func BenchmarkInterfaceConversion(b *testing.B) {
	p := &poodle{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dogBark = playWithDog(p)
	}
}
