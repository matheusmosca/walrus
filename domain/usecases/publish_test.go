package usecases

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func Test_useCase_publish(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		message vos.Message
	}

	type fields struct {
		wg        *sync.WaitGroup
		storage   Repository
		topicName vos.TopicName
	}

	tests := []struct {
		name      string
		args      args
		fields    fields
		beforeRun func(topic entities.Topic) chan vos.Message
		want      []byte
		wantErr   error
	}{
		{
			name: "publish should success",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "walrus",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				wg:        &sync.WaitGroup{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				subscriber.Subscribe()
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    []byte("hello world"),
			wantErr: nil,
		},
		{
			name: "short topic name message should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "xd",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				wg:        &sync.WaitGroup{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				subscriber.Subscribe()
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    nil,
			wantErr: vos.ErrTopicNameTooShort,
		},
		{
			name: "empty topic name message should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				wg:        &sync.WaitGroup{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				subscriber.Subscribe()
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    nil,
			wantErr: vos.ErrEmptyTopicName,
		},
		{
			name: "empty published by message should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "nicetopic",
					PublishedBy: "",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				wg:        &sync.WaitGroup{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				subscriber.Subscribe()
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    nil,
			wantErr: vos.ErrEmptyPublishedBy,
		},
		{
			name: "short published by message should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "nicetopic",
					PublishedBy: "xd",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				wg:        &sync.WaitGroup{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				subscriber.Subscribe()
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    nil,
			wantErr: vos.ErrPublishedByTooShort,
		},
		{
			name: "nonexistent topic should return error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "nicetopic",
					PublishedBy: "ddxd",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				wg:        &sync.WaitGroup{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				subscriber.Subscribe()
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    nil,
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

			// activate de topic
			topic.Activate()

			// subscribe in the topic
			subsCh := tt.beforeRun(topic)
			defer close(subsCh)

			tt.fields.storage = &RepositoryMock{
				GetTopicFunc: func(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
					if tt.fields.topicName != tt.args.message.TopicName {
						return entities.Topic{}, entities.ErrTopicNotFound
					}

					return topic, nil
				},
			}

			useCase := New(tt.fields.storage)
			err = useCase.Publish(tt.args.ctx, tt.args.message)
			if err != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			actualMsg := <-subsCh
			require.Equal(t, tt.want, actualMsg.Body)
		})
	}
}
