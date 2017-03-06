package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"unsafe"
)

const sampleStr = "It’s the job that’s never started as takes longest to finish"

var sampleBytes = []byte(sampleStr)

func BenchmarkWriteBytes(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	var w io.Writer
	w = buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.Write(sampleBytes)
	}
}

func BenchmarkWriteString(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	var w io.Writer
	w = buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.Write([]byte(sampleStr))
	}
}

func BenchmarkWriteUnafeString(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	var w io.Writer
	w = buf
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		w.Write(unsafeStrToByte(sampleStr))
	}
}

func unsafeStrToByte(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	byteHeader := reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&byteHeader))
}
