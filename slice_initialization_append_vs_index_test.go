package main

import "testing"

var sl = make([]interface{}, 10)

func BenchmarkSliceInitializationAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := make([]interface{}, 0, len(sl))
		for _, iface := range sl {
			b = append(b, iface)
		}
	}
}

func BenchmarkSliceInitializationIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := make([]interface{}, len(sl))
		for i, iface := range sl {
			b[i] = iface
		}
	}
}
