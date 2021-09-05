package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

type topics map[vos.TopicName]entities.Topic

type useCase struct {
	// TODO add a repository layer and add TopicRepository via
	// dependecy injection here
	topics topics
}

type UseCase interface {
	Publish(ctx context.Context, message vos.Message) error
	Subscribe(ctx context.Context, subscriberID string, topicName vos.TopicName) (<-chan vos.Message, error)
}

func New() UseCase {
	return useCase{
		topics: make(topics),
	}
}
