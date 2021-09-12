package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func (u useCase) createTopic(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
	topic, err := entities.NewTopic(topicName)
	if err != nil {
		return entities.Topic{}, err
	}
	err = u.storage.CreateTopic(ctx, topicName, topic)
	if err != nil {
		return entities.Topic{}, err
	}

	topic.Activate()

	return topic, nil
}
