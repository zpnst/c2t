package main

import (
	"log"

	"github.com/zpnst/c2t/cmd/c2t-instance/database"
	"github.com/zpnst/c2t/cmd/c2t-instance/transport"
	c2t "github.com/zpnst/c2t/internal/protocol"
)

func makeInstance(listenAddr string) *Instance {
	topts := transport.TCPTransportOpts{
		InstanceAddr: listenAddr,
		Encoding:     c2t.NewGOBEncoding(),
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
	if err := i.Run(); err != nil {
		log.Println(err)
	}
}
