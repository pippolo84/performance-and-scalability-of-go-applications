# Run Length Encoding

[Run-length encoding](https://en.wikipedia.org/wiki/Run-length_encoding) is a simple form of data compression, where a stream of bytes is stored as a sequence of <data_value, data_count>.

As an example, consider the following byte stream:

`... 1 1 1 1 1 15 15 15 15 7 7 7 7 ...`

The RLE encoding of that stream chunk would be:

`... 5-1 4-15 4-7 ...`

That is, five times the value `1`, 4 four times the value `15`, four times the value `7` and so on

The package `rle` contains an implementation of a very basic RLE encoder and decoder.
Plus, in `main.go` there is a sample program to demonstrate the use of the `rle` package on files.

In this exercise you should read the source code I wrote, profile it and optimize it.

## Example

As a meaningful example, download [La Divina Commedia](https://www.gutenberg.org/files/1012/1012-0.txt) from [Project Gutenberg](https://www.gutenberg.org/) and save it locally as `divina-commedia.txt`.

Then, build the `runlengthencoding` executable. From the `runlengthencoding` folder:

`go build`

### Encode a file

To encode a file:

`./runlengthencoding e divina-commedia.txt`

This will create the `encoded.rle` file

### Decode a file

To decode a file:

`./runlengthencoding d encoded.rle`

This will create the `decoded.out` file

### Check the process of encoding and decoding

To verify that the starting file and the decoded file contains the same data, compare their md5sum digest:

`md5sum divina-commedia.txt decoded.txt`

Inside the `rle` package you will also find some unit tests to check the correctness of your changes.

## Complete the exercise

Read the source code, try to launch the tests and use the program to encode and decode some files.
Then profile it and guess what can be optimizes.

Refer to the slides to see how to generate a profile and how to read it.

---

Source: [Dave Cheney - Two Go Programs, Three Different Profiling Techniques](https://www.youtube.com/watch?v=nok0aYiGiYA)
