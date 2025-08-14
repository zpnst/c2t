package main

import (
	c2t "github.com/zpnst/c2t/internal/protocol"
)

func (c Client) HandleAnswer(m c2t.Message) (any, error) {
	switch b := m.Body.(type) {
	case c2t.AnswerSignUp:
		if b.ErrorCode == c2t.ERR_OK {
			return nil, nil
		} else {
			return nil, ErrClientAlreadyExists
		}
	case c2t.AnswerGetBundle:
		if b.ErrorCode == c2t.ERR_OK {
			return b, nil
		} else {
			return nil, ErrNoSuchClient
		}
	case c2t.AnswerSetBundle:
		if b.ErrorCode == c2t.ERR_OK {
			return nil, nil
		} else {
			return nil, ErrNoSuchClient
		}
	}
	return nil, nil
}
