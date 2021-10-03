package usecases

import (
	"context"
	"testing"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubsribe(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		topicName vos.TopicName
	}

	type fields struct {
		storage   Repository
		topicName vos.TopicName
	}

	tests := []struct {
		name    string
		args    args
		fields  fields
		wantCh  chan vos.Message
		wantId  vos.SubscriberID
		wantErr error
	}{
		{
			name: "subscribe should succeed",
			args: args{
				ctx:       context.Background(),
				topicName: "walrus",
			},
			fields: fields{
				storage:   &RepositoryMock{},
				topicName: "walrus",
			},
			wantCh:  make(chan vos.Message),
			wantId:  vos.SubscriberID("id"),
			wantErr: nil,
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
			wantCh:  nil,
			wantId:  "",
			wantErr: vos.ErrTopicNameTooShort,
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
			wantCh:  nil,
			wantId:  "",
			wantErr: vos.ErrEmptyTopicName,
		},
		{
			name: "nonexistent topic should return error",
			args: args{
				ctx:       context.Background(),
				topicName: "testtopic",
			},
			fields: fields{
				storage:   &RepositoryMock{},
				topicName: "walrus",
			},
			wantCh:  nil,
			wantId:  "",
			wantErr: entities.ErrTopicNotFound,
		},
		{
			name: "topic already exists should return error",
			args: args{
				ctx:       context.Background(),
				topicName: "walrus",
			},
			fields: fields{
				storage:   &RepositoryMock{},
				topicName: "walrus",
			},
			wantCh:  nil,
			wantId:  "",
			wantErr: entities.ErrTopicNotFound,
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

			tt.fields.storage = &RepositoryMock{
				GetTopicFunc: func(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
					if tt.fields.topicName != tt.args.topicName {
						return entities.Topic{}, entities.ErrTopicNotFound
					}

					return topic, nil
				},
				CreateTopicFunc: func(ctx context.Context, topicName vos.TopicName, topic entities.Topic) error {
					if tt.fields.topicName == tt.args.topicName {
						return entities.ErrTopicAlreadyExists
					}

					return nil
				},
			}

			useCase := New(tt.fields.storage)
			subCh, id, err := useCase.Subscribe(tt.args.ctx, tt.args.topicName)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.IsType(t, tt.wantCh, subCh)
			assert.IsType(t, tt.wantId, id)
		})
	}
}
