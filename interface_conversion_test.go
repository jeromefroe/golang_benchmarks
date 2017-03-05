package main

import "testing"

type dog interface {
	Bark() string
}

type poodle struct{}

func (p *poodle) Bark() string { return "woof" }

func play(d dog) {
	d.Bark()
}

func BenchmarkInterfaceNoConversion(b *testing.B) {
	var d dog
	d = &poodle{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		play(d)
	}
}

func BenchmarkInterfaceConversion(b *testing.B) {
	p := &poodle{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		play(p)
	}
}
