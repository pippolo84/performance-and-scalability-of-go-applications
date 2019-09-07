package rle

import (
	"bufio"
	"io"
)

// Encode returns a RLE-encoded slice of bytes from the input stream r
// or an error if something goes wrong
func Encode(r io.Reader) ([]byte, error) {
	var encoded []byte
	var buf [1]byte
	var err error

	b := bufio.NewReader(r)

	// read the first byte from the stream
	_, err = b.Read(buf[:])

	// end of stream
	if err == io.EOF {
		return encoded, nil
	}

	// error while reading
	if err != nil {

		return nil, err
	}

	// init encoder status with the first byte
	prev := buf[0]
	run := 1

	// encode input stream until EOF or error
	for {
		// read next symbol
		_, err = b.Read(buf[:])

		// end of stream
		if err == io.EOF {
			break
		}

		// error while reading
		if err != nil {
			return nil, err
		}

		// read symbol is different or we reached the run size limit
		if buf[0] != prev || (run == 255 && buf[0] == prev) {
			// flush run length and run symbol
			encoded = append(encoded, byte(run))
			encoded = append(encoded, prev)

			// reset encoder status
			prev = buf[0]
			run = 1
		} else {
			// increment the run length
			run++
		}
	}

	// flush remaining data
	encoded = append(encoded, byte(run))
	encoded = append(encoded, prev)

	return encoded, nil
}

// Decode writes the RLE-encoded input stream r into w
// Returns an error if something goes wrong
func Decode(r io.Reader, w io.Writer) error {
	b := bufio.NewReader(r)

	for {
		var sym byte
		var run byte

		var buf [1]byte

		// read run symbol
		_, err := b.Read(buf[:])

		// end of stream
		if err == io.EOF {
			break
		}

		// error while reading
		if err != nil {
			return err
		}

		run = buf[0]

		// read run length
		_, err = b.Read(buf[:])

		// error while reading (end of stream here is an error, too)
		if err != nil {
			return err
		}

		sym = buf[0]

		// write the symbol for run times
		for i := 0; i < int(run); i++ {
			_, err := w.Write([]byte{sym})

			// error while writing
			if err != nil {
				return err
			}
		}
	}

	return nil
}
