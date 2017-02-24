package main

import (
	"bytes"
	"strings"
	"testing"
)

// TODO: benchmark string.Join

// const vs var

var (
	bulbasur  = "bulbasaur"
	ivysaur   = "ivysaur"
	venasaur  = "venasur"
	squirtle  = "squirtle"
	wartortle = "wartortle"
	blastoise = "blastoise"
)

func concat() string {
	return bulbasur + ivysaur + venasaur + squirtle + wartortle + blastoise
}

func concatShort() string {
	return bulbasur + squirtle
}

func buffer(b *bytes.Buffer) string {
	b.Reset()

	b.WriteString(bulbasur)
	b.WriteString(ivysaur)
	b.WriteString(venasaur)
	b.WriteString(squirtle)
	b.WriteString(wartortle)
	b.WriteString(blastoise)

	return b.String()
}

func join() string {
	return strings.Join([]string{bulbasur, ivysaur, venasaur, squirtle, wartortle, blastoise}, "")
}

func BenchmarkStringConcatenation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = concat()
	}
}

func BenchmarkStringBuffer(b *testing.B) {
	buf := bytes.NewBuffer(make([]byte, 0, len(bulbasur)+len(ivysaur)+len(venasaur)+len(squirtle)+len(wartortle)+len(blastoise)))
	for i := 0; i < b.N; i++ {
		_ = buffer(buf)
	}
}

func BenchmarkStringJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = join()
	}
}

func BenchmarkStringConcatenationShort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = concatShort()
	}
}

// buffer would have additional overhead from pool contention

// The constant is known to the compiler.
// 10	// There is no fundamental theory behind this number.
// 11	const tmpStringBufSize = 32
// 12
// 13	type tmpBuf [tmpStringBufSize]byte
// 14
// 15	// concatstrings implements a Go string concatenation x+y+z+...
// 16	// The operands are passed in the slice a.
// 17	// If buf != nil, the compiler has determined that the result does not
// 18	// escape the calling function, so the string data can be stored in buf
// 19	// if small enough.
// 20	func concatstrings(buf *tmpBuf, a []string) string {
