package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

type topics map[vos.TopicName]entities.Topic

type useCase struct {
	storage Repository
	topics  topics
}

//go:generate moq -fmt goimports -out repository_mock.go . Repository:RepositoryMock
type Repository interface {
	CreateTopic(ctx context.Context, name vos.TopicName, topic entities.Topic) error
	GetTopic(ctx context.Context, topicName vos.TopicName) (entities.Topic, error)
	ListTopics(ctx context.Context) ([]entities.Topic, error)
}

//go:generate moq -fmt goimports -out usecase_mock.go . UseCase:UseCaseMock
type UseCase interface {
	Publish(ctx context.Context, message vos.Message) error
	Unsubscribe(ctx context.Context, subscriberID vos.SubscriberID, topicName vos.TopicName) error
	Subscribe(ctx context.Context, topicName vos.TopicName) (chan vos.Message, vos.SubscriberID, error)
	ListTopics(ctx context.Context) ([]entities.Topic, error)
}

func New(storage Repository) UseCase {
	return useCase{
		storage: storage,
		topics:  make(topics),
	}
}
