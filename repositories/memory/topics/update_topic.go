package topics

import (
	"context"

	"github.com/matheusmosca/walrus/domain/entities"
)

func (r *repository) UpdateTopic(ctx context.Context, topicName entities.Topic) (entities.Topic, error) {
	return entities.Topic{}, nil
}
