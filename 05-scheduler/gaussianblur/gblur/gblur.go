package gblur

import (
	"image"
	"image/color"
	"math"
)

// GaussianKernel returns a gaussian filter of size sz x sz with sigma standard deviation
func GaussianKernel(sz int, sigma float64) [][]float64 {
	filter := make([][]float64, sz)
	for i := range filter {
		filter[i] = make([]float64, sz)
	}

	var sum float64

	// generate kernel
	s := 2 * sigma * sigma
	delta := sz / 2
	for x := -delta; x <= delta; x++ {
		for y := -delta; y <= delta; y++ {
			r := float64(x*x) + float64(y*y)
			filter[x+delta][y+delta] = math.Exp(-(r)/s) / (s * math.Pi)
			sum += filter[x+delta][y+delta]
		}
	}

	// normalize the kernel
	for i := range filter {
		for j := range filter[i] {
			filter[i][j] /= sum
		}
	}

	return filter
}

// BoundReflection returns a reflected index on the interval [min, max)
func BoundReflection(idx, min, max int) int {
	if idx < min {
		return -idx - 1
	}

	if idx >= max {
		return 2*max - idx - 1
	}

	return idx
}

// ConvolvePixel apply the specified filter to the pixel at coordinates (x, y)
func ConvolvePixel(img image.NRGBA, x, y int, filter [][]float64) color.NRGBA {
	var r, g, b, a float64

	delta := len(filter) / 2
	for k := -delta; k <= delta; k++ {
		for j := -delta; j <= delta; j++ {
			// reflect index to avoid errors at image borders
			boundedX := BoundReflection(x-j, img.Bounds().Min.X, img.Bounds().Max.X)
			boundedY := BoundReflection(y-k, img.Bounds().Min.Y, img.Bounds().Max.Y)

			color := img.NRGBAAt(boundedX, boundedY)
			r += float64(color.R) * filter[j+delta][k+delta]
			g += float64(color.G) * filter[j+delta][k+delta]
			b += float64(color.B) * filter[j+delta][k+delta]
			a += float64(color.A) * filter[j+delta][k+delta]
		}
	}

	return color.NRGBA{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: uint8(a),
	}
}

// Convolve apply the specified filter to the image with a 2D convolution
func Convolve(img image.NRGBA, filter [][]float64) *image.NRGBA {
	dst := image.NewNRGBA(img.Bounds())

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			dst.SetNRGBA(x, y, ConvolvePixel(img, x, y, filter))
		}
	}

	return dst
}

// ToNRGBA creates a copy of img converting its format in NRGBA (non-alpha-premultiplied 32-bit color)
func ToNRGBA(img image.Image) *image.NRGBA {
	dst := image.NewNRGBA(img.Bounds())

	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			dst.SetNRGBA(x, y, color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA))
		}
	}

	return dst
}
