package transport

import (
	"encoding/gob"
	"errors"
	"io"
	"log"
	"net"

	"github.com/zpnst/c2t/internal/protocol"
)

type TCPTransportOpts struct {
	InstanceAddr string
	Encoding     protocol.Encoding
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	recvChan chan protocol.Message
	Peers    map[string]protocol.Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		recvChan:         make(chan protocol.Message),
		Peers:            make(map[string]protocol.Peer),
	}
}

func (t TCPTransport) Peer(peerAddr string) protocol.Peer {
	return t.Peers[peerAddr]
}

func (t TCPTransport) Addr() string {
	return t.InstanceAddr
}

func (t TCPTransport) Close() error {
	return t.listener.Close()
}

func (t TCPTransport) Consume() <-chan protocol.Message {
	return t.recvChan
}

func (t TCPTransport) SendAnswer(src string, b any) error {
	p := t.Peer(src)
	ans := protocol.Message{
		From: t.Addr(),
		Body: b,
	}
	return gob.NewEncoder(p).Encode(ans)
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	if t.listener, err = net.Listen("tcp", t.InstanceAddr); err != nil {
		return err
	} else {
		log.Printf("[%s] :: instance is up\n", t.InstanceAddr)
		go t.AcceptLoop()
	}
	return nil
}

func (t TCPTransport) AcceptLoop() {
	for {
		if conn, err := t.listener.Accept(); errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
			log.Printf("[%s] :: connection rejected by peer: %s\n",
				t.InstanceAddr, conn.RemoteAddr().String())
			return
		} else if err != nil {
			log.Println(err)
			return
		} else {
			t.Peers[conn.RemoteAddr().String()] = TCPPerr{
				Conn: conn,
			}
			go t.HandleConn(conn)
		}
	}
}

func (t TCPTransport) HandleConn(conn net.Conn) {
	var err error

	defer func() {
		log.Printf("[%s] :: dropping peer connection: %s\n", t.Addr(), err)
		conn.Close()
	}()

	var m protocol.Message = protocol.Message{
		From: conn.RemoteAddr().String(),
	}
	for {
		if err = t.Encoding.Decode(conn, &m); err != nil {
			return
		}
		t.recvChan <- m
	}
}
