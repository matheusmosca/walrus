package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/vos"
)

func (u useCase) Publish(ctx context.Context, message vos.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	topic, err := u.storage.GetTopic(ctx, message.TopicName)
	if err != nil {
		return err
	}

	topic.Dispatch(message)

	return nil
}
