package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/zpnst/c2t/cmd/c2t-client/transport"
	c2t "github.com/zpnst/c2t/internal/protocol"
	"github.com/zpnst/libsignal-protocol-go/protocol"
	signal "github.com/zpnst/libsignal-protocol-go/protocol"
	"github.com/zpnst/libsignal-protocol-go/serialize"
	"github.com/zpnst/libsignal-protocol-go/session"
)

var (
	c *Client

	bob   *ClientSession
	alice *ClientSession
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

func aliceF() {
	name := "Alice"
	deviceId, _ := strconv.Atoi("21456789")
	serializer := serialize.NewJSONSerializer()
	alice = NewClientSession(name, uint32(deviceId), serializer)

	dstName := "Bob"

	bobsRawBundle, err := c.SendGetBundle(dstName)
	if err != nil {
		log.Panicln("Alice ::", err)
	}
	log.Println("Alice :: get budle ok!")

	bobsBundle := bobsRawBundle.ToSignalBundle()
	bobsAddress := signal.NewSignalAddress(dstName, bobsBundle.DeviceID())
	alice.BuildSession(bobsAddress, serializer)

	if err := alice.SessionBuilder.ProcessBundle(bobsBundle); err != nil {
		log.Panicln("Alice :: unable to process retrieved prekey bundle")
	}
	log.Println("Alice :: retrieved prekey bundle succesfully processed")

	plaintextMessage := []byte("Hello, Bob! I am Alice :)")
	log.Println("Alice :: plaintext message: ", string(plaintextMessage))
	sessionCipher := session.NewCipher(alice.SessionBuilder, bobsAddress)
	message, err := sessionCipher.Encrypt(plaintextMessage)
	if err != nil {
		log.Panicln("Alice :: unable to encrypt message: ", err)

	}
	log.Println("Alice :: encrypted message sha-256 hash: ", fmt.Sprintf("0x%X", sha256.Sum256(message.Serialize())))

	time.Sleep(time.Second * 1)

	if err := c.SendPostEncryptedMessage(dstName, c2t.EncryptedMessage{
		FromName:     name,
		FromDeviceID: uint32(deviceId),
		Message:      message.Serialize(),
	}); err != nil {
		log.Panicln("Alice :: unable to post encrypt message: ", err)
	}
	log.Println("Alice :: post encrypt message ok: ")

}

func bobF() {
	name := "Bob"
	password := "12345"
	deviceId, _ := strconv.Atoi("21456789")
	serializer := serialize.NewJSONSerializer()

	bob = NewClientSession(name, uint32(deviceId), serializer)

	if err := c.SendSignUp(name, password); err != nil {
		log.Panicln("Bob ::", err)
	}
	log.Println("Bob :: sign up ok!")

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
		log.Panicln("Bob ::", err)
	}
	log.Println("Bob :: set bundle ok!")

}

func bob2F() {
	name := "Bob"
	serializer := serialize.NewJSONSerializer()

	encMsg, err := c.SendFetchEncryptedMessage(name)
	if err != nil {
		log.Panicln("Bob ::", err)
	}
	log.Println("Bob :: enc msg ok!", encMsg[0].FromName, encMsg[0].FromDeviceID)
	log.Println("Bob :: encrypted message sha-256 hash: ", fmt.Sprintf("0x%X", sha256.Sum256(encMsg[0].Message)))

	aliceAddr := signal.NewSignalAddress(encMsg[0].FromName, encMsg[0].FromDeviceID)

	bob.BuildSession(aliceAddr, serializer)

	receivedMessage, err := protocol.NewPreKeySignalMessageFromBytes(encMsg[0].Message, serializer.PreKeySignalMessage, serializer.SignalMessage)
	if err != nil {
		log.Panicln("Bob ::", err)
	}

	unsignedPreKeyID, err := bob.SessionBuilder.Process(receivedMessage)
	if err != nil {
		log.Panicln("Bob :: unable to process prekeysignal message: ", err)
	}
	log.Println("Bob :: got PreKeyID: ", unsignedPreKeyID)

	bobSessionCipher := session.NewCipher(bob.SessionBuilder, aliceAddr)
	msg, err := bobSessionCipher.Decrypt(receivedMessage.WhisperMessage())
	if err != nil {
		log.Panicln("Bob :: unable to decrypt message: ", err)
	}
	log.Println("Bob :: ecrypted message: ", string(msg))
}

func main() {
	c = makeClient(":4042")
	if err := c.Init(); err != nil {
		log.Panicln(err)
	}
	bobF()
	time.Sleep(time.Second * 1)
	aliceF()
	time.Sleep(time.Second * 1)
	bob2F()
	time.Sleep(time.Second * 1)

	// if os.Args[1] == "alice" {
	// 	alice()
	// } else if os.Args[1] == "bob" {
	// 	bob()
	// }
}
