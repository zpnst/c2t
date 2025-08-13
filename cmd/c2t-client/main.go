package main

import (
	"log"

	"github.com/zpnst/c2t/cmd/c2t-client/transport"
	"github.com/zpnst/c2t/internal/protocol"
)

func makeClient(iaddr string) *Client {
	topts := transport.TCPTransportOpts{
		InstanceAddr: iaddr,
		Encoding:     protocol.NewGOBEncoding(),
	}
	t := transport.NewTCPTransport(topts)

	copts := ClientOpts{
		Transport: t,
	}

	c := NewClient(copts)
	return c
}

func main() {
	// name := os.Args[1]
	// deviceId, _ := strconv.Atoi(os.Args[2])
	// serializer := signal.NewSerializer()
	// user := signal.NewUser(name, uint32(deviceId), serializer)
	// log.Println(user)

	c := makeClient(":4042")

	if err := c.Init(); err != nil {
		log.Panicln(err)
	}

	m := protocol.Message{
		Body: protocol.MessageSignUp{
			Name:     "Alice",
			DeviceID: 2456,
		},
	}

	if err := c.Transport.Send(m); err != nil {
		log.Panicln(err)
	}
}
