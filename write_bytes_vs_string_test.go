package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
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

func unsafeByteToStr(b []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&strHeader))
}

func TestUnsafeStrToByte(t *testing.T) {
	s := "fizzbuzz"
	expected := []byte(s)
	assert.Equal(t, expected, unsafeStrToByte(s))
}

func TestUnsafeByteToStr(t *testing.T) {
	expected := "fizzbuzz"
	b := []byte(expected)
	assert.Equal(t, expected, unsafeByteToStr(b))
}
