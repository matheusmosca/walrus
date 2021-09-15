package topics

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func (r *repository) CreateTopic(ctx context.Context, name vos.TopicName, topic entities.Topic) error {
	if _, ok := r.storage[name]; ok {
		return entities.ErrTopicAlreadyExists
	}

	r.mu.Lock()
	r.storage[name] = topic
	r.mu.Unlock()

	return nil
}
