package gblur

import (
	"image"
	"image/color"
	"math"
	"sync"
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

// ConvolvePixel apply the specified filter to the pixel at coordinates (x, y) (sequential)
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

// ConvolvePixel apply the specified filter to the pixel at coordinates (x, y) (one goroutine per kernel row)
// func ConvolvePixel(img image.NRGBA, x, y int, filter [][]float64) color.NRGBA {
// 	var r, g, b, a atomic.Float64

// 	var wg sync.WaitGroup

// 	delta := len(filter) / 2
// 	for k := -delta; k <= delta; k++ {
// 		wg.Add(1)
// 		go func(k int) {
// 			defer wg.Done()
// 			for j := -delta; j <= delta; j++ {
// 				// reflect index to avoid errors at image borders
// 				boundedX := BoundReflection(x-j, img.Bounds().Min.X, img.Bounds().Max.X)
// 				boundedY := BoundReflection(y-k, img.Bounds().Min.Y, img.Bounds().Max.Y)

// 				color := img.NRGBAAt(boundedX, boundedY)
// 				r.Add(float64(color.R) * filter[j+delta][k+delta])
// 				g.Add(float64(color.G) * filter[j+delta][k+delta])
// 				b.Add(float64(color.B) * filter[j+delta][k+delta])
// 				a.Add(float64(color.A) * filter[j+delta][k+delta])
// 			}
// 		}(k)
// 	}

// 	wg.Wait()

// 	return color.NRGBA{
// 		R: uint8(r.Load()),
// 		G: uint8(g.Load()),
// 		B: uint8(b.Load()),
// 		A: uint8(a.Load()),
// 	}
// }

// Convolve apply the specified filter to the image with a 2D convolution (sequential)
// func Convolve(img image.NRGBA, filter [][]float64) *image.NRGBA {
// 	dst := image.NewNRGBA(img.Bounds())

// 	for y := 0; y < img.Bounds().Max.Y; y++ {
// 		for x := 0; x < img.Bounds().Max.X; x++ {
// 			dst.SetNRGBA(x, y, ConvolvePixel(img, x, y, filter))
// 		}
// 	}

// 	return dst
// }

// Convolve apply the specified filter to the image with a 2D convolution (One gouroutine per CPU)
func Convolve(img image.NRGBA, filter [][]float64) *image.NRGBA {
	dst := image.NewNRGBA(img.Bounds())

	const nGoroutines = 8
	stride := img.Bounds().Max.Y / nGoroutines

	var wg sync.WaitGroup
	for i := 0; i < nGoroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			yMax := (i + 1) * stride
			if i == nGoroutines-1 {
				yMax = img.Bounds().Max.Y
			}
			for y := i * stride; y < yMax; y++ {
				for x := 0; x < img.Bounds().Max.X; x++ {
					dst.SetNRGBA(x, y, ConvolvePixel(img, x, y, filter))
				}
			}
		}(i)
	}

	wg.Wait()

	return dst
}

// Convolve apply the specified filter to the image with a 2D convolution (One gouroutine per row)
// func Convolve(img image.NRGBA, filter [][]float64) *image.NRGBA {
// 	dst := image.NewNRGBA(img.Bounds())

// 	var wg sync.WaitGroup
// 	for y := 0; y < img.Bounds().Max.Y; y++ {
// 		wg.Add(1)
// 		go func(y int) {
// 			defer wg.Done()

// 			for x := 0; x < img.Bounds().Max.X; x++ {
// 				dst.SetNRGBA(x, y, ConvolvePixel(img, x, y, filter))
// 			}
// 		}(y)
// 	}

// 	wg.Wait()

// 	return dst
// }

// Convolve apply the specified filter to the image with a 2D convolution (One gouroutine per pixel)
// func Convolve(img image.NRGBA, filter [][]float64) *image.NRGBA {
// 	dst := image.NewNRGBA(img.Bounds())

// 	var wg sync.WaitGroup
// 	for y := 0; y < img.Bounds().Max.Y; y++ {
// 		for x := 0; x < img.Bounds().Max.X; x++ {
// 			wg.Add(1)
// 			go func(x, y int) {
// 				defer wg.Done()
// 				dst.SetNRGBA(x, y, ConvolvePixel(img, x, y, filter))
// 			}(x, y)
// 		}
// 	}
// 	wg.Wait()

// 	return dst
// }

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
