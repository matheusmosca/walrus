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

	type fields struct {
		storage   Repository
		topicName vos.TopicName
	}

	tests := []struct {
		name     string
		args     args
		fields   fields
		wantBody []byte
		wantId   vos.SubscriberID
		wantErr  error
	}{
		{
			name: "subscribe should succeed",
			args: args{
				ctx:       context.Background(),
				topicName: "walrus",
				message: vos.Message{
					TopicName:   "walrus",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage: &RepositoryMock{
					GetTopicFunc: func(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
						return entities.Topic{}, nil
					}},
				topicName: "walrus",
			},
			wantBody: []byte("hello world"),
			wantId:   vos.SubscriberID("id"),
			wantErr:  nil,
		},
		{
			name: "short topic name should return error",
			args: args{
				ctx:       context.Background(),
				topicName: "xd",
			},
			fields: fields{
				storage:   &RepositoryMock{},
				topicName: "walrus",
			},
			wantBody: nil,
			wantId:   "",
			wantErr:  vos.ErrTopicNameTooShort,
		},
		{
			name: "empty topic name message should return error",
			args: args{
				ctx:       context.Background(),
				topicName: "",
			},
			fields: fields{
				storage:   &RepositoryMock{},
				topicName: "walrus",
			},
			wantBody: nil,
			wantId:   "",
			wantErr:  vos.ErrEmptyTopicName,
		},
		{
			name: "nonexistent topic should succeed",
			args: args{
				ctx:       context.Background(),
				topicName: "newTopic",
				message: vos.Message{
					TopicName:   "newTopic",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage: &RepositoryMock{
					GetTopicFunc: func(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
						return entities.Topic{}, entities.ErrTopicNotFound
					},
					CreateTopicFunc: func(ctx context.Context, topicName vos.TopicName, topic entities.Topic) error {
						return nil
					}},
				topicName: "newTopic",
			},
			wantBody: []byte("hello world"),
			wantId:   vos.SubscriberID("id"),
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// init the topic
			topic, err := entities.NewTopic(tt.fields.topicName)
			require.NoError(t, err)

			// activate the topic
			topic.Activate()

			useCase := New(tt.fields.storage)
			subCh, id, err := useCase.Subscribe(tt.args.ctx, tt.args.topicName)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr != nil {
				assert.Empty(t, subCh)
				assert.Empty(t, id)
				return
			}

			err = topic.Dispatch(tt.args.message)
			require.NoError(t, err)

			gotMessage := <-subCh
			assert.Equal(t, tt.wantBody, gotMessage.Body)
			assert.NoError(t, err)
			assert.NotEmpty(t, id)
		})
	}
}
