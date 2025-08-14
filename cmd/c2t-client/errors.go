package main

import "errors"

var ErrClientAlreadyExists = errors.New("client with that name already exists")

var ErrNoSuchClient = errors.New("no such client")

var ErrStub = errors.New("some error")
