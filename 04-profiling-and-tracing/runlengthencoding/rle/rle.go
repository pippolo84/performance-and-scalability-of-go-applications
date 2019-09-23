package rle

import (
	"bytes"
	"errors"
	"io"
)

// Reader is an io.Reader that decode RLE input from a wrapped reader
type Reader struct {
	sym byte
	run byte
	r   io.Reader
}

// NewReader creates a new Reader wrapping r
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

func (z *Reader) nextRun() (err error) {
	var buf [1]byte

	// read run length
	n, err := z.r.Read(buf[:])

	// error while reading
	if n != 1 || err != nil {
		return err
	}

	z.run = buf[0]

	// read run symbol
	n, err = z.r.Read(buf[:])

	// unexpected EOF
	if err == io.EOF {
		return errors.New("invalid rle stream")
	}

	// error while reading
	if n != 1 || err != nil {
		return err
	}

	z.sym = buf[0]

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
			buf := []byte{z.sym}
			nwrite, err := w.Write(buf)
			if nwrite != 1 || err != nil {
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
	w    io.Writer
}

// NewWriter creates a new Writer wrapping w
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

func (z *Writer) write(r io.Reader) (n int, err error) {
	var buf [1]byte

	// read the first byte from the stream
	nRead, err := r.Read(buf[:])

	// end of stream
	if err == io.EOF {
		return n, nil
	}

	// error while reading
	if nRead != 1 || err != nil {
		return n, err
	}
	z.sym = buf[0]

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
		nRead, err := r.Read(buf[:])

		// end of stream
		if err == io.EOF {
			break
		}

		// error while reading
		if nRead != 1 || err != nil {
			return n, err
		}
		z.sym = buf[0]

		// update bytes read
		n++

		// read symbol is different or we reached the run size limit
		if z.sym != z.prev || (z.run == 255 && z.sym == z.prev) {
			wbuf := []byte{z.run, z.prev}
			// flush run length and run symbol
			nWrite, err := z.w.Write(wbuf)
			if nWrite != 2 || err != nil {
				return n, err
			}

			// reset encoder status
			z.prev = z.sym
			z.run = 1
		} else {
			// increment the run length
			z.run++
		}
	}

	// flush remaining data
	wbuf := []byte{z.run, z.prev}
	nWrite, err := z.w.Write(wbuf)
	if nWrite != 2 || err != nil {
		return n, err
	}

	return n, nil
}

// Write writes and encodes bytes to the wrapped stream, reading up to len(p) bytes from p
func (z *Writer) Write(p []byte) (n int, err error) {
	return z.write(bytes.NewReader(p))
}

// ReadFrom reads input from r and writes a RLE encoded stream to the wrapped writer
func (z *Writer) ReadFrom(r io.Reader) (n int64, err error) {
	written, err := z.write(r)
	return int64(written), err
}
