# Golang Benchmarks

Various benchmarks for different patterns in Go. Some of these are implementations of or were
inspired by
[this excellent article on performance in Go](http://bravenewgeek.com/so-you-wanna-go-fast/).
Furthemore, the Golang wiki provides a
[list of compiler optimizations](https://github.com/golang/go/wiki/CompilerOptimizations).

> Lies, damned lies, and benchmarks.

### Allocate on Stack vs Heap

[`allocate_stack_vs_heap_test.go`](allocate_stack_vs_heap_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkAllocateFooStack          | 1000000000 | 2.27 ns/op  |    0 B/op | 0 allocs/op
BenchmarkAllocateBarStack          | 1000000000 | 2.27 ns/op  |    0 B/op | 0 allocs/op
BenchmarkAllocateFooHeap           | 50000000   | 29.0 ns/op  |   32 B/op | 1 allocs/op
BenchmarkAllocateBarHeap           | 50000000   | 30.2 ns/op  |   32 B/op | 1 allocs/op
BenchmarkAllocateSliceHeapNoEscape | 50000000   | 32.3 ns/op  |    0 B/op | 0 allocs/op
BenchmarkAllocateSliceHeapEscape   | 5000000    |  260 ns/op  | 1024 B/op |1 allocs/op

Generated using go version go1.7.5 darwin/amd64

This benchmark just looks at the difference in performance between allocating a struct on the stack
versus on the heap. As expected, allocating a struct on the stack is much faster than allocating it
on the heap. The two structs I used in the benchmark are below:

```
type Foo struct {
	foo int64
	bar int64
	baz int64
}

type Bar struct {
	foo int64
	bar int64
	baz int64
	bah int64
}
```

One interesting thing is that although `Foo` is only 24 bytes, when we allocate it on the heap, 32 bytes are
allocated and when `Bar` is allocated on the heap, 32 bytes are allocated for it as well. When I first
saw this, my initial suspicion was that Go's memory allocator allocates memory in certain bin sizes instead
of the exact size of the struct, and there is no bin size between 24 and 32 bytes, so `Foo` was allocated
the next highest bin size, which was 32 bytes. This
[blog post](https://medium.com/@robertgrosse/optimizing-rc-memory-usage-in-rust-6652de9e119e#.x2kfg63oh)
examines a similar phenomenen in Rust and its memory allocator jemalloc. As for Go, I found the following
in the file
[runtime/malloc.go](https://github.com/golang/go/blob/01c6a19e041f6b316c17a065f7a42b8dab57c9da/src/runtime/malloc.go#L27):

```
// Allocating a small object proceeds up a hierarchy of caches:
//
//	1. Round the size up to one of the small size classes
//	   and look in the corresponding mspan in this P's mcache.
//  ...
```

The last two benchmarks look at an optimization the Go compiler performs. If it can prove through
[escape analysis](https://en.wikipedia.org/wiki/Escape_analysis) that a slice does not escape the calling
function, then it allocates the data for the slice on the stack instead of the heap. More information can
be found on this [golang-nuts post](https://groups.google.com/forum/#!topic/golang-nuts/KdbtOqNK6JQ).

### Append

[`append_test.go`](append_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkAppendLoop     |   500000 | 2456 ns/op | 0 B/op | 0 allocs/op
BenchmarkAppendVariadic | 20000000 | 97.1 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.8.1 darwin/amd64

This benchmark looks at the performance difference between appending the values
of one slice into another slice one by one, i.e. `dst = append(dst, src[i])`,
versus appending them all at once, i.e. `dst = append(dst, src...)`. As the
benchmarks show, using the variadic approach is faster. My suspicion is that this
is because the compiler can optimize this away into a single `memcpy`.

### Atomic Operations

[`atomic_operations_test.go`](atomic_operations_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkAtomicLoad32  | 2000000000 | 1.77 ns/op
BenchmarkAtomicLoad64  | 1000000000 | 1.79 ns/op
BenchmarkAtomicStore32 |   50000000 | 27.5 ns/op
BenchmarkAtomicStore64 |  100000000 | 25.2 ns/op
BenchmarkAtomicAdd32   |   50000000 | 27.1 ns/op
BenchmarkAtomicAdd64   |   50000000 | 27.8 ns/op
BenchmarkAtomicCAS32   |   50000000 | 28.8 ns/op
BenchmarkAtomicCAS64   |   50000000 | 28.6 ns/op

Generated using go version go1.7.5 darwin/amd64

These benchmarks look at various atomic operations on 32 and 64 bit integers. The only thing that
really stands out is that loads are significantly faster than all other operations. I suspect that
there's two reasons for this: there's no cache invalidation because only reads are performed and
[on x86_64 loads and stores using `movq` are atomic if performed on natural alignments](http://preshing.com/20130618/atomic-vs-non-atomic-operations/).
I took a look at the `Load64` function in
[src/sync/atomic/asm_amd64.go](https://github.com/golang/go/blob/master/src/sync/atomic/asm_amd64.s):

```
TEXT ·LoadInt64(SB),NOSPLIT,$0-16
	JMP	·LoadUint64(SB)

TEXT ·LoadUint64(SB),NOSPLIT,$0-16
	MOVQ	addr+0(FP), AX
	MOVQ	0(AX), AX
	MOVQ	AX, val+8(FP)
	RET
```

It uses [Go's assembly language](https://golang.org/doc/asm) which I'm not too familiar with, but
it appears to move the address of the integer into the AX register in the first function, move the
value pointed to by that address into the AX register in the second instruction, and then move that
value into the return value of the function in the third instruction. On x86_64 the Go assembly
likely can be translated exactly using the
[`movq` instruction](http://www.felixcloutier.com/x86/MOVQ.html) and since this instruction
is atomic if executed on natural alignments the load will be atomic as well.

### Bit Tricks

[`bit_tricks_test.go`](bit_tricks_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkBitTricksModPowerOfTwo       | 2000000000 | 0.84 ns/op
BenchmarkBitTricksModNonPowerOfTwo    | 2000000000 | 1.58 ns/op
BenchmarkBitTricksAnd                 | 2000000000 | 0.46 ns/op
BenchmarkBitTricksDividePowerOfTwo    | 2000000000 | 0.72 ns/op
BenchmarkBitTricksDivideNonPowerOfTwo | 2000000000 | 1.09 ns/op
BenchmarkBitTricksShift               | 2000000000 | 0.52 ns/op

Generated using go version go1.8.1 darwin/amd64

These benchmarks look at some micro optimizations that can be performed when doing division or
modulo division. The first three benchmarks show the overhead of doing modulus division and how
we can replace modulus division by a power of two with a bitwise and, which is a faster operation.
Likewise the last three benchmarks show the overhead of division and how we can improve the speed
of division by a power of two by performing a right shift.

### Buffered vs Synchronous Channel

[`buffered_vs_unbuffered_channel_test.go`](buffered_vs_unbuffered_channel_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkSynchronousChannel | 5000000  | 240 ns/op
BenchmarkBufferedChannel    | 10000000 | 108 ns/op

Generated using go version go1.7.5 darwin/amd64

This benchmark examines the speed with which one can put objects onto a channel and comes from this
[golang-nuts forum post](https://groups.google.com/forum/#!topic/golang-nuts/ec9G0MGjn48). Using a buffered
channels is over twice as fast as using a synchronous channel which makes sense since the goroutine that is
putting objects into the channel need not wait until the object is taken out of the channel before placing
another object into it.

### Channel vs Ring Buffer

[`channel_vs_ring_buffer_test.go`](channel_vs_ring_buffer_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkChannelSPSC    | 20000000 |      102 ns/op   |    8 B/op |    1 allocs/op
BenchmarkRingBufferSPSC | 20000000 |     72.2 ns/op   |    8 B/op |    1 allocs/op
BenchmarkChannelSPMC    |  3000000 |      464 ns/op   |    8 B/op |    1 allocs/op
BenchmarkRingBufferSPMC |  1000000 |     1065 ns/op   |    8 B/op |    1 allocs/op
BenchmarkChannelMPSC    |  3000000 |      447 ns/op   |    8 B/op |    1 allocs/op
BenchmarkRingBufferMPSC |   300000 |     5033 ns/op   |    9 B/op |    1 allocs/op
BenchmarkChannelMPMC    |    10000 |   193557 ns/op   | 8016 B/op | 1000 allocs/op
BenchmarkRingBufferMPMC |       30 | 34618237 ns/op   | 8000 B/op | 1000 allocs/op

Generated using go version go1.8.3 darwin/amd64

The blog post [So You Wanna Go Fast?](http://bravenewgeek.com/so-you-wanna-go-fast/) also took a look at using
channels versus using a
[lock-free ring buffer](https://github.com/Workiva/go-datastructures/blob/master/queue/ring.go). I decided to
run similar benchmarks myself and the results are above. The benchmarks SPSC, SPMC, MPSC, and MPMC refer to
Single Producer Single Consumer, Single Producer Mutli Consumer, Mutli Producer Single Consumer, and
Mutli Producer Mutli Consumer respectively. The blog post found that for the SPSC case, a channel was faster
than a ring buffer when the tests were run on a single thread (`GOMAXPROCS=1`) but the ring buffer was faster
when the tests were on multiple threads (`GOMAXPROCS=8`). The blog post also examined the the SPMC and MPMC
cases and found similar results. That is, channels were faster when run on a single thread and the ring buffer
was faster when the tests were run with multiple threads. I ran all the test with `GOMAXPROCS=4` which is
the number of CPU cores on the machine I ran the tests on (a 2015 MacBook Pro with a 3.1 GHz Intel Core i7
Processor, it has 2 physical CPUs,`sysctl hw.physicalcpu`, and 4 logical CPUs, `sysctl hw.logicalcpu`).
Ultimately, the benchmarks I ran produced different results. They show that in the SPSC and SPMC cases the
performance of a channel and ring buffer are similar with the ring buffer holding a small advantage. However,
in the MPSC and MPMC a channel performed much better than a ring buffer did.

### defer

[`defer_test.go`](defer_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkMutexUnlock     | 50000000 | 25.8 ns/op
BenchmarkMutexDeferUnlock| 20000000 | 92.1 ns/op

Generated using go version go1.7.5 darwin/amd64

`defer` carries a slight performance cost, so for simple use cases it may be preferable
to call any cleanup code manually. As this [blog post](http://bravenewgeek.com/so-you-wanna-go-fast/)
notes, `defer` can be called from within conditional blocks and must be called if a functions
panics as well. Therefore, the compiler can't simply add the deferred function
wherever the function returns and instead `defer` must be more nuanced, resuling in the performance
hit. There is, in fact, an [open issue](https://github.com/golang/go/issues/14939) to address the
performance cost of `defer`.
[Another discussion](http://grokbase.com/t/gg/golang-nuts/158zz5p42w/go-nuts-defer-performance)
suggests calling `defer mu.Unlock()` before one calls `mu.Lock()` so the defer call will
be moved out of the critical path:

```
defer mu.Unlock()
mu.Lock()
```

### False Sharing

[`false_sharing_test.go`](false_sharing_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkIncrementFalseSharing                |  3000 | 453087 ns/op
BenchmarkIncrementNoFalseSharing              |  5000 | 246124 ns/op
BenchmarkIncrementNoFalseSharingLocalVariable | 20000 | 71624 ns/op

go version go1.8.3 darwin/amd64

This example demonstrates the effects of false sharing when multiple goroutines are updating
a variable. In the first benchmark, although the goroutines are each updating different variables
because those variables are on the same cache line, the updates contend with one another. In the
second example, however, we introduce some padding to ensure the integers are on different cache
lines so the updates won't interfere with each other. Finally, the last example performs the
increments locally and then writes the variable to the shared slice.

### Function Call

[`function_call_test.go`](function_call_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkPointerToStructMethodCall | 2000000000 | 0.32 ns/op
BenchmarkInterfaceMethodCall       | 2000000000 | 1.90 ns/op
BenchmarkFunctionPointerCall       | 2000000000 | 1.91 ns/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at the overhead for three different kinds of function calls: calling a method
on a pointer to a struct, calling a method on an interface, and calling a function through a function
pointer field in a struct. As expected, the method call on the pointer to the struct is the fastest since
the compiler knows what function is being called at compile time, whereas the others do not. For example,
the interface method call relies on dynamic dispatch at runtime to determine which function call and
likewise the function pointer to call is determined at runtime as well and has almost identical performance
to the interface method call.

### Interface conversion

[`interface_conversion_test.go`](interface_conversion_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkInterfaceConversion   | 2000000000 | 1.32 ns/op | 0 B/op | 0 allocs/op
BenchmarkNoInterfaceConversion | 2000000000 | 0.85 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.8.1 darwin/amd64

This benchmark looks at the overhead of converting an interface to its concrete type. Surprisingly,
the overhead of the type assertion, while not zero, it pretty minimal at only about 0.5 nanoseconds.

### Memset optimization

[`memset_test.go`](memset_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkSliceClearZero/1K      | 100000000 |  12.9 ns/op
BenchmarkSliceClearZero/16K     |  10000000 |   167 ns/op
BenchmarkSliceClearZero/128K    |    300000 |  3994 ns/op
BenchmarkSliceClearNonZero/1K   |   3000000 |   497 ns/op
BenchmarkSliceClearNonZero/16K  |    200000 |  7891 ns/op
BenchmarkSliceClearNonZero/128K |     20000 | 79763 ns/op

Generated using go version go1.9.2 darwin/amd64

This benchmark looks at the
[Go compiler's optimization for clearing slices](https://github.com/golang/go/wiki/CompilerOptimizations#idioms)
to the respective type's zero value. Specifically, if `s` is a slice or
an array then the following loop is optimized with memclr calls:

```
for i := range s {
	a[i] = <zero value for element of s>
}
```

If the value is not the the zero value of the type though then the loop
is not optimized as the benchmarks show. The library
[`go-memset`](https://github.com/tmthrgd/go-memset) provides a function
which optimizes clearing byte slices with any value not just zero.

### Mutex

[`mutex_test.go`](mutex_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkNoMutexLock     | 2000000000 | 1.18 ns/op
BenchmarkRWMutexReadLock |   30000000 | 54.5 ns/op
BenchmarkRWMutexLock     |   20000000 | 96.0 ns/op
BenchmarkMutexLock       |   20000000 | 78.7 ns/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at the cost of acquiring different kinds of locks. In the first benchmark we
don't acquire any lock. In the second benchmark we acquire a read lock on a `RWMutex`. In the third
we acquire a write lock on a `RWMutex`. And in the last benchmark we acquire a regular `Mutex` lock.

### Non-cryptographic Hash functions

[`non_cryptographic_hash_functions_test.go`](non_cryptographic_hash_functions_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkHash32Fnv         |  20000000 |  70.3 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash32Fnva        |  20000000 |  70.4 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64Fnv         |  20000000 |  71.1 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64Fnva        |  20000000 |  77.1 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash32Crc         |  30000000 |  87.5 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64Crc         |  10000000 |   175 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash32Adler       |  30000000 |  40.3 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash32Xxhash      |  30000000 |  46.1 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64Xxhash      |  30000000 |  47.4 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash32Murmur3     |  20000000 |  59.4 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash128Murmur3    |  20000000 |  63.4 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64CityHash    |  30000000 |  57.4 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash128CityHash   |  20000000 |   113 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash32FarmHash    |  30000000 |  44.4 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64FarmHash    |  50000000 |  26.4 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash128FarmHash   |  30000000 |  40.3 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64SipHash     |  30000000 |  39.3 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash128SipHash    |  30000000 |  44.9 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64HighwayHash |  50000000 |  36.9 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash32SpookyHash  |  30000000 |  58.1 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64SpookyHash  |  20000000 |  62.7 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash128SpookyHash |  30000000 |  68.2 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHashMD5           |  10000000 |   169 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash64MetroHash   | 100000000 |  18.6 ns/op | 0 B/op  | 0 allocs/op
BenchmarkHash128MetroHash  |  30000000 |  48.8 ns/op | 0 B/op  | 0 allocs/op

Generated using go version go1.8.3 darwin/amd64

These benchmarks look at the speed of various non-cryptographic hash function implementations in Go.

### Pass By Value vs Reference

[`pass_by_value_vs_reference_test.go`](pass_by_value_vs_reference_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkPassByReferenceOneWord    |  1000000000  | 2.20 ns/op
BenchmarkPassByValueOneWord        |  1000000000  | 2.58 ns/op
BenchmarkPassByReferenceFourWords  |   500000000  | 2.71 ns/op
BenchmarkPassByValueFourWords      |  1000000000  | 2.78 ns/op
BenchmarkPassByReferenceEightWords |  1000000000  | 2.32 ns/op
BenchmarkPassByValueEightWords     |   300000000  | 4.35 ns/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at the performance cost of passing a variable by reference vs passing it by value. For
small structs there doesn't appear to be much of a difference, but as the structs gets larger we start to
see a bit of difference which is to be expected since the larger the struct is the more words that have to
be copied into the function's stack when passed by value.

### Pool

[`pool_test.go`](pool_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkAllocateBufferNoPool | 20000000 |  118 ns/op | 368 B/op | 2 allocs/op
BenchmarkChannelBufferPool    | 10000000 |  213 ns/op |  43 B/op | 0 allocs/op
BenchmarkSyncBufferPool       | 50000000 | 27.7 ns/op |   0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

This benchmark compares three different memory allocation schemes. The first approach just
allocates its buffer on the heap normally. After it's done using the buffer it will eventually be
garbage collected. The second approach uses Go's sync.Pool type to pool buffers which caches objects
between runs of the garbage collector. The last approach uses a channel to permanently pool objects.
The difference between the last two approaches is the sync.Pool dynamically resizes itself and clears
items from the pool during a GC run. Two good resources to learn more about pools in Go
are these blog posts:
[Using Buffer Pools with Go](https://elithrar.github.io/article/using-buffer-pools-with-go/)
and
[How to Optimize Garbage Collection in Go](https://www.cockroachlabs.com/blog/how-to-optimize-garbage-collection-in-go/).

### Pool Put Non Interface

[`pool_put_non_interface_test.go`](pool_put_non_interface_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkPoolM3XPutSlice           |  5000000 | 282 ns/op | 32 B/op | 1 allocs/op
BenchmarkPoolM3XPutPointerToSlice  |  5000000 | 327 ns/op |  0 B/op | 0 allocs/op
BenchmarkPoolSyncPutSlice          | 10000000 | 184 ns/op | 32 B/op | 1 allocs/op
BenchmarkPoolSyncPutPointerToSlice | 10000000 | 177 ns/op |  0 B/op | 0 allocs/op

Generated using go version go1.8.3 darwin/amd64

This benchmark looks at the cost of pooling slices. Since slices are three words they cannot be
coerced into interfaces without an allocation, see the
[comments on this CL](https://go-review.googlesource.com/c/24371) for more details. Consequently,
putting a slice in a pool requires the three words for the slice to be allocated on the heap.
This cost will admittedly likely be offset by the savings from pooling the actual data backing
the slice, however this test looks performs the benchmarks to look at just that. Indeed we see
that although the putting a slice on a pool does require an additional allocation, there does not
appear to be a significant cost in speed.

### Rand

[`rand_test.go`](rand_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkGlobalRandInt63   |  20000000 |    115 ns/op | 0 B/op | 0 allocs/op
BenchmarkLocalRandInt63    | 300000000 |   3.95 ns/op | 0 B/op | 0 allocs/op
BenchmarkGlobalRandFloat64 |  20000000 |   96.5 ns/op | 0 B/op | 0 allocs/op
BenchmarkLocalRandFloat64  | 200000000 |   6.00 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.8.3 darwin/amd64

Go's [math/rand package](https://golang.org/pkg/math/rand/) exposes various functions for generating
random numbers, for example [`Int63`](https://golang.org/pkg/math/rand/#Int63). These functions use a
global [`Rand` struct](https://golang.org/pkg/math/rand/#Rand) which is
[created by the package](https://github.com/golang/go/blob/master/src/math/rand/rand.go#L235). This
struct
[uses a lock to serialize access to its random number source](https://github.com/golang/go/blob/master/src/math/rand/rand.go#L316)
though which can lead to contention if multiple goroutines are all trying to generate random numbers
using the global struct. Consequently, these benchmarks look at the performance improvement that comes from
giving each goroutine its own `Rand` struct so they don't need to acquire a lock. This
[blog post](http://blog.sgmansfield.com/2016/01/the-hidden-dangers-of-default-rand/) explores the similar
optimizations for using the math/rand package for those who are interested.

### Random Bounded Numbers

[`random_bounded_test.go`](random_bounded_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkStandardBoundedRandomNumber     | 100000000 | 18.3 ns/op
BenchmarkBiasedFastBoundedRandomNumber   | 100000000 | 11.0 ns/op
BenchmarkUnbiasedFastBoundedRandomNumber | 50000000  | 40.5 ns/op

Generated using go version go1.8.1 darwin/amd64

Benchmarks for three different algorithims for generating a random bounded number
as discussed in the blog post
[Fast random shuffling](http://lemire.me/blog/2016/06/30/fast-random-shuffling/).
The top result is the standard approach of generating a random number and taking
its modulus of the bound. The second approach implements the algorithim discussed
in the aforementioned blog post which avoids using the modulus operator. The third
algorithim is an implementation of the second algorithim which is unbiased. As
mentioned in the article, for most applications the second algorithim will be
sufficient enough as the bias introduced by the algorithim is likely less than
the bias from the pseudo-random number generator which is used.

### Range over Arrays and Slices

[`range_test.go`](range_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkIndexRangeArray         | 100000000 | 10.6 ns/op | 0 B/op | 0 allocs/op
BenchmarkIndexValueRangeArray    | 100000000 | 14.1 ns/op | 0 B/op | 0 allocs/op
BenchmarkIndexValueRangeArrayPtr | 100000000 | 10.1 ns/op | 0 B/op | 0 allocs/op
enchmarkIndexSlice               | 100000000 | 10.4 ns/op | 0 B/op | 0 allocs/op
BenchmarkIndexValueSlice         | 100000000 | 10.3 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.8.3 darwin/amd64

These tests look at three different ways to range over an array or slice. The first three benchmarks
range over an array. The first uses just the index into the array (`for i := range a`), the second
uses both the index and the value (`for i, v := range a`), and the third uses the index and value
while ranging over a pointer to an array (`for i, v := range &a`). What's interesting to note is
that the second benchmark is noticably slower than the other two. This is because go
[makes a copy of the array when you range over the index and value](https://groups.google.com/forum/#!topic/golang-dev/35W8LvT51vg).
Another example of this can been in
[this tweet by Damian Gryski](https://twitter.com/dgryski/status/816226596835225600) and there is
even [a linter to catch this](https://github.com/mdempsky/rangerdanger). The last two benchmarks
look at ranging over a slice. The first uses just the index into the slice and the second uses
both the index and the value. Unlike in the case of an array, there is no difference in performance
here.

### Reducing an Integer

[`reduction_test.go`](reduction_test.go)

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkReduceModuloPowerOfTwo         | 500000000  | 3.41 ns/op
BenchmarkReduceModuloNonPowerOfTwo      | 500000000  | 3.44 ns/op
BenchmarkReduceAlternativePowerOfTwo    | 2000000000 | 0.84 ns/op
BenchmarkReduceAlternativeNonPowerOfTwo | 2000000000 | 0.84 ns/op

Generated using go version go1.8.1 darwin/amd64

This benchmark compares two different approaches for reducing an integer into a given range. The first
two benchmarks use the traditional approach of taking the modulus of a given integer with the length
of the range that we want to reduce the integer into. The latter two benchmarks implement an alternative
approach that was described in
[A fast alternative to the modulo reduction](http://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/).
As the benchmarks show the alternative approach provides superior performance. This alternative
implementation was invented because modulus division is a slow instruction on modern processors in
comparison to other common instructions, and while one could replace a modulus division
by a power of two with a bitwise AND one cannot do the same for a value which is not a power of two.
This alternative approach is fair in that every integer in the range [0,N) will have either
ceil(2^32/N) or floor(s^32/N) integers in the range [0,2^32) mapped to it. However, unlike modulus
division which preserves the lower bits of information (so that k and k+1 will be mapped to different
integers if N != 1) the alternative implementation instead preserves the higher order bits (so k and
k+1 have a much higher likelohood of being mapped to the same integer) which means it can't be used
in hashmaps which use probing to resolve collisions since probing often adds the probe bias to the
lower bits (for example, linear probing adds 1 to the hash value) though one can certainly imagine
using a probing function which adds the probe bias to the higher order bits.

### Slice Initialization Append vs Index

[`slice_initialization_append_vs_index_test.go`](slice_initialization_append_vs_index_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkSliceInitializationAppend | 10000000 | 132 ns/op | 160 B/op | 1 allocs/op
BenchmarkSliceInitializationIndex  | 10000000 | 119 ns/op | 160 B/op | 1 allocs/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at slice initialization with `append` versus using an explicit index. I ran this benchmark
a few times and it seesawed back and forth. Ultimately, I think they compile down into the same code so there
probably isn't any actual performance difference. I'd like to take an actual look at the assembly
that they are compiled to and update this section in the future.

### String Concatenation

[`string_concatenation_test.go`](string_concatenation_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkStringConcatenation      | 20000000 |  83.9 ns/op |  64 B/op | 1 allocs/op
BenchmarkStringBuffer             | 10000000 |   131 ns/op |  64 B/op | 1 allocs/op
BenchmarkStringJoin               | 10000000 |   144 ns/op | 128 B/op | 2 allocs/op
BenchmarkStringConcatenationShort | 50000000 |  25.4 ns/op |   0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at the three different ways to perform string concatenation, the first uses the builtin `+`
operator, the second uses a `bytes.Buffer` and the third uses `string.Join`. It seems using `+` is preferable
to either of the other approaches which are similar in performance.

The last benchmark highlights a neat optimization Go performs when concatenating strings with `+`. The
documentation for the string concatenation function in
[runtime/string.go](https://github.com/golang/go/blob/d7b34d5f29324d77fad572676f0ea139556235e0/src/runtime/string.go)
states:

```
// The constant is known to the compiler.
// There is no fundamental theory behind this number.
const tmpStringBufSize = 32

type tmpBuf [tmpStringBufSize]byte

// concatstrings implements a Go string concatenation x+y+z+...
// The operands are passed in the slice a.
// If buf != nil, the compiler has determined that the result does not
// escape the calling function, so the string data can be stored in buf
// if small enough.
func concatstrings(buf *tmpBuf, a []string) string {
  ...
}
```

That is, if the compiler determines that the resulting string does not escape the calling function it will
allocate a 32 byte buffer on the stack which can be used as the underlying buffer for the string if it 32
bytes or less. In the last benchmark, the resulting string is in fact less than 32 bytes so it can be stored
on the stack saving a heap allocation.

### Type Assertion

[`type_assertion_test.go`](type_assertion_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkTypeAssertion | 2000000000 | 0.97 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at the performance cost of a type assertion. I was a little surprised to find
it was so cheap.

### Write Bytes vs String

[`write_bytes_vs_string_test.go`](write_bytes_vs_string_test.go)

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkWriteBytes       | 100000000 | 18.7 ns/op |  0 B/op | 0 allocs/op
BenchmarkWriteString      | 20000000  | 63.3 ns/op | 64 B/op | 1 allocs/op
BenchmarkWriteUnafeString | 100000000 | 21.1 ns/op |  0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

Go's [`io.Writer` interface](https://golang.org/pkg/io/#Writer) only has one `Write` method which
takes a byte slice as an argument. To pass a string to it though requires a conversion to a byte
slice which entails a heap allocation. These benchmarks look at the performance cost of writing a
byte slice, converting a string to a byte slice and then writing it, and using the `unsafe` and
`reflect` packages to create a byte slice to the data underlying the string without an allocation.
