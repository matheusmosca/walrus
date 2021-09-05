package topics

import (
	"context"
	"testing"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTopic(t *testing.T) {
	type args struct {
		topicName vos.TopicName
	}
	type testCase struct {
		description string
		args        args
		beforeRun   func(storage map[vos.TopicName]entities.Topic)
		wantErr     error
	}

	tests := []testCase{
		{
			description: "should not create a new topic, topic already exists",
			args: args{
				topicName: vos.TopicName("topic_name"),
			},
			beforeRun: func(storage map[vos.TopicName]entities.Topic) {
				topicName := vos.TopicName("topic_name")
				topic, _ := entities.NewTopic(topicName)
				storage[topicName] = topic
			},
			wantErr: entities.ErrTopicAlreadyExists,
		},
		{
			description: "should create a topic successfully",
			args: args{
				topicName: vos.TopicName("topic_name"),
			},
			beforeRun: func(storage map[vos.TopicName]entities.Topic) {},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			storage := make(map[vos.TopicName]entities.Topic)
			tt.beforeRun(storage)

			newTopic, err := entities.NewTopic(tt.args.topicName)
			require.NoError(t, err)

			repository := NewMemoryRepository(storage)
			err = repository.CreateTopic(context.TODO(), tt.args.topicName, newTopic)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				topic, ok := storage[tt.args.topicName]
				assert.True(t, ok)
				assert.NotEmpty(t, topic)
			}
		})
	}
}
