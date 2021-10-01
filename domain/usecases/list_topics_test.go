package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

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
						return nil, entities.ErrTopicsNotFound
					},
				},
			},
			want:    nil,
			wantErr: entities.ErrTopicsNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			useCase := New(tt.fields.storage)
			topics, err := useCase.ListTopics(tt.args.ctx)

			var actualTopics []vos.TopicName
			for _, topic := range topics {
				actualTopics = append(actualTopics, topic.GetName())
			}

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, actualTopics)
		})
	}
}
