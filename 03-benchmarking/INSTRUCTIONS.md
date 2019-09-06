# Insert and Remove

Starting from an empty sequence, insert N random integers into a sorted sequence, then remove those elements one by one, as determined by a random sequence of positions

## Example

### _Insert phase_

Suppose to generate the random sequence `4, 7, 2, 8`

your container should hold, after each step:

```
4
4, 7
2, 4, 7
2, 4, 7, 8
```

### _Remove phase_

Suppose to generate the random sequence `3, 0, 1, 0` (zero-based indexing)

your container should hold, after each step:

```
2, 4, 7
4, 7
4
```

N.B.: while generating random positions to remove elements, at each step you must pick a number in the range [0, container_length)

## Complete the exercise

Complete the functions `InsertRemoveSlice` and `InsertRemoveList` inside `insertremove/insertremove.go`.
Use a [slice](https://golang.org/ref/spec#Slice_types) in the former and a [list](https://golang.org/pkg/container/list/) in the latter.

Finally, write some benchmarks to compare their performance in `insertremove/insertremove_test.go`.

## Running the benchmarks

To run the benchmarks run the following command from within the `03-benchmarking/insertremove` directory:

`go test -v -run=^$ -bench=.`

While you try to improve your code, you should use [benchstat](https://godoc.org/golang.org/x/perf/cmd/benchstat) to do a before/after comparison.

---

Source: [Are lists evil?](https://isocpp.org/blog/2014/06/stroustrup-lists)
