package database

import (
	c2t "github.com/zpnst/c2t/internal/protocol"
)

type Database interface {
	GetAllClients() ([]c2t.Client, error)
	GetClient(name string) (c2t.Client, error)
	AddClient(name string, client c2t.Client) error
	UpdateClientBundle(name string, bundle c2t.RawBundle) error
	AddClientMessage(name string, message c2t.EncryptedMessage) error
	GetClientMessages(name string) ([]c2t.EncryptedMessage, error)
}
