package main

import (
	"bytes"
	"strings"
	"testing"
)

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
