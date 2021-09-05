package rpc

import (
	"github.com/matheusmosca/walrus/domain/entities"
)

type RPC struct {
	dispatcher entities.Dispatcher
}

func New(dispatcher entities.Dispatcher) RPC {
	return RPC{
		dispatcher: dispatcher,
	}
}
