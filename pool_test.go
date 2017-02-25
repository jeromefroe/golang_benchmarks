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
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			buf := bytes.NewBuffer(make([]byte, 0, 256))
			buf.WriteString("gotta catch 'em all")
		}
	})
}

func BenchmarkChannelBufferPool(b *testing.B) {
	p := NewChannelBufferPool(1, 256)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			buf := p.Get()
			buf.Reset()
			buf.WriteString("gotta catch 'em all")
			p.Put(buf)
		}
	})
}

func BenchmarkSyncBufferPool(b *testing.B) {
	var p = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 256))
		},
	}

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {
			buf := p.Get().(*bytes.Buffer)
			buf.Reset()
			buf.WriteString("gotta catch 'em all")
			p.Put(buf)
		}
	})
}
