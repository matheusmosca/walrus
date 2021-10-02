package topics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func TestGetTopic_Success(t *testing.T) {
	type args struct {
		topicName vos.TopicName
	}
	type testCase struct {
		description string
		args        args
		beforeRun   func(storage map[vos.TopicName]entities.Topic)
		want        vos.TopicName
		wantErr     error
	}

	tests := []testCase{
		{
			description: "should return a topic without error",
			args: args{
				topicName: vos.TopicName("topic1"),
			},
			beforeRun: func(storage map[vos.TopicName]entities.Topic) {
				topicName := vos.TopicName("topic1")
				topic, _ := entities.NewTopic(topicName)
				storage[topicName] = topic
			},
			want:    vos.TopicName("topic1"),
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			storage := make(map[vos.TopicName]entities.Topic)
			tt.beforeRun(storage)
			repository := NewMemoryRepository(storage)

			got, err := repository.GetTopic(context.TODO(), tt.args.topicName)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.Equal(t, tt.want, got.GetName())
		})
	}
}

func TestGetTopic_Failure(t *testing.T) {
	type args struct {
		topicName vos.TopicName
	}
	type testCase struct {
		description string
		args        args
		want        vos.TopicName
		wantErr     error
	}

	tests := []testCase{
		{
			description: "should return an error when topic is not found",
			args: args{
				topicName: vos.TopicName("topic1"),
			},
			want:    vos.TopicName(""),
			wantErr: entities.ErrTopicNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			storage := make(map[vos.TopicName]entities.Topic)
			repository := NewMemoryRepository(storage)

			got, err := repository.GetTopic(context.TODO(), tt.args.topicName)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.Equal(t, tt.want, got.GetName())
		})
	}
}
