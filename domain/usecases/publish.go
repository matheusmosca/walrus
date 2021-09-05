package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func (u useCase) Publish(ctx context.Context, message vos.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	topic, ok := u.topics[message.TopicName]
	if !ok {
		return entities.ErrTopicNotFound
	}

	topic.Dispatch(message)

	return nil
}
