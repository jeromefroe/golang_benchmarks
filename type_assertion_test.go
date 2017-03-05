package main

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkTypeAssertion(b *testing.B) {
	var r io.Reader
	r = new(bytes.Buffer)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := r.(*bytes.Buffer); !ok {
			b.Fatal()
		}
	}
}
