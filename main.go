package main

import (
	"fmt"
	"math"
)

const shardSpace = 4096

func hash(hash, clientsLen, replicas int) []int {
	indexes := make([]int, 0, replicas)

	// Use virtual nodes to rotate by 1/n when changing the backing caches
	vnodesPerClient := int(math.Ceil(float64(shardSpace) / float64(clientsLen)))
	currIndex := (hash % shardSpace) / vnodesPerClient
	indexes = append(indexes, currIndex)

	rotateBy := int(math.Ceil(float64(clientsLen) / float64(replicas)))
	for i := 0; i < replicas-1; i++ {
		currIndex = (currIndex + rotateBy) % clientsLen
		containedAlready := false
		for _, index := range indexes {
			if index == currIndex {
				containedAlready = true
				break
			}
		}
		if !containedAlready {
			indexes = append(indexes, currIndex)
		}
	}

	return indexes
}

func main() {
	// test with only 1 replica first
	replicas := 1

	initialClientsLen := 20
	// we only need to look at shardSpace hashes because we take the hash modulo shardSpace
	initialMappings := make(map[int][]int, shardSpace)
	for i := 0; i < shardSpace; i++ {
		targets := hash(i, initialClientsLen, replicas)
		initialMappings[i] = targets
	}

	updatedClientsLen := 21
	updatedMappings := make(map[int][]int, shardSpace)
	for i := 0; i < shardSpace; i++ {
		targets := hash(i, updatedClientsLen, replicas)
		updatedMappings[i] = targets
	}

	var moved int

	for i := 0; i < shardSpace; i++ {
		moved += missing(initialMappings[i], updatedMappings[i])
	}

	fmt.Println(moved)
	fmt.Println(float64(moved) / float64(shardSpace))
	fmt.Println(1.0 / 21.0)

	// now test with 3 replicas
	replicas = 3

	initialMappings = make(map[int][]int, shardSpace)
	for i := 0; i < shardSpace; i++ {
		initialMappings[i] = hash(i, initialClientsLen, replicas)
	}

	updatedMappings = make(map[int][]int, shardSpace)
	for i := 0; i < shardSpace; i++ {
		updatedMappings[i] = hash(i, updatedClientsLen, replicas)
	}

	moved = 0

	for i := 0; i < shardSpace; i++ {
		moved += missing(initialMappings[i], updatedMappings[i])
	}

	fmt.Println(moved)
	fmt.Println(float64(moved) / float64(3*shardSpace))
	fmt.Println(1.0 / 21.0)
}

func missing(left, right []int) int {
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
