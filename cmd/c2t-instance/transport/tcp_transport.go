package transport

import (
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

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	if t.listener, err = net.Listen("tcp", t.InstanceAddr); err != nil {
		return err
	} else {
		log.Printf("[%s] :: node is up\n", t.InstanceAddr)
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
			go func() {
				var m protocol.Message = protocol.Message{
					From: conn.RemoteAddr().String(),
				}
				t.Encoding.Decode(conn, &m)
				t.recvChan <- m
			}()
		}
	}
}
