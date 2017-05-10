# Golang Benchmarks

Various benchmarks for different patterns in Go. Many of these are implementations of or were
inspired by
[this excellent article on performance in Go](http://bravenewgeek.com/so-you-wanna-go-fast/).

> Lies, damned lies, and benchmarks.

### Allocate on Stack vs Heap

`allocate_stack_vs_heap_test.go`

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

### Atomic Operations

`atomic_operations_test.go`

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

`bit_tricks_test.go`

Benchmark Name|Iterations|Per-Iteration
----|----|----
BenchmarkBitTricksModPowerOfTwo-4               2000000000               0.84 ns/op
BenchmarkBitTricksModNonPowerOfTwo-4            2000000000               1.58 ns/op
BenchmarkBitTricksAnd-4                         2000000000               0.46 ns/op
BenchmarkBitTricksDividePowerOfTwo-4            2000000000               0.72 ns/op
BenchmarkBitTricksDivideNonPowerOfTwo-4         2000000000               1.09 ns/op
BenchmarkBitTricksShift-4                       2000000000               0.52 ns/op

Generated using go version go1.8.1 darwin/amd64

These benchmarks look at some micro optimizations that can be performed when doing division or
modulo division. The first three benchmarks show the overhead of doing modulus division and how
we can replace modulus division by a power of two with a bitwise and, which is a faster operation.
Likewise the last three benchmarks show the overhead of division and how we can improve the speed
of division by a power of two by performing a right shift.

### Buffered vs Synchronous Channel

`buffered_vs_unbuffered_channel_test.go`

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

`channel_vs_ring_buffer_test.go`

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkChannelSPSC    | 10000000  |        140 ns/op  |     8 B/op |    1 allocs/op
BenchmarkRingBufferSPSC | 20000000  |        106 ns/op  |     8 B/op |    1 allocs/op
BenchmarkChannelSPMC    |  5000000  |        369 ns/op  |     8 B/op |    1 allocs/op
BenchmarkRingBufferSPMC |  5000000  |        341 ns/op  |     8 B/op |    1 allocs/op
BenchmarkChannelMPSC    |  3000000  |        417 ns/op  |     8 B/op |    1 allocs/op
BenchmarkRingBufferMPSC |   500000  |       8387 ns/op  |     8 B/op |    1 allocs/op
BenchmarkChannelMPMC    |       20  |   70112786 ns/op  | 10532 B/op | 1031 allocs/op
BenchmarkRingBufferMPMC |        1  | 1228960979 ns/op  | 14256 B/op | 1015 allocs/op

Generated using go version go1.7.5 darwin/amd64

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

`defer_test.go`

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

### Function Call

`function_call_test.go`

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

`interface_conversion_test.go`

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkInterfaceNoConversion | 300000000 | 4.01 ns/op | 0 B/op | 0 allocs/op
BenchmarkInterfaceConversion   | 300000000 | 4.12 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at the overhead of converting a pointer to a struct to an interface when passing
it to a function which expects an interface. I was little surprised to find there is almost no
overhead.

### Mutex

`mutex_test.go`

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

`non_cryptographic_hash_functions_test.go`

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkHash32Fnv         | 20000000 |  72.1 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash32Fnva        | 20000000 |  70.4 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64Fnv         | 20000000 |  77.8 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64Fnva        | 20000000 |  68.9 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash32Crc         | 30000000 |  69.4 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64Crc         | 10000000 |   163 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash32Adler       | 30000000 |  39.0 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash32Xxhash      | 30000000 |  62.8 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64Xxhash      | 50000000 |  31.4 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash32Murmur3     | 30000000 |  53.6 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash128Murmur3    | 30000000 |  49.5 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64CityHash    | 50000000 |  28.9 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash128CityHash   | 20000000 |   109 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash32FarmHash    | 30000000 |  46.2 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64FarmHash    | 50000000 |  25.3 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash128FarmHash   | 50000000 |  37.3 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64SipHash     | 50000000 |  37.1 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash128SipHash    | 30000000 |  44.9 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64HighwayHash | 50000000 |  38.8 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash32SpookyHash  | 30000000 |  54.6 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash64SpookyHash  | 30000000 |  53.4 ns/op | 0 B/op | 0 allocs/op
BenchmarkHash128SpookyHash | 30000000 |  47.4 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

These benchmarks look at the speed of various non-cryptographic hash function implementations in Go.

### Pass By Value vs Reference

`pass_by_value_vs_reference_test.go`

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

`pool_test.go`

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

### Rand

`rand_test.go`

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkGlobalRand |  20000000 | 101 ns/op | 0 B/op | 0 allocs/op
BenchmarkLocalRand  | 200000000 |5.79 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

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

### Slice Initialization Append vs Index

`slice_intialization_append_vs_index_test.go`

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

`string_concatenation_test.go`

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

`type_assertion_test.go`

Benchmark Name|Iterations|Per-Iteration|Bytes Allocated per Operation|Allocations per Operation
----|----|----|----|----
BenchmarkTypeAssertion | 2000000000 | 0.97 ns/op | 0 B/op | 0 allocs/op

Generated using go version go1.7.5 darwin/amd64

This benchmark looks at the performance cost of a type assertion. I was a little surprised to find
it was so cheap.

### Write Bytes vs String

`write_bytes_vs_string_test.go`

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
