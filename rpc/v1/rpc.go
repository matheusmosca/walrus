package rpc

import (
	"github.com/matheusmosca/walrus/domain/usecases"
	"github.com/sirupsen/logrus"
)

type RPC struct {
	useCase usecases.UseCase
	log     *logrus.Entry
}

func New(useCase usecases.UseCase, log *logrus.Entry) RPC {
	return RPC{
		useCase: useCase,
		log:     log,
	}
}
