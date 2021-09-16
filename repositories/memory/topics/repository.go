package topics

import (
	"context"
	"sync"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

type repository struct {
	storage map[vos.TopicName]entities.Topic
	mu      sync.Mutex
}

type MemoryRepository interface {
	GetTopic(ctx context.Context, topicName vos.TopicName) (entities.Topic, error)
	CreateTopic(ctx context.Context, name vos.TopicName, topic entities.Topic) error
}

func NewMemoryRepository(storage map[vos.TopicName]entities.Topic) MemoryRepository {
	return &repository{
		storage: storage,
	}
}
