package salloc

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// StrcatStrAppend uses the + operator to append the string pieces
func StrcatStrAppend(id int, name string, t time.Time) string {
	s := fmt.Sprintf("%d", id)
	s += " "
	s += name
	s += " "
	s += t.String()
	return s
}

// StrcatStrBuilder uses a strings.Builder
func StrcatStrBuilder(id int, name string, t time.Time) string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("%d", id))
	s.WriteByte(' ')
	s.WriteString(name)
	s.WriteByte(' ')
	s.WriteString(t.String())

	return s.String()
}

// StrcatSprintf uses fmt.Sprintf
func StrcatSprintf(id int, name string, t time.Time) string {
	return fmt.Sprintf("%d %s %s", id, name, t.String())
}

// StrcatBytesBuffer uses a bytes.Buffer
func StrcatBytesBuffer(id int, name string, t time.Time) string {
	var b bytes.Buffer

	fmt.Fprintf(&b, "%d %s %s", id, name, t.String())

	return b.String()
}

// StrcatSlicePreallocated uses a slice of bytes with a starting capacity of 64 bytes
func StrcatSlicePreallocated(id int, name string, t time.Time) string {
	b := make([]byte, 0, 64)

	b = append(b, fmt.Sprintf("%d", id)...)
	b = append(b, ' ')
	b = append(b, name...)
	b = append(b, ' ')
	b = append(b, t.String()...)

	return string(b)
}

// StrcatByteSlice uses a slice of bytes initially empty
func StrcatByteSlice(id int, name string, t time.Time) string {
	var b []byte

	b = append(b, fmt.Sprintf("%d", id)...)
	b = append(b, ' ')
	b = append(b, name...)
	b = append(b, ' ')
	b = append(b, t.String()...)

	return string(b)
}
