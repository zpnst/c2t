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

// Get Bundle Message
type MessageGetBundle struct {
	Name string
}

type AnswerGetBundle struct {
	ErrorCode uint8
	Bundle    RawBundle
}

// Set Bundle Message
type MessageSetBundle struct {
	Name   string
	Bundle RawBundle
}

type AnswerSetBundle struct {
	ErrorCode uint8
}

// GOB Register
func init() {
	// Sign Up Message
	gob.Register(MessageSignUp{})
	gob.Register(AnswerSignUp{})

	// Get/Set Bundle
	gob.Register(RawBundle{})
	gob.Register(MessageGetBundle{})
	gob.Register(AnswerGetBundle{})
	gob.Register(MessageSetBundle{})
	gob.Register(AnswerSetBundle{})

	// Signal
	gob.Register(ecc.DjbECPublicKey{})
	gob.Register(ecc.DjbECPrivateKey{})
	gob.Register(serialize.JSONPreKeyRecordSerializer{})
	gob.Register(serialize.JSONSignedPreKeyRecordSerializer{})
}
