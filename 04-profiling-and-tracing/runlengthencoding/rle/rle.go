package rle

import (
	"io"
)

// Encode returns a RLE-encoded slice of bytes from the input stream r
// or an error if something goes wrong
func Encode(r io.Reader) ([]byte, error) {
	var encoded []byte
	var buf [1]byte
	var err error

	// read the first byte from the stream
	_, err = r.Read(buf[:])

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
		_, err = r.Read(buf[:])

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

// Decode takes the RLE-encoded input stream and decodes it returning data in a slice
// or an error if something goes wrong
func Decode(r io.Reader) ([]byte, error) {
	var decoded []byte

	for {
		var sym byte
		var run byte

		var buf [1]byte

		// read run symbol
		_, err := r.Read(buf[:])

		// end of stream
		if err == io.EOF {
			break
		}

		// error while reading
		if err != nil {
			return nil, err
		}

		run = buf[0]

		// read run length
		_, err = r.Read(buf[:])

		// error while reading (end of stream here is an error, too)
		if err != nil {
			return nil, err
		}

		sym = buf[0]

		// repeat the symbol for run times
		for i := 0; i < int(run); i++ {
			decoded = append(decoded, sym)
		}
	}

	return decoded, nil
}
