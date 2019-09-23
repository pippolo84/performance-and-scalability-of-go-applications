package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"performance-and-scalability-of-go-applications/04-profiling-and-tracing/runlengthencoding/solution/rle"
	"runtime/pprof"
)

func encodeFile(f string) (n int64, err error) {
	in, err := os.Open(f)
	if err != nil {
		return n, err
	}
	defer in.Close()

	encoded, err := os.Create("encoded.rle")
	if err != nil {
		return n, err
	}
	defer encoded.Close()

	rleWriter := rle.NewWriter(encoded)

	return io.Copy(rleWriter, in)
}

func decodeFile(f string) (n int64, err error) {
	in, err := os.Open(f)
	if err != nil {
		return n, err
	}
	defer in.Close()

	decoded, err := os.Create("decoded.out")
	if err != nil {
		return n, err
	}
	defer decoded.Close()

	rleReader := rle.NewReader(in)

	return io.Copy(decoded, rleReader)
}

func main() {
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(flag.Args()) < 2 {
		fmt.Println("missing arguments")
		return
	}

	switch flag.Args()[0] {
	case "e":
		_, err := encodeFile(flag.Args()[1])
		if err != nil {
			panic(err)
		}
	case "d":
		_, err := decodeFile(flag.Args()[1])
		if err != nil {
			panic(err)
		}
	default:
		fmt.Printf("unknown option %s\n", flag.Args()[0])
	}
}
