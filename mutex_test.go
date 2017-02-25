package main

import (
	"sync"
	"testing"
)

func BenchmarkNoMutexLock(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
		}
	})
}

func BenchmarkRWMutexReadLock(b *testing.B) {
	var rw sync.RWMutex
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rw.RLock()
			rw.RUnlock()
		}
	})
}

func BenchmarkRWMutexLock(b *testing.B) {
	var rw sync.RWMutex
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rw.Lock()
			rw.Unlock()
		}
	})
}

func BenchmarkMutexLock(b *testing.B) {
	var mu sync.Mutex
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			mu.Unlock()
		}
	})
}
