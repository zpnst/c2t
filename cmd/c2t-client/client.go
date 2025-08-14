package main

import (
	c2t "github.com/zpnst/c2t/internal/protocol"
)

type ClientOpts struct {
	Transport c2t.CTransport
}

type Client struct {
	ClientOpts
}

func NewClient(opts ClientOpts) *Client {
	return &Client{
		ClientOpts: opts,
	}
}

func (c *Client) Init() error {
	return c.Transport.DialInstance()
}

func (c Client) WaitForInstance() (any, error) {
	var m c2t.Message
	if err := c.Transport.DecodeAnswer(&m); err != nil {
		return nil, err
	}
	return c.HandleAnswer(m)
}

func (c Client) Send(m c2t.Message) (any, error) {
	if err := c.Transport.EncodeMessage(&m); err != nil {
		return nil, err
	}
	return c.WaitForInstance()
}

func (c Client) OnlyErrorSend(m c2t.Message) error {
	if _, err := c.Send(m); err != nil {
		return err
	}
	return nil
}

func (c Client) SendSignUp(name, password string) error {
	m := c2t.Message{
		Body: c2t.MessageSignUp{
			Name:     name,
			Password: password,
		},
	}
	return c.OnlyErrorSend(m)
}

func (c Client) SendSignIn(name, password string) error {
	return nil
}

func (c Client) SendGetBundle(name string) (c2t.RawBundle, error) {
	var m c2t.Message = c2t.Message{
		Body: c2t.MessageGetBundle{
			Name: name,
		},
	}
	answ, err := c.Send(m)
	if err != nil {
		return c2t.RawBundle{}, err
	}
	return answ.(c2t.AnswerGetBundle).Bundle, nil
}

func (c Client) SendSetBundle(name string, bundle c2t.RawBundle) error {
	var m c2t.Message = c2t.Message{
		Body: c2t.MessageUpdateBundle{
			Name:   name,
			Bundle: bundle,
		},
	}
	return c.OnlyErrorSend(m)
}

func (c Client) SendPostEncryptedMessage(to string, message c2t.EncryptedMessage) error {
	var m c2t.Message = c2t.Message{
		Body: c2t.MessagePostEncrypted{
			To:      to,
			Message: message,
		},
	}
	return c.OnlyErrorSend(m)
}

func (c Client) SendFetchEncryptedMessage(me string) ([]c2t.EncryptedMessage, error) {
	var m c2t.Message = c2t.Message{
		Body: c2t.MessageFetchEncrypted{
			Me: me,
		},
	}
	answ, err := c.Send(m)
	if err != nil {
		return nil, err
	}
	return answ.(c2t.AnswerFetchEncrypted).Messages, nil
}
