package main

import (
	"flag"
	"image/png"
	"log"
	"os"
	"performance-and-scalability-of-go-applications/05-scheduler/gaussianblur/solution/gblur"
	"runtime/trace"
)

func main() {
	input := flag.String("in", "lenna.png", "name of the input PNG file")
	output := flag.String("out", "blurred.png", "name of the output PNG file")
	kernelSz := flag.Int("size", 13, "size of the Gaussian kernel")
	sigma := flag.Float64("sigma", 2.5, "standard deviation")
	traceFile := flag.String("trace", "", "trace file name")
	flag.Parse()

	if *traceFile != "" {
		f, err := os.Create(*traceFile)
		if err != nil {
			log.Fatalf("could not create trace output file %q: %v", *traceFile, err)
		}
		if err := trace.Start(f); err != nil {
			log.Fatalf("could not start trace: %v", err)
		}
		defer func() {
			trace.Stop()
		}()
	}

	// check requested kernel size
	if *kernelSz%2 == 0 {
		log.Fatal("kernel size must be odd")
	}

	// read source image
	in, err := os.Open(*input)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	// Decode PNG source image
	img, err := png.Decode(in)
	if err != nil {
		log.Fatal(err)
	}

	// apply gaussian filter to blur image
	filter := gblur.GaussianKernel(*kernelSz, *sigma)
	blurredImg := gblur.Convolve(*gblur.ToNRGBA(img), filter)

	// save transformed image encoding it in PNG format
	out, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := png.Encode(out, blurredImg); err != nil {
		log.Fatal(err)
	}
}
