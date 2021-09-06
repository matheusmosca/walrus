package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/vos"
)

func (u useCase) Unsubscribe(ctx context.Context, subscriberID string, topicName vos.TopicName) error {
	topic, err := u.storage.GetTopic(ctx, topicName)
	if err != nil {
		return err
	}

	topic.RemoveSubscriber(subscriberID)

	return nil
}
