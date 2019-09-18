#TODO: heap allocations profiling (runtime.MemProfileRate, inuse_space and inuse_objects, alloc_space and alloc_objects, ecc)
see https://www.freecodecamp.org/news/how-i-investigated-memory-leaks-in-go-using-pprof-on-a-large-codebase-4bec4325e192/

* How memory for the user stack is allocated

Still through the heap (with various caches) but it is not accounted by the GC (the GC starts after a threshold is passed) -> need verification!!!

How can you say that it is allocated on the heap?

starting from src/runtime/proc.go, at line 3251

// create new goroutine
func newproc(siz int32, fn *funcval)

// switch to systemstack to call:

// Handle all intricacies of creating a goroutine, the stack allocation too!
func newproc1(fn *funcval, argp *uint8, narg int32, callergp *g, callerpc uintptr)

// call malg with a stack size of _StackMin (from src/runtime/stack.go we can see that _StackMin = 2048)
newg = malg(_StackMin)

// malg will call stackalloc with a size of round2(_StackSystem + stacksize)

// stackalloc allocates an n byte stack.
//go:systemstack
func stackalloc(n uint32) stack

  Here, we allocate stack using:

  - for small stacks ("Small stacks are allocated with a fixed-size free-list allocator.") normally we execute 1), less often 2):
    1) // get a stack from the global pool of stack. This is done because we cannot access the local stack cache during exitsyscall, procresize or gc
       x = stackpoolalloc(order)
    2) // get a reference to the stack cache (This cache is allocated with allocmcache, search for it in the code)
       c := thisg.m.mcache

       // get a stack
       x = c.stackcache[order].list
	   if x.ptr() == nil {
	       stackcacherefill(c, order)
		   x = c.stackcache[order].list
	   }
  - for large stacks ("If we need a stack of a bigger size, we fall back on allocating a dedicated span.")
    var s *mspan
    // Try to get a stack from the large stack cache.
    if !stackLarge.free[log2npage].isEmpty() {
        s = stackLarge.free[log2npage].first
        stackLarge.free[log2npage].remove(s)
    }

    if s == nil {
        // Allocate a new stack from the heap.
        s = mheap_.allocManual(npage, &memstats.stacks_inuse)
        if s == nil {
            throw("out of memory")
        }
        osStackAlloc(s)
        //...
    }

// stackcacherefill calls stackpoolalloc()

// allocmcache calls mheap_.cachealloc.alloc(), so the cache for small stacks is allocating from the heap
// stackalloc calls mheap_.allocManual, so it is allocating from the heap
// stackLarge.free (cache of larger stacks) is replenished in func stackfree(stk stack)

So we have stack span and heap span. The cache structures for both heap and stack allocations are similar (maybe equal?)

So, why stack allocations is cheaper than heap allocations?

From https://dave.cheney.net/2014/06/07/five-things-that-make-go-fast, in the Section "Escape Analysis":

"Because the numbers slice is only referenced inside Sum, the compiler will arrange to store the 100 integers for that slice on the stack, rather than the heap.
There is no need to garbage collect numbers, it is automatically freed when Sum returns.
Because escape analysis is performed at compile time, not run time, stack allocation will always be faster than heap allocation, no matter how efficient your garbage collector is."

for both stack and heap allocations we pay the following prices:
- bigger and frequent allocations defy the per P caches and optimizations in Go memory allocator
    - more contention while taking lock on shared allocator data structures
    - more code to be executed to search for free slots while going up the allocator data structures hierarchy (spans, arena, etc.)
    - syscall overhead when using sbrk/mmap to ask for more memory to the kernel

But for heap allocations, we pay ALSO these prices:
    - more allocation makes the GC starts often (or for longer? read about GC pacing and GOGC threshold)
        - STW pauses (worse latency)
        - more work in the concurrent marking phase, since the heap to scan is greater
        - though the marking is concurrent, we pay the price of the write barrier on
        - mark assist for goroutine that allocate a lot
        - more work for the scavenger ???

Plus, apart from the allocations that uses the per P caches, when you allocate on the heap, you need to lock, since the structures are shared

* Go for ETL and data manipulation

Quote from https://syslog.ravelin.com/further-dangers-of-large-heaps-in-go-7a267b57d487:

So what are the lessons to learn here? If you are using Go for data-processing then you either can’t have any long-term large heap allocations or you must ensure that they don’t contain any pointers. And this means no strings, no slices, no time.Time (it contains a pointer to a locale), no nothing with a hidden pointer in it. I hope to follow up with some blog posts about tricks I’ve used to do that.

* Why garbage collection?

tweet from Roberto Clapis (about Go vs Rust)

