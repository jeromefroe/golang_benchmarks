package main

import "testing"

var (
	myDog    dog
	myPoodle *poodle
)

type dog interface {
	bark() string
}

type poodle struct{}

func (p *poodle) bark() string { return "woof" }

func BenchmarkInterfaceConversion(b *testing.B) {
	var d dog = &poodle{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myDog = d.(*poodle)
	}
}

func BenchmarkNoInterfaceConversion(b *testing.B) {
	p := &poodle{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		myPoodle = p
	}
}
