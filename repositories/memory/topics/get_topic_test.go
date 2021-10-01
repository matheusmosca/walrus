package topics

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
	"github.com/stretchr/testify/require"
)

type args struct {
	topicName vos.TopicName
}
type testCase struct {
	description string
	args        args
	beforeRun   func(storage map[vos.TopicName]entities.Topic)
	wantErr     error
}

func TestGetTopic_Success(t *testing.T) {
	// positiveTestcase creates a topic and stores it in the repository.storage as a part of its beforeRun routine
	// hence these test cases should run with positive results for GetTopic method
	positiveTestcase := []testCase{
		{
			description: "The first positive topic",
			args: args{
				topicName: vos.TopicName("pos_topic_1"),
			},
			beforeRun: func(storage map[vos.TopicName]entities.Topic) {
				topicName := vos.TopicName("pos_topic_1")
				topic, _ := entities.NewTopic(topicName)
				storage[topicName] = topic
			},
			wantErr: nil,
		},
	}

	for _, tt := range positiveTestcase {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			storage := make(map[vos.TopicName]entities.Topic)
			tt.beforeRun(storage)

			newTopic, err := entities.NewTopic(tt.args.topicName)
			require.NoError(t, err)
			require.NotEmpty(t, newTopic)

			repository := NewMemoryRepository(storage)
			getTopic, err := repository.GetTopic(context.TODO(), tt.args.topicName)

			if tt.wantErr == nil {
				require.NoError(t, err)
			}
			assert.NotEmpty(t, getTopic)
		})
	}
}

func TestGetTopic_Negative(t *testing.T) {
	// negativeTestcase neither creates a topic nor stores it in the repository.storage as a part of its beforeRun routine
	// hence these test cases should run with negative results for GetTopic method
	negativeTestcase := []testCase{
		{
			description: "The first negative topic",
			args: args{
				topicName: vos.TopicName("neg_topic_1"),
			},
			beforeRun: func(storage map[vos.TopicName]entities.Topic) {
			},
			wantErr: entities.ErrTopicNotFound,
		},
	}

	for _, tt := range negativeTestcase {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			storage := make(map[vos.TopicName]entities.Topic)
			tt.beforeRun(storage)

			newTopic, err := entities.NewTopic(tt.args.topicName)
			require.NoError(t, err)
			assert.NotEmpty(t, newTopic)

			repository := NewMemoryRepository(storage)
			getTopic, err := repository.GetTopic(context.TODO(), tt.args.topicName)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			}
			assert.Empty(t, getTopic)
		})
	}
}
