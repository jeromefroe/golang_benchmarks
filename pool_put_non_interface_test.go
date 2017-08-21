package main

import (
	"sync"
	"testing"

	"github.com/m3db/m3x/pool"
)

func BenchmarkPoolM3XPutSlice(b *testing.B) {
	var p = pool.NewObjectPool(pool.NewObjectPoolOptions())
	p.Init(func() interface{} {
		return make([]byte, 0, 1024)
	})
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := p.Get().([]byte)
			buf = buf[:0]

			// Simulate work.
			for i := 0; i < 256; i++ {
				buf = append(buf, 100)
			}

			p.Put(buf)
		}
	})
}

func BenchmarkPoolM3XPutPointerToSlice(b *testing.B) {
	var p = pool.NewObjectPool(pool.NewObjectPoolOptions())
	p.Init(func() interface{} {
		b := make([]byte, 0, 1024)
		return &b
	})
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := p.Get().(*[]byte)
			*buf = (*buf)[:0]

			// Simulate work.
			for i := 0; i < 256; i++ {
				*buf = append(*buf, 100)
			}

			p.Put(buf)
		}
	})
}

func BenchmarkPoolSyncPutSlice(b *testing.B) {
	var p = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 256)
		},
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := p.Get().([]byte)
			buf = buf[:0]

			// Simulate work.
			for i := 0; i < 256; i++ {
				buf = append(buf, 100)
			}

			p.Put(buf)
		}
	})
}

func BenchmarkPoolSyncPutPointerToSlice(b *testing.B) {
	var p = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 0, 256)
			return &b
		},
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := p.Get().(*[]byte)
			*buf = (*buf)[:0]

			// Simulate work.
			for i := 0; i < 256; i++ {
				*buf = append(*buf, 100)
			}

			p.Put(buf)
		}
	})
}
