package topics

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func (r repository) GetTopic(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
	topic, ok := r.storage[topicName]
	if !ok {
		return nil, entities.ErrTopicNotFound
	}

	return topic, nil
}
