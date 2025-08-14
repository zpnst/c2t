package main

import (
	"log"

	"github.com/zpnst/c2t/cmd/c2t-instance/database"
	c2t "github.com/zpnst/c2t/internal/protocol"
)

type InstanceOpts struct {
	Database  database.Database
	Transport c2t.ITransport
}

type Instance struct {
	InstanceOpts
}

func NewInstance(opts InstanceOpts) *Instance {
	return &Instance{
		InstanceOpts: opts,
	}
}

func (i Instance) Run() error {
	if err := i.Transport.ListenAndAccept(); err != nil {
		return err
	}
	i.MainLoop()
	return nil
}

func (i Instance) MainLoop() {
	defer func() {
		log.Printf("[%s] :: instance stopped due to error\n",
			i.Transport.Addr())
		i.Transport.Close()
	}()

	for {
		select {
		case m := <-i.Transport.Consume():
			go i.HandleMessage(m)
		}
	}
}
