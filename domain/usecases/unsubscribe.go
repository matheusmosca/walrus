package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/vos"
)

func (u useCase) Unsubscribe(ctx context.Context, subscriberID vos.SubscriberID, topicName vos.TopicName) error {
	topic, err := u.storage.GetTopic(ctx, topicName)
	if err != nil {
		return err
	}

	err = topic.RemoveSubscriber(subscriberID)
	if err != nil {
		return err
	}

	//u.storage.UpdateTopic(ctx, topic)

	return nil
}
