package topics

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

type repository struct {
	storage map[vos.TopicName]entities.Topic
}

type MemoryRepository interface {
	GetTopic(ctx context.Context, topicName vos.TopicName) (entities.Topic, error)
	CreateTopic(ctx context.Context, name vos.TopicName, topic entities.Topic) error
}

func NewMemoryRepository(storage map[vos.TopicName]entities.Topic) MemoryRepository {
	return repository{
		storage: storage,
	}
}
