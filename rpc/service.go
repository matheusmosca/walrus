package rpc

import (
	"github.com/matheusmosca/walrus/domain/entities"
)

type RPCServer struct {
	dispatcher entities.Dispatcher
}

func NewServer(dispatcher entities.Dispatcher) RPCServer {
	return RPCServer{
		dispatcher: dispatcher,
	}
}
