package main

import "testing"

type Doer interface {
	Do()
}

type Po func()

type A struct {
	Po Po
}

//go:noinline
func (s *A) Do() {}

//go:noinline
func po() {}

func BenchmarkPointerToStructMethodCall(b *testing.B) {
	a := &A{Po: po}
	for i := 0; i < b.N; i++ {
		a.Do()
	}
}

func BenchmarkInterfaceMethodCall(b *testing.B) {
	var a Doer = &A{Po: po}
	for i := 0; i < b.N; i++ {
		a.Do()
	}
}

func BenchmarkFunctionPointerCall(b *testing.B) {
	a := &A{Po: po}
	for i := 0; i < b.N; i++ {
		a.Po()
	}
}
