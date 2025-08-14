package transport

import (
	"net"

	c2t "github.com/zpnst/c2t/internal/protocol"
)

type TCPTransportOpts struct {
	Encoding     c2t.Encoding
	InstanceAddr string
}

type TCPTransport struct {
	TCPTransportOpts
	Conn net.Conn
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) EncodeMessage(m *c2t.Message) error {
	return t.Encoding.Encode(t.Conn, m)
}

func (t *TCPTransport) DecodeAnswer(m *c2t.Message) error {
	return t.Encoding.Decode(t.Conn, m)
}

func (t *TCPTransport) DialInstance() error {
	var err error
	t.Conn, err = net.Dial("tcp", t.InstanceAddr)
	if err != nil {
		return err
	}
	return nil
}
