package database

import "github.com/zpnst/libsignal-protocol-go/protocol"

type Database interface {
	GetAllClients() []*protocol.SignalAddress
	GetClient(name string) *protocol.SignalAddress

	AddClient(name string, deviceID uint32) error
}
