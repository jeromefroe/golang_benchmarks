package main

import (
	"runtime"
	"sync"
	"testing"
)

var (
	numFalseSharingIters = 100000
	numCPUs              = runtime.GOMAXPROCS(0)
)

type paddedInt64 struct {
	n int64
	_ [8]uint64
}

func BenchmarkIncrementFalseSharing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			wg   sync.WaitGroup
			ints = make([]int64, numCPUs)
		)
		wg.Add(numCPUs)

		for j := 0; j < numCPUs; j++ {
			go func(j int) {
				for k := 0; k < numFalseSharingIters; k++ {
					ints[j]++
				}
				wg.Done()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkIncrementNoFalseSharing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			wg   sync.WaitGroup
			ints = make([]paddedInt64, numCPUs)
		)
		wg.Add(numCPUs)

		for j := 0; j < numCPUs; j++ {
			go func(j int) {
				for k := 0; k < numFalseSharingIters; k++ {
					ints[j].n++
				}
				wg.Done()
			}(j)
		}

		wg.Wait()
	}
}

func BenchmarkIncrementNoFalseSharingLocalVariable(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			wg   sync.WaitGroup
			ints = make([]int64, numCPUs)
		)
		wg.Add(numCPUs)

		for j := 0; j < numCPUs; j++ {
			go func(j int) {
				var tmp int64
				for k := 0; k < numFalseSharingIters; k++ {
					tmp++
				}
				ints[j] = tmp
				wg.Done()
			}(j)
		}

		wg.Wait()
	}
}
