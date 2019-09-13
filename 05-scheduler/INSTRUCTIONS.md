# Scheduling exercises

## Gaussian blur

[Gaussian blur](https://en.wikipedia.org/wiki/Gaussian_blur) is a widely used effect in graphics software, typically to reduce image noise and detail. The visual effect of this blurring technique is a smooth blur resembling that of viewing the image through a translucent screen.

Mathematically, applying a Gaussian blur to an image is the same as [convolving](https://www.programming-techniques.com/2013/02/calculating-convolution-of-image-with-c_2.html) the image with a [Gaussian filter](https://www.geeksforgeeks.org/gaussian-filter-generation-c/)

### Example

Build the `gaussianblur` executable. From the `gaussianblur` folder

`go build`

Then, apply a gaussian blur effect to the test image `img/Lenna.png`

`./gaussianblur -in ../../img/Lenna.png -out blurred.png`

Use the `-h` flag on the executable to see all the available options

### Complete the exercise

Read the source code, try to understand how a 2D convolution with a gaussian kernel works.
Time and profile the program.

Is it possible to exploit Go concurrency support to make it faster? And, if so, to which extent?

Experiment with different solutions and try to point out the advantages and the limitations of each one.

Inspired by: [mandelbrot](https://github.com/campoy/mandelbrot)

---

## Word Line Count

We want to analyze multiple text files downloaded from the web, countint the total number of lines where a certain word appears, matching it case insensitive.

For the sake of simplicity, we simulate it downloading the same text file multiple times, pretending it is different every time we download it.

### Example

Build the `wordlinecount` executable. From the `wordlinecount` folder

`go build`

To calculate the number of lines with the word `This` in them

`./wordlinecount -word This`

The program will issue 8 GET requests to https://www.gutenberg.org/files/16/16-0.txt, downloading the text files and calculating the total count

### Complete the exercise

Read the source code, try to understand if the code may benefit from concurrency and/or parallelization
Time, profile and trace the program. You may want to use the `GOMAXPROCS` env var to tune the number of OS threads to use.

Is it possible to exploit Go concurrency support to make it faster? And, if so, to which extent?

Experiment with different solutions and try to point out the advantages and the limitations of each one.
