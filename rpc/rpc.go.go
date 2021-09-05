package rpc

import (
	"github.com/matheusmosca/walrus/domain/entities"
)

type RPC struct {
	dispatcher entities.Topic
}

func New(dispatcher entities.Topic) RPC {
	return RPC{
		dispatcher: dispatcher,
	}
}
