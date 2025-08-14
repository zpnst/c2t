package database

import (
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"sync"

	c2t "github.com/zpnst/c2t/internal/protocol"
	groupRecord "github.com/zpnst/libsignal-protocol-go/groups/state/record"
	"github.com/zpnst/libsignal-protocol-go/keys/identity"
	signal "github.com/zpnst/libsignal-protocol-go/protocol"
	"github.com/zpnst/libsignal-protocol-go/serialize"
	"github.com/zpnst/libsignal-protocol-go/state/record"
	"github.com/zpnst/libsignal-protocol-go/util/bytehelper"
)

// Database

// In Memory Database
type InMemoryDatabase struct {
	ClientsLock sync.Mutex
	Clients     map[string]c2t.Client
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		Clients: make(map[string]c2t.Client),
	}
}

func (i *InMemoryDatabase) GetAllClients() ([]c2t.Client, error) {
	i.ClientsLock.Lock()
	defer i.ClientsLock.Unlock()
	var clinets []c2t.Client
	for _, c := range i.Clients {
		clinets = append(clinets, c)
	}
	return clinets, nil
}

func (i *InMemoryDatabase) GetClient(name string) (c2t.Client, error) {
	i.ClientsLock.Lock()
	defer i.ClientsLock.Unlock()
	return i.Clients[name], nil
}

func (i *InMemoryDatabase) AddClient(name string, client c2t.Client) error {
	i.ClientsLock.Lock()
	defer i.ClientsLock.Unlock()
	if _, exists := i.Clients[name]; !exists {
		client.Password = string(
			bytehelper.ArrayToSlice(sha256.Sum256([]byte(client.Password))))
		i.Clients[name] = client
	} else {
		return errors.New("client with that name already exists")
	}
	return nil
}

func (i *InMemoryDatabase) UpdateClientBundle(name string, bundle c2t.RawBundle) error {
	i.ClientsLock.Lock()
	defer i.ClientsLock.Unlock()
	if _, exists := i.Clients[name]; !exists {
		return errors.New("client with that name already exists")
	} else {
		c := i.Clients[name]
		c.Bundle = bundle
		i.Clients[name] = c
	}
	return nil
}

// Stores

// IdentityKeyStore
type InMemoryIdentityKey struct {
	TrustedKeys         map[*signal.SignalAddress]*identity.Key
	IdentityKeyPair     *identity.KeyPair
	LocalRegistrationID uint32
}

func NewInMemoryIdentityKey(identityKey *identity.KeyPair, localRegistrationID uint32) *InMemoryIdentityKey {
	return &InMemoryIdentityKey{
		TrustedKeys:         make(map[*signal.SignalAddress]*identity.Key),
		IdentityKeyPair:     identityKey,
		LocalRegistrationID: localRegistrationID,
	}
}

func (i *InMemoryIdentityKey) GetIdentityKeyPair() *identity.KeyPair {
	return i.IdentityKeyPair
}

func (i *InMemoryIdentityKey) GetLocalRegistrationId() uint32 {
	return i.LocalRegistrationID
}

func (i *InMemoryIdentityKey) SaveIdentity(address *signal.SignalAddress, identityKey *identity.Key) {
	i.TrustedKeys[address] = identityKey
}

func (i *InMemoryIdentityKey) IsTrustedIdentity(address *signal.SignalAddress, identityKey *identity.Key) bool {
	trusted := i.TrustedKeys[address]
	return (trusted == nil || trusted.Fingerprint() == identityKey.Fingerprint())
}

// PreKeyStore
type InMemoryPreKey struct {
	store map[uint32]*record.PreKey
}

func NewInMemoryPreKey() *InMemoryPreKey {
	return &InMemoryPreKey{
		store: make(map[uint32]*record.PreKey),
	}
}

func (i *InMemoryPreKey) LoadPreKey(preKeyID uint32) *record.PreKey {
	return i.store[preKeyID]
}

func (i *InMemoryPreKey) StorePreKey(preKeyID uint32, preKeyRecord *record.PreKey) {
	i.store[preKeyID] = preKeyRecord
}

