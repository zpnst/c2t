package main

import (
	c2t "github.com/zpnst/c2t/internal/protocol"
	"github.com/zpnst/c2t/internal/utils"
)

func (i Instance) HandleMessage(m c2t.Message) {
	switch b := m.Body.(type) {
	case c2t.MessageSignUp:
		utils.LogErr(i.handleSignUp(m.From, b))
	case c2t.MessageGetBundle:
		utils.LogErr(i.handleGetBundle(m.From, b))
	case c2t.MessageUpdateBundle:
		utils.LogErr(i.handleSetBundle(m.From, b))
	case c2t.MessagePostEncrypted:
		utils.LogErr(i.handlePostEncrypted(m.From, b))
	case c2t.MessageFetchEncrypted:
		utils.LogErr(i.handleFetchEncrypted(m.From, b))
	}
}

func (i Instance) handleFetchEncrypted(src string, m c2t.MessageFetchEncrypted) error {
	var b c2t.AnswerFetchEncrypted
	if msgs, err := i.Database.GetClientMessages(m.Me); err != nil {
		b = c2t.AnswerFetchEncrypted{
			ErrorCode: c2t.EXIT_NO_SUCH_CLIENT,
			Messages:  []c2t.EncryptedMessage{},
		}
	} else {
		b = c2t.AnswerFetchEncrypted{
			ErrorCode: c2t.EXIT_OK,
			Messages:  msgs,
		}
	}
	return i.Transport.SendAnswer(src, b)
}

func (i Instance) handlePostEncrypted(src string, m c2t.MessagePostEncrypted) error {
	var b c2t.AnswerPostEncrypted
	if err := i.Database.AddClientMessage(m.To, m.Message); err != nil {
		b = c2t.AnswerPostEncrypted{
			ErrorCode: c2t.EXIT_NO_SUCH_CLIENT,
		}
	} else {
		b = c2t.AnswerPostEncrypted{
			ErrorCode: c2t.EXIT_OK,
		}
	}
	return i.Transport.SendAnswer(src, b)
}

func (i Instance) handleGetBundle(src string, m c2t.MessageGetBundle) error {
	var b c2t.AnswerGetBundle
	if c, err := i.Database.GetClient(m.Name); err != nil {
		b = c2t.AnswerGetBundle{
			ErrorCode: c2t.EXIT_NO_SUCH_CLIENT,
			Bundle:    c2t.RawBundle{},
		}
	} else {
		b = c2t.AnswerGetBundle{
			ErrorCode: c2t.EXIT_OK,
			Bundle:    c.Bundle,
		}
	}
	return i.Transport.SendAnswer(src, b)
}

func (i Instance) handleSetBundle(src string, m c2t.MessageUpdateBundle) error {
	var b c2t.AnswerUpdateBundle
	if err := i.Database.UpdateClientBundle(m.Name, m.Bundle); err != nil {
		b = c2t.AnswerUpdateBundle{
			ErrorCode: c2t.EXIT_NO_SUCH_CLIENT,
		}
	} else {
		b = c2t.AnswerUpdateBundle{
			ErrorCode: c2t.EXIT_OK,
		}
	}
	return i.Transport.SendAnswer(src, b)
}

func (i Instance) handleSignUp(src string, m c2t.MessageSignUp) error {
	var b c2t.AnswerSignUp
	var client c2t.Client
	if err := i.Database.AddClient(m.Name, client); err != nil {
		b = c2t.AnswerSignUp{
			ErrorCode: c2t.EXIT_CLIENT_ALREADY_EXISTS,
		}
	} else {
		b = c2t.AnswerSignUp{
			ErrorCode: c2t.EXIT_OK,
		}
	}
	return i.Transport.SendAnswer(src, b)
}
