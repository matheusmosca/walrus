package rpc

import (
	"github.com/matheusmosca/walrus/domain/usecases"
)

type RPC struct {
	useCase usecases.UseCase
}

func New(useCase usecases.UseCase) RPC {
	return RPC{
		useCase: useCase,
	}
}
