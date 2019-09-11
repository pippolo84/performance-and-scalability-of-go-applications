# Gaussian blur

[Gaussian blur](https://en.wikipedia.org/wiki/Gaussian_blur) is a widely used effect in graphics software, typically to reduce image noise and detail. The visual effect of this blurring technique is a smooth blur resembling that of viewing the image through a translucent screen.

Mathematically, applying a Gaussian blur to an image is the same as [convolving](https://www.programming-techniques.com/2013/02/calculating-convolution-of-image-with-c_2.html) the image with a [Gaussian filter](https://www.geeksforgeeks.org/gaussian-filter-generation-c/)

## Example

Build the `gaussianblur` executable. From the `gaussianblur` folder

`go build`

Then, apply a gaussian blur effect to the test image `img/Lenna.png`

`./gaussianblur -in ../../img/Lenna.png -out blurred.png`

Use the `-h` flag on the executable to see all the available options

## Complete the exercise

Read the source code, try to understand how a 2D convolution with a gaussian kernel works.
Time and profile the program.

Is it possible to exploit Go concurrency support to make it faster? And, if so, to which extent?

Experiment with different solutions and try to point out the advantages and the limitations of each one.

---

Inspired by: [mandelbrot](https://github.com/campoy/mandelbrot)
