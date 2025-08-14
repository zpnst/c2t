package protocol

import (
	"github.com/zpnst/libsignal-protocol-go/ecc"
	"github.com/zpnst/libsignal-protocol-go/keys/identity"
	"github.com/zpnst/libsignal-protocol-go/keys/prekey"
	"github.com/zpnst/libsignal-protocol-go/util/optional"
)

type EncryptedMessage struct {
	FromName     string
	FromDeviceID uint32
	Message      []byte
}

type Client struct {
	Name       string
	Password   string
	Bundle     RawBundle
	MessagesIn []EncryptedMessage
}

type RawBundle struct {
	RegistrationID        uint32
	DeviceID              uint32
	PreKeyID              uint32
	SignedPreKeyID        uint32
	PreKeyPublic          [32]byte
	SignedPreKeyPublic    [32]byte
	SignedPreKeySignature [64]byte
	IdentityKeyPublic     [32]byte
}

func NewRawBundle(registrationID, deviceID uint32, preKeyID uint32, signedPreKeyID uint32,
	preKeyPublic, signedPreKeyPublic, identityKey [32]byte, signedPreKeySig [64]byte) *RawBundle {

	bundle := RawBundle{
		RegistrationID:        registrationID,
		DeviceID:              deviceID,
		PreKeyID:              preKeyID,
		PreKeyPublic:          preKeyPublic,
		SignedPreKeyID:        signedPreKeyID,
		SignedPreKeyPublic:    signedPreKeyPublic,
		SignedPreKeySignature: signedPreKeySig,
		IdentityKeyPublic:     identityKey,
	}

	return &bundle
}

func (rb *RawBundle) ToSignalBundle() *prekey.Bundle {
	bundle := prekey.NewBundle(
		rb.RegistrationID,
		rb.DeviceID,
		&optional.Uint32{
			Value:   rb.PreKeyID,
			IsEmpty: false,
		},
		rb.SignedPreKeyID,
		ecc.NewDjbECPublicKey(rb.PreKeyPublic),
		ecc.NewDjbECPublicKey(rb.SignedPreKeyPublic),
		rb.SignedPreKeySignature,
		identity.NewKey(ecc.NewDjbECPublicKey(rb.IdentityKeyPublic)),
	)
	return bundle
}
