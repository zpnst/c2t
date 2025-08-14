package main

import (
	"log"
	"strconv"
	"time"

	"github.com/zpnst/c2t/cmd/c2t-client/transport"
	c2t "github.com/zpnst/c2t/internal/protocol"
	signal "github.com/zpnst/libsignal-protocol-go/protocol"
	"github.com/zpnst/libsignal-protocol-go/serialize"
)

var (
	c *Client
)

func makeClient(iaddr string) *Client {
	topts := transport.TCPTransportOpts{
		InstanceAddr: iaddr,
		Encoding:     c2t.NewGOBEncoding(),
	}
	t := transport.NewTCPTransport(topts)

	copts := ClientOpts{
		Transport: t,
	}

	c := NewClient(copts)
	return c
}

func alice() {
	name := "Alice"
	deviceId, _ := strconv.Atoi("21456789")
	serializer := serialize.NewJSONSerializer()
	alice := NewClientSession(name, uint32(deviceId), serializer)

	dstName := "Bob"

	bobsBundle, err := c.SendGetBundle(dstName)
	if err != nil {
		log.Println("Alice ::", err)
	} else {
		log.Println("Alice :: get budle ok!")
	}
	alice.BuildSession(signal.NewSignalAddress(dstName, bobsBundle.DeviceID), serializer)

	log.Println(bobsBundle.DeviceID)
	log.Println(bobsBundle.SignedPreKeyPublic)
	log.Println(bobsBundle.SignedPreKeySignature)
}

func bob() {
	name := "Bob"
	password := "12345"
	deviceId, _ := strconv.Atoi("21456789")
	serializer := serialize.NewJSONSerializer()

	bob := NewClientSession(name, uint32(deviceId), serializer)

	if err := c.SendSignUp(name, password); err != nil {
		log.Println("Bob ::", err)
	} else {
		log.Println("Bob :: sign up ok!")
	}

	time.Sleep(time.Second * 3)

	bobsRawBundle := c2t.NewRawBundle(
		bob.RegistrationID,
		bob.DeviceID,
		bob.PreKeys[0].ID().Value,
		bob.SignedPreKey.ID(),
		bob.PreKeys[0].KeyPair().PublicKey().PublicKey(),
		bob.SignedPreKey.KeyPair().PublicKey().PublicKey(),
		bob.IdentityKeyPair.PublicKey().PublicKey().PublicKey(),
		bob.SignedPreKey.Signature(),
	)

	if err := c.SendSetBundle(name, *bobsRawBundle); err != nil {
		log.Println("Bob ::", err)
	} else {
		log.Println("Bob :: set bundle ok!")
	}

	log.Println(bob.DeviceID)
	log.Println(bob.SignedPreKey.KeyPair().PublicKey().PublicKey())
	log.Println(bob.SignedPreKey.Signature())

}

func main() {
	c = makeClient(":4042")
	if err := c.Init(); err != nil {
		log.Panicln(err)
	}
	bob()
}
