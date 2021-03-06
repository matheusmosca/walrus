package topics

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
)

func (r *repository) ListTopics(_ context.Context) ([]entities.Topic, error) {
	topics := make([]entities.Topic, 0, len(r.storage))

	for _, topic := range r.storage {
		topics = append(topics, topic)
	}

	if len(topics) < 1 {
		return nil, entities.ErrTopicsNotFound
	}

	return topics, nil
}
