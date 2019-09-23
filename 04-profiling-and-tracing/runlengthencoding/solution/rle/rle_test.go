package rle

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	// setup
	for _, count := range []int{4, 16, 64, 256, 512} {
		cmd := exec.Command(
			"dd",
			"if=/dev/urandom",
			fmt.Sprintf("of=./encode%dk.golden", count),
			"bs=1K",
			fmt.Sprintf("count=%d", count),
		)
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}

	// test
	res := m.Run()

	// teardown
	for _, pattern := range []string{"*.golden", "*.output"} {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			panic(err)
		}
		for _, match := range matches {
			err := os.Remove(match)
			if err != nil {
				panic(err)
			}
		}
	}

	os.Exit(res)
}

var readTestCases = []struct {
	descr  string
	input  []byte
	output []byte
}{
	{
		descr:  "empty stream",
		input:  []byte{},
		output: []byte{},
	},
	{
		descr:  "single run",
		input:  []byte{6, 10},
		output: []byte{10, 10, 10, 10, 10, 10},
	},
	{
		descr:  "multiple runs",
		input:  []byte{2, 10, 2, 20, 1, 30},
		output: []byte{10, 10, 20, 20, 30},
	},
}

func TestRead(t *testing.T) {
	for _, tc := range readTestCases {
		t.Run(tc.descr, func(t *testing.T) {
			rleReader := NewReader(bytes.NewReader(tc.input))

			// first partial read
			buf1 := make([]byte, len(tc.output)/2+1)
			_, err := rleReader.Read(buf1)
			if err != nil && err != io.EOF {
				t.Fatalf("rleReader.Read(...) returned unexpected error: %v", err)
			}

			// second partial read
			buf2 := make([]byte, len(tc.output)/2)
			_, err = rleReader.Read(buf2)
			if err != nil && err != io.EOF {
				t.Fatalf("rleReader.Read(...) returned unexpected error: %v", err)
			}

			buf := append(buf1, buf2...)
			if bytes.Compare(tc.output, buf[:len(tc.output)]) != 0 {
				t.Fatalf("got %v want %v", buf, tc.output)
			}
		})
	}
}

var writeTestCases = []struct {
	descr  string
	input  []byte
	output []byte
}{
	{
		descr:  "empty stream",
		input:  []byte{},
		output: []byte{},
	},
	{
		descr:  "single run",
		input:  []byte{10, 10, 10, 10, 10, 10},
		output: []byte{6, 10},
	},
	{
		descr:  "multiple runs",
		input:  []byte{10, 10, 20, 20, 30},
		output: []byte{2, 10, 2, 20, 1, 30},
	},
}

func TestWrite(t *testing.T) {
	for _, tc := range writeTestCases {
		t.Run(tc.descr, func(t *testing.T) {
			var buf bytes.Buffer
			rleWriter := NewWriter(&buf)

			n, err := rleWriter.Write(tc.input)
			if err != nil && err != io.EOF {
				t.Fatalf("rleWriter.Write(...) returned unexpected error: %v", err)
			}
			if n != len(tc.input) {
				t.Fatalf("rleWriter.Write(...) = %d, want %d", n, len(tc.input))
			}

			if bytes.Compare(tc.output, buf.Bytes()) != 0 {
				t.Fatalf("got %v want %v", buf.Bytes(), tc.output)
			}
		})
	}
}

var testCases = []struct {
	descr string
	seq   []byte
}{
	{
		descr: "empty sequence",
		seq:   []byte{},
	},
	{
		descr: "single symbol",
		seq:   []byte{15},
	},
	{
		descr: "all equal symbols",
		seq:   []byte{20, 20, 20, 20, 20, 20, 20, 20, 20, 20},
	},
	{
		descr: "two runs",
		seq:   []byte{20, 20, 20, 20, 20, 15, 15, 15},
	},
	{
		descr: "last run has a single symbol",
		seq:   []byte{20, 20, 18, 18, 18, 18, 15},
	},
	{
		descr: "random sequence",
		seq:   []byte{20, 20, 20, 20, 18, 10, 9, 8, 8, 8, 8, 8, 8, 250, 250, 250, 250, 250, 38, 38, 38, 41, 41},
	},
	{
		descr: "run with more than 255 elements",
		seq:   bytes.Repeat([]byte{20}, 2500),
	},
}

func TestReadFromWriteTo(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.descr, func(t *testing.T) {
			var b bytes.Buffer
			rleWriter := NewWriter(&b)
			n, err := rleWriter.ReadFrom(bytes.NewReader(tc.seq))
			if n != int64(len(tc.seq)) || err != nil {
				t.Fatalf("rleWriter.ReadFrom(...) got %d %v, want %d %v", n, err, len(tc.seq), nil)
			}

			var out bytes.Buffer
			rleReader := NewReader(&b)
			n, err = rleReader.WriteTo(&out)
			if n != int64(len(out.Bytes())) || err != nil {
				t.Fatalf("rleReader.WriteTo(...) got %d %v, want %d %v", n, err, len(out.Bytes()), nil)
			}

			if bytes.Compare(tc.seq, out.Bytes()) != 0 {
				t.Fatalf("got %v want %v", out.Bytes(), tc.seq)
			}
		})
	}
}

var encodeCases = []struct {
	descr string
	src   string
	dst   string
}{
	{
		descr: "4 KB",
		src:   "encode4k.golden",
		dst:   "decode4k.golden",
	},
	{
		descr: "16 KB",
		src:   "encode16k.golden",
		dst:   "decode16k.golden",
	},
	{
		descr: "64 KB",
		src:   "encode64k.golden",
		dst:   "decode64k.golden",
	},
	{
		descr: "256 KB",
		src:   "encode256k.golden",
		dst:   "decode256k.golden",
	},
	{
		descr: "512 KB",
		src:   "encode512k.golden",
		dst:   "decode512k.golden",
	},
}

func BenchmarkEncode(b *testing.B) {
	for _, bc := range encodeCases {
		b.Run(bc.descr, func(b *testing.B) {
			src, err := os.Open(bc.src)
			if err != nil {
				b.Fatal(err)
			}
			defer src.Close()

			dst, err := os.Create(bc.dst)
			if err != nil {
				b.Fatal(err)
			}
			defer dst.Close()

			rleWriter := NewWriter(dst)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rleWriter.ReadFrom(src)
			}
		})
	}
}

var decodeCases = []struct {
	descr string
	src   string
	dst   string
}{
	{
		descr: "4 KB",
		src:   "decode4k.golden",
		dst:   "decode4k.output",
	},
	{
		descr: "16 KB",
		src:   "decode16k.golden",
		dst:   "decode16k.output",
	},
	{
		descr: "64 KB",
		src:   "decode64k.golden",
		dst:   "decode64k.output",
	},
	{
		descr: "256 KB",
		src:   "decode256k.golden",
		dst:   "decode256k.output",
	},
	{
		descr: "512 KB",
		src:   "decode512k.golden",
		dst:   "decode512k.output",
	},
}

func BenchmarkDecode(b *testing.B) {
	for _, bc := range encodeCases {
		b.Run(bc.descr, func(b *testing.B) {
			src, err := os.Open(bc.src)
			if err != nil {
				b.Fatal(err)
			}
			defer src.Close()

			dst, err := os.Create(bc.dst)
			if err != nil {
				b.Fatal(err)
			}
			defer dst.Close()

			rleReader := NewReader(src)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rleReader.WriteTo(dst)
			}
		})
	}
}
