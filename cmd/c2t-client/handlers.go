package main

import (
	c2t "github.com/zpnst/c2t/internal/protocol"
)

func (c Client) HandleAnswer(m c2t.Message) (any, error) {
	switch b := m.Body.(type) {
	case c2t.AnswerSignUp:
		if b.ErrorCode == c2t.EXIT_OK {
			return nil, nil
		} else if b.ErrorCode == c2t.EXIT_CLIENT_ALREADY_EXISTS {
			return nil, ErrClientAlreadyExists
		}
	case c2t.AnswerGetBundle:
		if b.ErrorCode == c2t.EXIT_OK {
			return b, nil
		} else if b.ErrorCode == c2t.EXIT_NO_SUCH_CLIENT {
			return nil, ErrNoSuchClient
		}
	case c2t.AnswerUpdateBundle:
		if b.ErrorCode == c2t.EXIT_OK {
			return nil, nil
		} else if b.ErrorCode == c2t.EXIT_NO_SUCH_CLIENT {
			return nil, ErrNoSuchClient
		}
	case c2t.AnswerFetchEncrypted:
		if b.ErrorCode == c2t.EXIT_OK {
			return b, nil
		} else if b.ErrorCode == c2t.EXIT_NO_SUCH_CLIENT {
			return nil, ErrNoSuchClient
		}
	}
	return nil, nil
}
