package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func TestListTopics(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
	}

	type fields struct {
		storage Repository
	}

	tests := []struct {
		name      string
		args      args
		fields    fields
		beforeRun func(topic entities.Topic) chan vos.Message
		want      []vos.TopicName
		wantErr   error
	}{
		{
			name: "list topics should success",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				storage: &RepositoryMock{
					ListTopicsFunc: func(ctx context.Context) ([]entities.Topic, error) {
						walrus, _ := entities.NewTopic("walrus")
						walrus2, _ := entities.NewTopic("walrus2")
						return []entities.Topic{walrus, walrus2}, nil
					},
				},
			},
			want:    []vos.TopicName{"walrus", "walrus2"},
			wantErr: nil,
		},
		{
			name: "no topics in storage topic should returns error",
			args: args{
				ctx: context.Background(),
			},
			fields: fields{
				storage: &RepositoryMock{
					ListTopicsFunc: func(ctx context.Context) ([]entities.Topic, error) {
						return nil, entities.ErrNoTopicsFound
					},
				},
			},
			want:    nil,
			wantErr: entities.ErrNoTopicsFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			useCase := New(tt.fields.storage)
			topics, err := useCase.ListTopics(tt.args.ctx)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Empty(t, topics)
				return
			}

			require.NoError(t, err)
			topicNames := make([]vos.TopicName, 0, len(tt.want))
			for _, topic := range topics {
				topicNames = append(topicNames, topic.GetName())
			}

			assert.Equal(t, tt.want, topicNames)
		})
	}
}
