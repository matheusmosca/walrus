package usecases

import (
	"context"
	"errors"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func (u useCase) Subscribe(ctx context.Context, topicName vos.TopicName) (chan vos.Message, string, error) {
	if err := topicName.Validate(); err != nil {
		return nil, "", err
	}

	topic, err := u.storage.GetTopic(ctx, topicName)
	if err != nil {
		if !errors.Is(err, entities.ErrTopicNotFound) {
			return nil, "", err
		}

		topic, err = u.createTopic(ctx, topicName)
		if err != nil {
			return nil, "", err
		}

		subscriber := entities.NewSubscriber(topic)

		subscriptionCh, id := subscriber.Subscribe()
		return subscriptionCh, id, nil
	}

	subscriber := entities.NewSubscriber(topic)

	subscriptionCh, id := subscriber.Subscribe()
	return subscriptionCh, id, nil
}
