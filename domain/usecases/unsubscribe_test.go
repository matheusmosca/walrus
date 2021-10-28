package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func TestUnsubscribe(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx     context.Context
		message vos.Message
	}

	type fields struct {
		storage   Repository
		topicName vos.TopicName
	}

	tests := []struct {
		name           string
		args           args
		fields         fields
		beforeRun      func(topic entities.Topic) (chan vos.Message, vos.SubscriberID)
		wantErr        error
		wantSubscriber bool
	}{
		{
			name: "unsubscribe should succeed",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName: "unsubscribing",
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				topicName: "unsubscribing",
			},
			beforeRun: func(topic entities.Topic) (chan vos.Message, vos.SubscriberID) {
				subscriber := entities.NewSubscriber(topic)
				ch, ID := subscriber.Subscribe()

				// Assert that the subscriber has been created successfully
				sub, err := topic.GetSubscriber(ID)
				require.NoError(t, err)
				assert.NotNil(t, sub)

				return ch, ID
			},
			wantErr:        nil,
			wantSubscriber: true,
		},
		{
			name: "try to unsubscribe a subscriber that doesnt exist should fail",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName: "unsubscribing",
				},
			},
			fields: fields{
				storage:   &RepositoryMock{},
				topicName: "unsubscribing",
			},
			beforeRun: func(topic entities.Topic) (chan vos.Message, vos.SubscriberID) {
				msgChan := make(chan vos.Message)

				return msgChan, vos.SubscriberID("")
			},
			wantErr:        entities.ErrSubscriberNotFound,
			wantSubscriber: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			topic, err := entities.NewTopic(tt.fields.topicName)
			require.NoError(t, err)

			topic.Activate()

			subsCh, subsID := tt.beforeRun(topic)
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
			err = useCase.Unsubscribe(context.Background(), subsID, tt.fields.topicName)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			// Assert the unsubscription
			_, err = topic.GetSubscriber(subsID)
			assert.ErrorIs(t, entities.ErrSubscriberNotFound, err)
		})
	}
}
