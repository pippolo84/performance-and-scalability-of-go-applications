package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"performance-and-scalability-of-go-applications/04-profiling/runlengthencoding/rle"
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

	err = ioutil.WriteFile("decoded.rle", decoded.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("syntax: %s e|d file\n", os.Args[0])
		return
	}

	switch os.Args[1] {
	case "e":
		err := encodeFile(os.Args[2])
		if err != nil {
			panic(err)
		}
	case "d":
		err := decodeFile(os.Args[2])
		if err != nil {
			panic(err)
		}
	default:
		fmt.Printf("unknown option %s\n", os.Args[1])
	}
}
