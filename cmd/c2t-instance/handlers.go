package main

import (
	"encoding/gob"
	"log"

	"github.com/zpnst/c2t/internal/protocol"
)

func (i Instance) HandleMessage(m *protocol.Message) {
	switch b := m.Body.(type) {
	case protocol.MessageSignUp:
		i.handleSignUp(m.From, b)
	case protocol.SMessage:
		i.handleSMessage(m.From, b)
	}
}

func (i Instance) handleSignUp(src string, m protocol.MessageSignUp) {
	p := i.Transport.Peer(src)
	ans := protocol.Message{
		From: i.Transport.Addr(),
	}
	if err := i.Database.AddClient(m.Name, m.DeviceID); err != nil {
		ans.Body = protocol.AnswerSignUp{
			ErrorCode: protocol.USER_ALREADY_EXISTS,
		}
		if err := gob.NewEncoder(p).Encode(ans); err != nil {
			log.Println(err)
		}
	} else {
		ans.Body = protocol.AnswerSignUp{
			ErrorCode: protocol.EXIT_OK,
		}
		if err := gob.NewEncoder(p).Encode(ans); err != nil {
			log.Println(err)
		}
	}
}

func (i Instance) handleSMessage(src string, m protocol.SMessage) {
	log.Println("string:", m.S, src)
}
