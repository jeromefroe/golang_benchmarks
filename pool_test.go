package main

import (
	"bytes"
	"sync"
	"testing"
)

type ChannelBufferPool struct {
	c chan *bytes.Buffer
	n int
}

func NewChannelBufferPool(size, n int) (bp *ChannelBufferPool) {
	return &ChannelBufferPool{
		c: make(chan *bytes.Buffer, size),
		n: n,
	}
}

func (p *ChannelBufferPool) Get() (b *bytes.Buffer) {
	select {
	case b = <-p.c:
	default:
		b = bytes.NewBuffer(make([]byte, 0, p.n))
	}
	return
}

func (p *ChannelBufferPool) Put(b *bytes.Buffer) {
	b.Reset()

	select {
	case p.c <- b:
	default:
	}
}

func BenchmarkAllocateBufferNoPool(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(10)
	b.ResetTimer()

	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < b.N/10; i++ {
				buf := bytes.NewBuffer(make([]byte, 0, 256))
				buf.WriteString("gotta catch 'em all")
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkChannelBufferPool(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(10)

	p := NewChannelBufferPool(1, 256)
	b.ResetTimer()

	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < b.N/10; i++ {
				buf := p.Get()
				buf.Reset()
				buf.WriteString("gotta catch 'em all")
				p.Put(buf)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkSyncBufferPool(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(10)

	var p = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 256))
		},
	}
	b.ResetTimer()

	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < b.N/10; i++ {
				buf := p.Get().(*bytes.Buffer)
				buf.Reset()
				buf.WriteString("gotta catch 'em all")
				p.Put(buf)
			}

			wg.Done()
		}()
	}

	wg.Wait()
}