"It’s garbage collected [...] this harms performance" → just a little but grants that unreferenced memory is freed, which Rust doesn't grant. So a pure safe Rust executable can run for a while and then die for an OOM. Good luck debugging that.

My thoughts: remember that Go is primarily target to high-performance web service, so OOMs are to be particularly feared.

* Another example of optimization: ZAP package

- no interface{}
- no allocations

https://medium.com/@blanchon.vincent/go-how-zap-package-is-optimized-dbf72ef48f2d

* Cache false sharing

Inside Go runtime (cpu.go)

    // CacheLinePad is used to pad structs to avoid false sharing.
    type CacheLinePad struct{ _ [CacheLinePadSize]byte }

    // CacheLineSize is the CPU's assumed cache line size.
    // There is currently no runtime detection of the real cache line size
    // so we use the constant per GOARCH CacheLinePadSize as an approximation.
    var CacheLineSize uintptr = CacheLinePadSize

Note that, CacheLinePad is set 64 bytes for amd64, no runtime detection.
Try also to search for CacheLinePad to see some padding examples in the runtime

Cache false sharing is addressed in part 1 of scheduler series by Bill kennedy and in this talk https://www.youtube.com/watch?v=WDIkqP4JbkE

and in this post (regarding Go): https://medium.com/@genchilu/whats-false-sharing-and-how-to-solve-it-using-golang-as-example-ef978a305e10

Exercise on cache false sharing by Bill Kennedy: https://github.com/ardanlabs/gotraining/blob/master/topics/go/testing/benchmarks/falseshare/README.md

* Escape analysis

Post on Ardan Labs about Escape analysis flaws (test again with Go 1.13)
Go 1.13 release notes (escape analysis update):
```
The compiler has a new implementation of escape analysis that is more precise.
For most Go code should be an improvement (in other words, more Go variables and
expressions allocated on the stack instead of heap).
However, this increased precision may also break invalid code that happened to work before
(for example, code that violates the unsafe.Pointer safety rules).
If you notice any regressions that appear related, the old escape analysis pass can be re-enabled with
go build -gcflags=all=-newescape=false.
The option to use the old escape analysis will be removed in a future release.
```

* Channels

see this (with all comments and Kavya Joshi talk about channels): https://www.reddit.com/r/golang/comments/apv6gj/why_is_a_channel_faster_than_a_mutex_in_this_test/

* Bad Go series

https://philpearl.github.io/

* [Go 1.13] performance improvements

See https://docs.google.com/presentation/d/1RiZmupILuIQQ1Y-psDb1SzXNjCWh-I_-wagthdcwlq8/edit#slide=id.g604d13147b_0_88

nice examples for mid-stack inlining and sync.Pool with GC

See also this tweets https://twitter.com/bradfitz/status/1121551757971087360

* [Go 1.13] mid-stack inlining and sync.Mutex

AFAIK the mid stack inlining allowed this patch https://github.com/golang/go/commit/41cb0aedffdf4c5087de82710c4d016a3634b4ac to inline the fast-path
(note that the slow path has been "outlined" manually in a separate function `lockSlow`)

* [Go.1.13] sync.Pool

https://medium.com/a-journey-with-go/go-understand-the-design-of-sync-pool-2dde3024e277

see also this example: https://docs.google.com/presentation/d/1RiZmupILuIQQ1Y-psDb1SzXNjCWh-I_-wagthdcwlq8/edit#slide=id.g604d13147b_0_106
(it should show that the object won't be garbage collected?)

* [Go 1.13] defer

https://go-review.googlesource.com/c/go/+/171758/

* Bounds Checks

???

* Performance analysis methodology

Steps to take to correctly analize performance

* Exercises:

1) [Parallel Letter Frequencies](https://exercism.io/my/solutions/c8fdbd363cca40999a1e12c9fca17875)
2) [GraphBlog](https://syslog.ravelin.com/making-something-faster-56dd6b772b83)
(see also https://github.com/ardanlabs/gotraining/tree/master/topics/go/profiling/exercises/graphblog)
3) [Stream Elvis replace](https://github.com/ardanlabs/gotraining/tree/master/topics/go/profiling/memcpu)
4) [http/pprox](https://github.com/ardanlabs/gotraining/tree/master/topics/go/profiling/pprof)
5) [Profiling exercises with XML topics search](https://github.com/ardanlabs/gotraining/tree/master/topics/go/profiling/trace)
6) [Exercises with CPU bound and I/O bound applications](https://www.ardanlabs.com/blog/2018/12/scheduling-in-go-part3.html)
7) [Large Web Service](https://github.com/ardanlabs/gotraining/tree/master/topics/go/profiling/project)