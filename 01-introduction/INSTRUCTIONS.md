# Isogram

Determine if a word or phrase is an [isogram](https://en.wikipedia.org/wiki/Isogram).

An isogram (also known as a "nonpattern word") is a word or phrase without a repeating letter, however spaces and hyphens are allowed to appear multiple times.

Examples of isograms:

- lumberjacks
- background
- downstream
- six-year-old

The word "isograms", however, is not an isogram, because the s repeats.

## Complete the exercise

Complete the function `IsIsogram` inside `isogram/isogram.go`. Do not change any other file. When you're done, try running the tests.

## Running the tests

To run the tests run the following command from within the `01-introduction/isogram` directory:

`go test -v`

## Running the benchmarks

To run the benchmarks run the following command from within the `01-introduction/isogram` directory:

`go test -v -run=^$ -bench=.`

---

Source: [exercism.io](https://exercism.io)