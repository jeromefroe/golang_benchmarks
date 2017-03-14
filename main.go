package main

import (
	"fmt"
	"math"

	"github.com/renstrom/go-jump-consistent-hash"
)

const shardSpace = 4096

type hashFn func(hash, clientsLen, replicas int) []int64

func main() {
	fmt.Println("start  end  replicas  numMoved  moved % ideal %  diff %")
	for i := 32; i < 40; i++ {
		moved(i, hashOld)
	}
	for i := 32; i < 40; i++ {
		moved(i, hashJump)
	}
}

func hashOld(hash, clientsLen, replicas int) []int64 {
	indexes := make([]int64, 0, replicas)

	// Use virtual nodes to rotate by 1/n when changing the backing caches
	vnodesPerClient := int(math.Ceil(float64(shardSpace) / float64(clientsLen)))
	currIndex := (hash % shardSpace) / vnodesPerClient
	indexes = append(indexes, int64(currIndex))

	rotateBy := int(math.Ceil(float64(clientsLen) / float64(replicas)))
	for i := 0; i < replicas-1; i++ {
		currIndex = (currIndex + rotateBy) % clientsLen
		containedAlready := false
		for _, index := range indexes {
			if index == int64(currIndex) {
				containedAlready = true
				break
			}
		}
		if !containedAlready {
			indexes = append(indexes, int64(currIndex))
		}
	}

	return indexes
}

func hashJump(hash, clientsLen, replicas int) []int64 {
	// would need to check clientsLen > replicas in real life

	indexes := make([]int64, 0, replicas)

	// take modulo 4096 to be consistent with hashOld
	hash = hash % shardSpace

	idx := jump.Hash(uint64(hash), int32(clientsLen))
	for i := 0; i < replicas; i++ {
		indexes = append(indexes, int64((idx+int32(i))%shardSpace))
	}

	return indexes
}

func moved(numClients int, hashFn hashFn) {
	for replicas := 1; replicas <= 5; replicas++ {
		// we only need to look at shardSpace hashes because we take the hash modulo shardSpace
		initialMappings := make(map[int][]int64, shardSpace)
		for i := 0; i < shardSpace; i++ {
			targets := hashFn(i, numClients, replicas)
			initialMappings[i] = targets
		}

		updatedMappings := make(map[int][]int64, shardSpace)
		for i := 0; i < shardSpace; i++ {
			targets := hashFn(i, numClients+1, replicas)
			updatedMappings[i] = targets
		}

		var moved int

		for i := 0; i < shardSpace; i++ {
			moved += missing(initialMappings[i], updatedMappings[i])
		}

		movedPercentage := float64(moved) / float64(shardSpace*replicas)
		idealPercentage := 1.0 / float64(numClients+1)
		diffPercentage := movedPercentage - idealPercentage
		fmt.Printf("%-8d%-8d%-8d%-8d%-8.4f%-8.4f%-8.4f\n", numClients, numClients+1, replicas, moved, movedPercentage, idealPercentage, diffPercentage)
	}
}

func missing(left, right []int64) int {
	var missing int
	for _, l := range left {
		var found bool
		for _, r := range right {
			if l == r {
				found = true
				break
			}
		}
		if !found {
			missing++
		}
	}
	return missing
}
