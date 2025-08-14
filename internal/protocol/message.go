package protocol

import (
	"encoding/gob"

	"github.com/zpnst/libsignal-protocol-go/ecc"
	"github.com/zpnst/libsignal-protocol-go/serialize"
)

type Message struct {
	From string
	Body any
}

// Sign Up Message
type MessageSignUp struct {
	Name     string
	Password string
}

type AnswerSignUp struct {
	ErrorCode uint8
}

// Update/Get Bundle Message
type MessageUpdateBundle struct {
	Name   string
	Bundle RawBundle
}

type AnswerUpdateBundle struct {
	ErrorCode uint8
}

type MessageGetBundle struct {
	Name string
}

type AnswerGetBundle struct {
	ErrorCode uint8
	Bundle    RawBundle
}

// Post/Get Encrypted Message
type MessagePostEncrypted struct {
	To      string
	Message EncryptedMessage
}

type AnswerPostEncrypted struct {
	ErrorCode uint8
}

type MessageFetchEncrypted struct {
	Me string
}

type AnswerFetchEncrypted struct {
	ErrorCode uint8
	Messages  []EncryptedMessage
}

// GOB Register
func init() {
	// Sign Up Message
	gob.Register(MessageSignUp{})
	gob.Register(AnswerSignUp{})

	// Update/Get Bundle
	gob.Register(RawBundle{})
	gob.Register(MessageGetBundle{})
	gob.Register(AnswerGetBundle{})
	gob.Register(MessageUpdateBundle{})
	gob.Register(AnswerUpdateBundle{})

	// Post/Fetch Encrypted Message
	gob.Register(MessagePostEncrypted{})
	gob.Register(AnswerPostEncrypted{})
	gob.Register(MessageFetchEncrypted{})
	gob.Register(AnswerFetchEncrypted{})

	// Signal
	gob.Register(ecc.DjbECPublicKey{})
	gob.Register(ecc.DjbECPrivateKey{})
	gob.Register(serialize.JSONPreKeyRecordSerializer{})
	gob.Register(serialize.JSONSignedPreKeyRecordSerializer{})
}
