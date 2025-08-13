package protocol

import "encoding/gob"

type Message struct {
	From string
	Body any
}

// From Client to Instance
type MessageSignUp struct {
	Name     string
	DeviceID uint32
}

type SMessage struct {
	S string
}

// From Instance to Client
type AnswerSignUp struct {
	ErrorCode uint8
}

// GOB Register
func init() {
	gob.Register(SMessage{})
	gob.Register(MessageSignUp{})
	gob.Register(AnswerSignUp{})
}
