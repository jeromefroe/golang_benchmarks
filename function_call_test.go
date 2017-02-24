package main

import "testing"

// go test -bench Call -benchmem

type Doer interface {
	Do()
}

type Po func()

type A struct {
	Po Po
}

func (s *A) Do() {}

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

// TODO(jeromefroe): method call vs function field call, ie x.Foo() where Foo is a field, not a method
