package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func (u useCase) Subscribe(ctx context.Context, subscriberID string, topicName vos.TopicName) (<-chan vos.Message, error) {
	if err := topicName.Validate(); err != nil {
		return nil, err
	}

	topic, ok := u.topics[topicName]
	if ok {
		subscriber := entities.NewSubscriber(subscriberID, topic)

		subscriptionCh := subscriber.Subscribe()
		return subscriptionCh, nil
	}

	topic, err := entities.NewTopic(topicName)
	if err != nil {
		return nil, err
	}
	u.topics[topicName] = topic

	topic.Activate()
	subscriber := entities.NewSubscriber(subscriberID, topic)

	subscriptionCh := subscriber.Subscribe()
	return subscriptionCh, nil
}
