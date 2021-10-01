package usecases

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
)

func (u useCase) ListTopics(ctx context.Context) ([]entities.Topic, error) {
	topics, err := u.storage.ListTopics(ctx)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
