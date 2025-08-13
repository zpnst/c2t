package main

import (
	"github.com/zpnst/c2t/internal/protocol"
)

type ClientOpts struct {
	Transport protocol.CTransport
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
	if err := c.Transport.DialInstance(); err != nil {
		return err
	}
	return nil
}
