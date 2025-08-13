package database

import (
	"errors"

	"github.com/zpnst/libsignal-protocol-go/protocol"
)

type InMemoryDatabase struct {
	clients map[string]*protocol.SignalAddress
}

func NewInMemoryDatabase() *InMemoryDatabase {
	return &InMemoryDatabase{
		clients: make(map[string]*protocol.SignalAddress),
	}
}

func (i InMemoryDatabase) GetAllClients() []*protocol.SignalAddress {
	var clinets []*protocol.SignalAddress
	for _, c := range i.clients {
		clinets = append(clinets, c)
	}
	return clinets
}

func (i InMemoryDatabase) GetClient(name string) *protocol.SignalAddress {
	return i.clients[name]
}

func (i *InMemoryDatabase) AddClient(name string, deviceID uint32) error {
	addr := protocol.NewSignalAddress(name, deviceID)
	if _, exists := i.clients[name]; !exists {
		i.clients[name] = addr
	} else {
		return errors.New("a client with that name already exists")
	}
	return nil
}
