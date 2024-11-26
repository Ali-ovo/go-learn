//go:build wireinject
// +build wireinject

package main

import "github.com/google/wire"

func initEvent() Event {
	wire.Build(NewEvent, NewGreeter, NewMessage)
	return Event{}
}
