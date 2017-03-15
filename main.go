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
	Resolve(hash int, replica uint) string
}

func main() {
	srvLow := 32
	srvHigh := 40
	replicasLow := 3
	replicasHigh := 3

	fmt.Println("hashringmap")
	fmt.Println("--")
	fmt.Println("start  end  replicas  numMoved  moved % ideal %  diff %")
	for i := srvLow; i < srvHigh; i++ {
		moved(i, replicasLow, replicasHigh, func(servers []string, replicas int) hasher {
			ring := newHashRingMap(replicas)
			ring.Add(servers...)
			return ring
		})
	}
	fmt.Println("hashjump")
	fmt.Println("--")
	fmt.Println("start  end  replicas  numMoved  moved % ideal %  diff %")
	for i := srvLow; i < srvHigh; i++ {
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

func (m *hashRingMap) Resolve(v int, replica uint) string {
	return m.hashMap[m.instances[m.index(v, replica)]]
}

func (m *hashRingMap) index(v int, replica uint) int {
	if replica > 0 {
		v += int(replica)
	}

	var buffer [8]byte
	key := buffer[:8]
	binary.LittleEndian.PutUint64(key, uint64(v))
	hash := int(m.hash(key))

	idx := sort.SearchInts(m.instances, hash)
	if idx == len(m.instances) {
		idx = 0
	}

	return idx
}

type hashJump struct {
	hash     func(b []byte) uint32
	servers  []string
	replicas int
}

func newHashJump(replicas int) *hashJump {
	return &hashJump{
		hash:     murmur3.Sum32,
		replicas: replicas,
	}
}

func (j *hashJump) Add(instances ...string) {
	j.servers = append(j.servers, instances...)
}

func (j *hashJump) Resolve(hash int, replica uint) string {
	return j.servers[int(jump.Hash(uint64(hash+int(replica)), int32(len(j.servers))))]
}

func moved(numServers int, replicasLow, replicasHigh int, newHasherFn newHasherFn) {
	for replicas := replicasLow; replicas <= replicasHigh; replicas++ {
		// we only need to look at shardSpace hashes because we take the hash modulo shardSpace
		servers := newServers(numServers)
		hasher := newHasherFn(servers, replicas)

		initialMappings := make(map[int][]string, shardSpace)
		for i := 0; i < shardSpace; i++ {
			var targets []string
			for r := uint(0); r < uint(replicas); r++ {
				targets = append(targets, hasher.Resolve(i, r))
			}
			initialMappings[i] = targets
		}

		hasher = newHasherFn(newServers(numServers+1), replicas)

		updatedMappings := make(map[int][]string, shardSpace)
		for i := 0; i < shardSpace; i++ {
			var targets []string
			for r := uint(0); r < uint(replicas); r++ {
				targets = append(targets, hasher.Resolve(i, r))
			}
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
