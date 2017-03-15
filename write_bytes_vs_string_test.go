package main

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const sampleStr = "It’s the job that’s never started as takes longest to finish"

var (
	sampleBytes                   = []byte(sampleStr)
	errInvalidByteSliceConversion = errors.New("invalid byte slice passed to conversion, slice have same the same length and capacity")
)

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
	var b []byte
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	// we need to take the address of s's Data field and assign it to b's Data field in one
	// expression as it as a uintptr and in the future Go may have a compacting GC that moves
	// pointers but it will not update uintptr values, but single expressions should be safe.
	// For more details see https://groups.google.com/forum/#!msg/golang-dev/rd8XgvAmtAA/p6r28fbF1QwJ
	byteHeader.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data

	// need to take the length of s here to ensure s is live until after we update b's Data
	// field since the garbage collector can collect a variable once it is no longer used
	// not when it goes out of scope, for more details see https://github.com/golang/go/issues/9046
	l := len(s)
	byteHeader.Len = l
	byteHeader.Cap = l
	return b
}

func unsafeByteToStr(b []byte) string {
	// need to assert that the slice's length and capacity are equal to avoid a memory leak
	// when converting to a string
	if len(b) != cap(b) {
		panic(errInvalidByteSliceConversion)
	}

	var s string
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))

	// we need to take the address of b's Data field and assign it to s's Data field in one
	// expression as it as a uintptr and in the future Go may have a compacting GC that moves
	// pointers but it will not update uintptr values, but single expressions should be safe.
	// For more details see https://groups.google.com/forum/#!msg/golang-dev/rd8XgvAmtAA/p6r28fbF1QwJ
	strHeader.Data = (*reflect.SliceHeader)(unsafe.Pointer(&b)).Data

	// need to take the length of b here to ensure b is live until after we update s's Data
	// field since the garbage collector can collect a variable once it is no longer used
	// not when it goes out of scope, for more details see https://github.com/golang/go/issues/9046
	strHeader.Len = len(b)
	return s
}

func TestUnsafeStrToByte(t *testing.T) {
	s := "fizzbuzz"
	expected := []byte(s)
	assert.Equal(t, expected, unsafeStrToByte(s))
}

func TestUnsafeByteToStr(t *testing.T) {
	b := []byte{'f', 'i', 'z', 'z', 'b', 'u', 'z', 'z'}
	expected := string(b)
	assert.Equal(t, expected, unsafeByteToStr(b))
}
