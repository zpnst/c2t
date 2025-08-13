package transport

import (
	"log"
	"net"

	"github.com/zpnst/c2t/internal/protocol"
)

type TCPTransportOpts struct {
	Encoding     protocol.Encoding
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

func (t *TCPTransport) DialInstance() error {
	var err error
	t.Conn, err = net.Dial("tcp", t.InstanceAddr)
	if err != nil {
		return err
	}
	return nil
}

func (t *TCPTransport) WaitForInstance() error {
	var m protocol.Message
	if err := t.Encoding.Decode(t.Conn, &m); err != nil {
		return err
	}
	log.Println(m)
	return nil
}

func (t TCPTransport) Send(m protocol.Message) error {
	if err := t.Encoding.Encode(t.Conn, &m); err != nil {
		return err
	}
	return t.WaitForInstance()
}
