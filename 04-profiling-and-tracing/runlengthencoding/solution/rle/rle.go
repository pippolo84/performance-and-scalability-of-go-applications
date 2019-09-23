package rle

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

// Reader is an io.Reader that decode RLE input from a wrapped reader
type Reader struct {
	sym byte
	run byte
	rb  *bufio.Reader
}

// NewReader creates a new Reader wrapping r
func NewReader(r io.Reader) *Reader {
	z := new(Reader)
	z.rb = bufio.NewReader(r)
	return z
}

func (z *Reader) nextRun() (err error) {
	// read run length
	z.run, err = z.rb.ReadByte()

	// error while reading
	if err != nil {
		return err
	}

	// read run symbol
	z.sym, err = z.rb.ReadByte()

	// unexpected EOF
	if err == io.EOF {
		return errors.New("invalid rle stream")
	}

	// error while reading
	if err != nil {
		return err
	}

	return nil
}

// Read reads and decodes bytes from the wrapped stream, writing up to len(p) bytes into p
func (z *Reader) Read(p []byte) (n int, err error) {
	for len(p) > 0 {
		// read new data if needed
		if z.run == 0 {
			err = z.nextRun()

			// error or EOF while reading
			if err != nil {
				return n, err
			}
		}

		// write the symbol for z.run times
		for len(p) > 0 && z.run > 0 {
			p[0] = z.sym
			p = p[1:]
			n++
			z.run--
		}
	}

	return 0, nil
}

// WriteTo writes a decoded form of the wrapped input stream to w
func (z *Reader) WriteTo(w io.Writer) (n int64, err error) {
	wb := bufio.NewWriter(w)
	defer wb.Flush()

	for {
		// read next run (length and symbol)
		err = z.nextRun()

		// EOF
		if err == io.EOF {
			break
		}

		// error while reading
		if err != nil {
			return n, err
		}

		// repeat the symbol for run times
		for i := 0; i < int(z.run); i++ {
			err = wb.WriteByte(z.sym)
			if err != nil {
				return n, err
			}
			n++
		}
	}

	return n, nil
}

// Writer is an io.Writer that encode RLE input to a wrapped writer
type Writer struct {
	sym  byte
	prev byte
	run  byte
	wb   *bufio.Writer
}

// NewWriter creates a new Writer wrapping w
func NewWriter(w io.Writer) *Writer {
	z := new(Writer)
	z.wb = bufio.NewWriter(w)
	return z
}

func (z *Writer) write(rb *bufio.Reader) (n int, err error) {
	defer z.wb.Flush()

	// read the first byte from the stream
	z.sym, err = rb.ReadByte()

	// end of stream
	if err == io.EOF {
		return n, nil
	}

	// error while reading
	if err != nil {
		return n, err
	}

	// update bytes read
	n++

	// init encoder status with the first byte
	z.prev = z.sym
	z.run = 1

	// encode input stream until EOF or error
	for {
		// read next symbol
		z.sym, err = rb.ReadByte()

		// end of stream
		if err == io.EOF {
			break
		}

		// error while reading
		if err != nil {
			return n, err
		}

		// update bytes read
		n++

		// read symbol is different or we reached the run size limit
		if z.sym != z.prev || (z.run == 255 && z.sym == z.prev) {
			// flush run length and run symbol
			z.wb.WriteByte(z.run)
			z.wb.WriteByte(z.prev)

			// reset encoder status
			z.prev = z.sym
			z.run = 1
		} else {
			// increment the run length
			z.run++
		}
	}

	// flush remaining data
	z.wb.WriteByte(z.run)
	z.wb.WriteByte(z.prev)

	return n, nil
}

// Write writes and encodes bytes to the wrapped stream, reading up to len(p) bytes from p
func (z *Writer) Write(p []byte) (n int, err error) {
	return z.write(bufio.NewReader(bytes.NewReader(p)))
}

// ReadFrom reads input from r and writes a RLE encoded stream to the wrapped writer
func (z *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	written, err := z.write(bufio.NewReader(r))
	return int64(written), err
}
