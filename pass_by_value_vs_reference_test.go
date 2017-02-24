package main

import "testing"

type OneWord int64

type FourWords struct {
	a OneWord
	b OneWord
	c OneWord
	d OneWord
}

type EightWords struct {
	a FourWords
	b FourWords
}

func BenchmarkPassByReferenceOneWord(b *testing.B) {
	f := func(s *OneWord) {}
	s := OneWord(0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(&s)
	}
}

func BenchmarkPassByValueOneWord(b *testing.B) {
	f := func(s OneWord) {}
	s := OneWord(0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(s)
	}
}

func BenchmarkPassByReferenceFourWords(b *testing.B) {
	f := func(s *FourWords) {}
	s := FourWords{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(&s)
	}
}

func BenchmarkPassByValueFourWords(b *testing.B) {
	f := func(s FourWords) {}
	s := FourWords{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(s)
	}
}

func BenchmarkPassByReferenceEightWords(b *testing.B) {
	f := func(s *EightWords) {}
	s := EightWords{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(&s)
	}
}

func BenchmarkPassByValueEightWords(b *testing.B) {
	f := func(s EightWords) {}
	s := EightWords{}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(s)
	}
}
