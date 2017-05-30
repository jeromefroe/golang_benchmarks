package main

import (
	"math/rand"
	"testing"
	"time"
)

const rangeBound int32 = 1357

var globalVar int32

func BenchmarkStandardBoundedRandomNumber(b *testing.B) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		globalVar = r.Int31n(rangeBound)
	}
}

func BenchmarkBiasedFastBoundedRandomNumber(b *testing.B) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		random := int64(r.Int31())
		multiResult := random * int64(rangeBound)
		globalVar = int32(multiResult >> 32)
	}
}

func BenchmarkUnbiasedFastBoundedRandomNumber(b *testing.B) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		random := int64(r.Int31())
		multiResult := random * int64(rangeBound)
		leftover := int32(multiResult)
		if leftover < rangeBound {
			threshold := -rangeBound % rangeBound
			for leftover < threshold {
				random = int64(r.Int31())
				multiResult = random * int64(rangeBound)
				leftover = int32(multiResult)
			}
		}
		globalVar = int32(multiResult >> 32)
	}
}
