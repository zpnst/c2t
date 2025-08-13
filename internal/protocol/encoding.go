package protocol

import (
	"encoding/gob"
	"io"
)

type Encoding interface {
	Encode(w io.Writer, m *Message) error
	Decode(r io.Reader, m *Message) error
}

type GOBEncoding struct{}

func NewGOBEncoding() GOBEncoding {
	return GOBEncoding{}
}

func (g GOBEncoding) Decode(r io.Reader, m *Message) error {
	return gob.NewDecoder(r).Decode(m)
}

func (g GOBEncoding) Encode(w io.Writer, m *Message) error {
	return gob.NewEncoder(w).Encode(m)
}
