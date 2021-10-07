package usecases

import (
	"context"
	"testing"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscribe(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		topicName vos.TopicName
		message   vos.Message
	}

	tests := []struct {
		name    string
		args    args
		fields  func(t *testing.T) (Repository, entities.Topic)
		wantErr error
	}{
		{
			name: "subscribe should succeed",
			args: args{
				ctx:       context.Background(),
				topicName: "walrus",
				message: vos.Message{
					TopicName:   vos.TopicName("walrus"),
					PublishedBy: "test_publisher",
					Body:        []byte("hello world"),
				},
			},
			fields: func(t *testing.T) (Repository, entities.Topic) {
				topic, err := entities.NewTopic("walrus")
				require.NoError(t, err)

				repoMock := &RepositoryMock{
					GetTopicFunc: func(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
						return topic, nil
					},
				}

				return repoMock, topic
			},
			wantErr: nil,
		},
		{
			name: "short topic name should return error",
			args: args{
				ctx:       context.Background(),
				topicName: "xy",
				message:   vos.Message{},
			},
			fields:  func(t *testing.T) (Repository, entities.Topic) { return &RepositoryMock{}, entities.Topic{} },
			wantErr: vos.ErrTopicNameTooShort,
		},
		{
			name: "empty topic name message should return error",
			args: args{
				ctx:       context.Background(),
				topicName: "",
				message:   vos.Message{},
			},
			fields:  func(t *testing.T) (Repository, entities.Topic) { return &RepositoryMock{}, entities.Topic{} },
			wantErr: vos.ErrEmptyTopicName,
		},
		{
			name: "nonexistent topic subscribe should succeed",
			args: args{
				ctx:       context.Background(),
				topicName: "newTopic",
				message: vos.Message{
					TopicName:   vos.TopicName("newTopic"),
					PublishedBy: "test_publisher",
					Body:        []byte("hello world"),
				},
			},
			fields: func(t *testing.T) (Repository, entities.Topic) {
				topic, err := entities.NewTopic("newTopic")
				require.NoError(t, err)

				repoMock := &RepositoryMock{
					GetTopicFunc: func(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
						return topic, entities.ErrTopicNotFound
					},
					CreateTopicFunc: func(ctx context.Context, topicName vos.TopicName, topic entities.Topic) error {
						return nil
					},
				}

				return repoMock, topic
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo, topic := tt.fields(t)
			topic.Activate()

			useCase := New(repo)
			subCh, id, err := useCase.Subscribe(tt.args.ctx, tt.args.topicName)
			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr != nil {
				assert.Empty(t, subCh)
				assert.Empty(t, id)
				return
			}
			defer close(subCh)

			assert.NotEmpty(t, id)
		})
	}
}
