package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"performance-and-scalability-of-go-applications/04-profiling/runlengthencoding/rle"
	"runtime/pprof"
)

func encodeFile(f string) error {
	in, err := os.Open(f)
	if err != nil {
		return err
	}
	defer in.Close()

	encoded, err := rle.Encode(in)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("encoded.rle", encoded, 0644)
	if err != nil {
		return err
	}

	return nil
}

func decodeFile(f string) error {
	in, err := os.Open(f)
	if err != nil {
		return err
	}
	defer in.Close()

	var decoded bytes.Buffer
	err = rle.Decode(in, &decoded)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("decoded.out", decoded.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
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
		err := encodeFile(flag.Args()[1])
		if err != nil {
			panic(err)
		}
	case "d":
		err := decodeFile(flag.Args()[1])
		if err != nil {
			panic(err)
		}
	default:
		fmt.Printf("unknown option %s\n", flag.Args()[0])
	}
}
