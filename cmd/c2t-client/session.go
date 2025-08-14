package main

import (
	"github.com/zpnst/c2t/cmd/c2t-instance/database"
	"github.com/zpnst/libsignal-protocol-go/groups"
	"github.com/zpnst/libsignal-protocol-go/groups/state/store"
	"github.com/zpnst/libsignal-protocol-go/keys/identity"
	"github.com/zpnst/libsignal-protocol-go/protocol"
	"github.com/zpnst/libsignal-protocol-go/serialize"
	"github.com/zpnst/libsignal-protocol-go/session"
	"github.com/zpnst/libsignal-protocol-go/state/record"
	stateStore "github.com/zpnst/libsignal-protocol-go/state/store"
	"github.com/zpnst/libsignal-protocol-go/util/keyhelper"
)

type ClientSession struct {
	Name     string
	DeviceID uint32
	Address  *protocol.SignalAddress

	IdentityKeyPair *identity.KeyPair
	RegistrationID  uint32

	PreKeys      []*record.PreKey
	SignedPreKey *record.SignedPreKey

	SessionStore      stateStore.Session
	PreKeyStore       stateStore.PreKey
	SignedPreKeyStore stateStore.SignedPreKey
	IdentityStore     stateStore.IdentityKey
	SenderKeyStore    store.SenderKey

	SessionBuilder *session.Builder
	GroupBuilder   *groups.SessionBuilder
}

func NewClientSession(name string, deviceID uint32, serializer *serialize.Serializer) *ClientSession {
	signalUser := &ClientSession{}

	signalUser.IdentityKeyPair, _ = keyhelper.GenerateIdentityKeyPair()
	signalUser.RegistrationID = keyhelper.GenerateRegistrationID()
	signalUser.PreKeys, _ = keyhelper.GeneratePreKeys(0, 100, serializer.PreKeyRecord)
	signalUser.SignedPreKey, _ = keyhelper.GenerateSignedPreKey(signalUser.IdentityKeyPair, 0, serializer.SignedPreKeyRecord)

	signalUser.SessionStore = database.NewInMemorySession(serializer)
	signalUser.PreKeyStore = database.NewInMemoryPreKey()
	signalUser.SignedPreKeyStore = database.NewInMemorySignedPreKey()
	signalUser.IdentityStore = database.NewInMemoryIdentityKey(signalUser.IdentityKeyPair, signalUser.RegistrationID)
	signalUser.SenderKeyStore = database.NewInMemorySenderKey()

	for i := range signalUser.PreKeys {
		signalUser.PreKeyStore.StorePreKey(
			signalUser.PreKeys[i].ID().Value,
			record.NewPreKey(signalUser.PreKeys[i].ID().Value, signalUser.PreKeys[i].KeyPair(), serializer.PreKeyRecord),
		)
	}

	signalUser.SignedPreKeyStore.StoreSignedPreKey(
		signalUser.SignedPreKey.ID(),
		record.NewSignedPreKey(
			signalUser.SignedPreKey.ID(),
			signalUser.SignedPreKey.Timestamp(),
			signalUser.SignedPreKey.KeyPair(),
			signalUser.SignedPreKey.Signature(),
			serializer.SignedPreKeyRecord,
		),
	)

	signalUser.Name = name
	signalUser.DeviceID = deviceID
	signalUser.Address = protocol.NewSignalAddress(name, deviceID)

	signalUser.BuildGroupSession(serializer)

	return signalUser
}

func (c *ClientSession) BuildSession(address *protocol.SignalAddress, serializer *serialize.Serializer) {
	c.SessionBuilder = session.NewBuilder(
		c.SessionStore,
		c.PreKeyStore,
		c.SignedPreKeyStore,
		c.IdentityStore,
		address,
		serializer,
	)
}

func (c *ClientSession) BuildGroupSession(serializer *serialize.Serializer) {
	c.GroupBuilder = groups.NewGroupSessionBuilder(c.SenderKeyStore, serializer)
}
