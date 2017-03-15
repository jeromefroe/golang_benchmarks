package main

import (
	"encoding/binary"
	"fmt"
	"sort"
	"strconv"

	"github.com/renstrom/go-jump-consistent-hash"
	"github.com/spaolacci/murmur3"
)

const shardSpace = 4096

type newHasherFn func(servers []string, replicas int) hasher

type hasher interface {
	Resolve(hash int) []string
}

func main() {
	replicasLow := 3
	replicasHigh := 3

	fmt.Println("hashringmap")
	fmt.Println("--")
	fmt.Println("start  end  replicas  numMoved  moved % ideal %  diff %")
	for i := 32; i < 40; i++ {
		moved(i, replicasLow, replicasHigh, func(servers []string, replicas int) hasher {
			ring := newHashRingMap(replicas)
			ring.Add(servers...)
			return ring
		})
	}
	fmt.Println("hashjump")
	fmt.Println("--")
	fmt.Println("start  end  replicas  numMoved  moved % ideal %  diff %")
	for i := 32; i < 40; i++ {
		moved(i, replicasLow, replicasHigh, func(servers []string, replicas int) hasher {
			jump := newHashJump(replicas)
			jump.Add(servers...)
			return jump
		})
	}
}

type hashRingMap struct {
	hash      func(b []byte) uint32
	replicas  int
	instances []int
	hashMap   map[int]string
}

func newHashRingMap(replicas int) *hashRingMap {
	m := &hashRingMap{
		replicas: replicas,
		hash:     murmur3.Sum32,
		hashMap:  make(map[int]string),
	}
	return m
}

func (m *hashRingMap) Add(instances ...string) {
	for _, instance := range instances {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + instance)))
			m.instances = append(m.instances, hash)
			m.hashMap[hash] = instance
		}
	}
	sort.Ints(m.instances)
}

func (m *hashRingMap) Resolve(v int) []string {
	var buffer [8]byte
	key := buffer[:8]
	binary.LittleEndian.PutUint64(key, uint64(v))
	hash := int(m.hash(key))

	servers := make([]string, 0, m.replicas)

	idx := sort.Search(len(m.instances), func(i int) bool { return m.instances[i] >= hash })
	if idx == len(m.instances) {
		idx = 0
	}

	servers = append(servers, m.hashMap[m.instances[idx]])
	for i := 0; i < m.replicas-1; i++ {
		idx = (idx + 1) % len(m.instances)
		servers = append(servers, m.hashMap[m.instances[idx]])
	}

	return servers
}

type hashJump struct {
	servers  []string
	replicas int
}

func newHashJump(replicas int) *hashJump {
	return &hashJump{replicas: replicas}
}

func (j *hashJump) Add(instances ...string) {
	j.servers = append(j.servers, instances...)
}

func (j *hashJump) Resolve(hash int) []string {
	servers := make([]string, 0, j.replicas)

	idx := int(jump.Hash(uint64(hash), int32(len(j.servers))))
	for i := 0; i < j.replicas; i++ {
		servers = append(servers, j.servers[idx%len(j.servers)])
		idx++
	}

	return servers
}

func moved(numServers int, replicasLow, replicasHigh int, newHasherFn newHasherFn) {
	for replicas := replicasLow; replicas <= replicasHigh; replicas++ {
		// we only need to look at shardSpace hashes because we take the hash modulo shardSpace
		hasher := newHasherFn(newServers(numServers), replicas)

		initialMappings := make(map[int][]string, shardSpace)
		for i := 0; i < shardSpace; i++ {
			targets := hasher.Resolve(i)
			initialMappings[i] = targets
		}

		hasher = newHasherFn(newServers(numServers+1), replicas)

		updatedMappings := make(map[int][]string, shardSpace)
		for i := 0; i < shardSpace; i++ {
			targets := hasher.Resolve(i)
			updatedMappings[i] = targets
		}

		var moved int

		for i := 0; i < shardSpace; i++ {
			moved += missing(initialMappings[i], updatedMappings[i])
		}

		movedPercentage := float64(moved) / float64(shardSpace*replicas)
		idealPercentage := 1.0 / float64(numServers+1)
		diffPercentage := movedPercentage - idealPercentage
		fmt.Printf("%-8d%-8d%-8d%-8d%-8.4f%-8.4f%-8.4f\n", numServers, numServers+1, replicas, moved, movedPercentage, idealPercentage, diffPercentage)
	}
}

func newServers(numServers int) []string {
	servers := make([]string, 0, numServers)
	for i := 0; i < numServers; i++ {
		servers = append(servers, fmt.Sprintf("server%02d", i))
	}
	return servers
}

func missing(left, right []string) int {
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
