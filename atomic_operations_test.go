package main

import (
	"sync/atomic"
	"testing"
)

func BenchmarkAtomicLoad32(b *testing.B) {
	var v int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.LoadInt32(&v)
		}
	})
}

func BenchmarkAtomicLoad64(b *testing.B) {
	var v int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.LoadInt64(&v)
		}
	})
}

func BenchmarkAtomicStore32(b *testing.B) {
	var v int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.StoreInt32(&v, 1)
		}
	})
}

func BenchmarkAtomicStore64(b *testing.B) {
	var v int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.StoreInt64(&v, 1)
		}
	})
}

func BenchmarkAtomicAdd32(b *testing.B) {
	var v int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt32(&v, 1)
		}
	})
}

func BenchmarkAtomicAdd64(b *testing.B) {
	var v int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt64(&v, 1)
		}
	})
}

func BenchmarkAtomicCAS32(b *testing.B) {
	var v int32
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.CompareAndSwapInt32(&v, 0, 0)
		}
	})
}

func BenchmarkAtomicCAS64(b *testing.B) {
	var v int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.CompareAndSwapInt64(&v, 0, 0)
		}
	})
}
