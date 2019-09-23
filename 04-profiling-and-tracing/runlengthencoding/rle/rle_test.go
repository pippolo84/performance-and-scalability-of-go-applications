package rle

import (
	"bytes"
	"io"
	"testing"
)

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