func (i *InMemoryPreKey) ContainsPreKey(preKeyID uint32) bool {
	_, ok := i.store[preKeyID]
	return ok
}

func (i *InMemoryPreKey) RemovePreKey(preKeyID uint32) {
	delete(i.store, preKeyID)
}

// SessionStore
type InMemorySession struct {
	Sessions   map[*signal.SignalAddress]*record.Session
	Serializer *serialize.Serializer
}

func NewInMemorySession(serializer *serialize.Serializer) *InMemorySession {
	return &InMemorySession{
		Sessions:   make(map[*signal.SignalAddress]*record.Session),
		Serializer: serializer,
	}
}

func (i *InMemorySession) LoadSession(address *signal.SignalAddress) *record.Session {
	if i.ContainsSession(address) {
		return i.Sessions[address]
	}
	sessionRecord := record.NewSession(i.Serializer.Session, i.Serializer.State)
	i.Sessions[address] = sessionRecord

	return sessionRecord
}

func (i *InMemorySession) GetSubDeviceSessions(name string) []uint32 {
	var deviceIDs []uint32

	for key := range i.Sessions {
		if key.Name() == name && key.DeviceID() != 1 {
			deviceIDs = append(deviceIDs, key.DeviceID())
		}
	}

	return deviceIDs
}

func (i *InMemorySession) StoreSession(remoteAddress *signal.SignalAddress, record *record.Session) {
	i.Sessions[remoteAddress] = record
}

func (i *InMemorySession) ContainsSession(remoteAddress *signal.SignalAddress) bool {
	_, ok := i.Sessions[remoteAddress]
	return ok
}

func (i *InMemorySession) DeleteSession(remoteAddress *signal.SignalAddress) {
	delete(i.Sessions, remoteAddress)
}

func (i *InMemorySession) DeleteAllSessions() {
	i.Sessions = make(map[*signal.SignalAddress]*record.Session)
}

// SignedPreKeyStore
type InMemorySignedPreKey struct {
	Store map[uint32]*record.SignedPreKey
}

func NewInMemorySignedPreKey() *InMemorySignedPreKey {
	return &InMemorySignedPreKey{
		Store: make(map[uint32]*record.SignedPreKey),
	}
}

func (i *InMemorySignedPreKey) LoadSignedPreKey(signedPreKeyID uint32) *record.SignedPreKey {
	return i.Store[signedPreKeyID]
}

func (i *InMemorySignedPreKey) LoadSignedPreKeys() []*record.SignedPreKey {
	var preKeys []*record.SignedPreKey

	for _, record := range i.Store {
		preKeys = append(preKeys, record)
	}

	return preKeys
}

func (i *InMemorySignedPreKey) StoreSignedPreKey(signedPreKeyID uint32, record *record.SignedPreKey) {
	i.Store[signedPreKeyID] = record
}

func (i *InMemorySignedPreKey) ContainsSignedPreKey(signedPreKeyID uint32) bool {
	_, ok := i.Store[signedPreKeyID]
	return ok
}

func (i *InMemorySignedPreKey) RemoveSignedPreKey(signedPreKeyID uint32) {
	delete(i.Store, signedPreKeyID)
}

func NewInMemorySenderKey() *InMemorySenderKey {
	return &InMemorySenderKey{
		Store: make(map[*signal.SenderKeyName]*groupRecord.SenderKey),
	}
}

type InMemorySenderKey struct {
	Store map[*signal.SenderKeyName]*groupRecord.SenderKey
}

func (i *InMemorySenderKey) StoreSenderKey(senderKeyName *signal.SenderKeyName, keyRecord *groupRecord.SenderKey) {
	i.Store[senderKeyName] = keyRecord
}

func (i *InMemorySenderKey) LoadSenderKey(senderKeyName *signal.SenderKeyName) *groupRecord.SenderKey {
	return i.Store[senderKeyName]
}

func init() {
	gob.Register(InMemorySession{})
}
