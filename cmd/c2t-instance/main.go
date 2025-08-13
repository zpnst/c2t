package main

import (
	"log"

	"github.com/zpnst/c2t/cmd/c2t-instance/database"
	"github.com/zpnst/c2t/cmd/c2t-instance/transport"
	"github.com/zpnst/c2t/internal/protocol"
)

func makeInstance(listenAddr string) *Instance {
	topts := transport.TCPTransportOpts{
		InstanceAddr: listenAddr,
		Encoding:     protocol.NewGOBEncoding(),
	}

	iopts := InstanceOpts{
		Database:  database.NewInMemoryDatabase(),
		Transport: transport.NewTCPTransport(topts),
	}

	i := NewInstance(iopts)
	return i
}

func main() {
	i := makeInstance(":4042")
	log.Fatal(i.Run())
}
