package main

import (
	"sync"
	"testing"

	"github.com/Workiva/go-datastructures/queue"
)

func BenchmarkChannelSPSC(b *testing.B) {
	ch := make(chan interface{}, 128)
	var wg sync.WaitGroup
	wg.Add(1)
	b.ResetTimer()

	go func() {
		for i := 0; i < b.N; i++ {
			ch <- 100
		}
	}()

	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
		wg.Done()
	}()

	wg.Wait()
}

func BenchmarkRingBufferSPSC(b *testing.B) {
	q := queue.NewRingBuffer(128)
	var wg sync.WaitGroup
	wg.Add(1)
	b.ResetTimer()

	go func() {
		for i := 0; i < b.N; i++ {
			q.Put(100)
		}
	}()

	go func() {
		for i := 0; i < b.N; i++ {
			q.Get()
		}
		wg.Done()
	}()

	wg.Wait()
}

func BenchmarkChannelSPMC(b *testing.B) {
	ch := make(chan interface{}, 128)
	var wg sync.WaitGroup
	wg.Add(1000)
	b.ResetTimer()

	go func() {
		for i := 0; i < b.N; i++ {
			ch <- 100
		}
	}()

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N/1000; i++ {
				<-ch
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkRingBufferSPMC(b *testing.B) {
	q := queue.NewRingBuffer(128)
	var wg sync.WaitGroup
	wg.Add(1000)
	b.ResetTimer()

	go func() {
		for i := 0; i < b.N; i++ {
			q.Put(100)
		}
	}()

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N/1000; i++ {
				q.Get()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkChannelMPSC(b *testing.B) {
	ch := make(chan interface{}, 128)
	var wg sync.WaitGroup
	wg.Add(1)
	b.ResetTimer()

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				ch <- 100
			}
		}()
	}

	go func() {
		for i := 0; i < b.N; i++ {
			<-ch
		}
		wg.Done()
	}()

	wg.Wait()
}

func BenchmarkRingBufferMPSC(b *testing.B) {
	q := queue.NewRingBuffer(128)
	var wg sync.WaitGroup
	wg.Add(1)
	b.ResetTimer()

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				q.Put(100)
			}
		}()
	}

	go func() {
		for i := 0; i < b.N; i++ {
			q.Get()
		}
		wg.Done()
	}()

	wg.Wait()
}

func BenchmarkChannelMPMC(b *testing.B) {
	ch := make(chan interface{}, 128)
	var wg sync.WaitGroup
	wg.Add(1000)
	b.ResetTimer()

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				ch <- 100
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				<-ch
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkRingBufferMPMC(b *testing.B) {
	q := queue.NewRingBuffer(128)
	var wg sync.WaitGroup
	wg.Add(1000)
	b.ResetTimer()

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				q.Put(100)
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				q.Get()
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
