package entities

import (
	"sync"
	"testing"

	"github.com/matheusmosca/walrus/domain/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTopic(t *testing.T) {
	type args struct {
		topicName vos.TopicName
	}
	type testCase struct {
		description   string
		args          args
		wantErr       error
		wantTopicName vos.TopicName
	}

	tests := []testCase{
		{
			description: "should return a topic without errors",
			args: args{
				topicName: vos.TopicName("any"),
			},
			wantErr:       nil,
			wantTopicName: vos.TopicName("any"),
		},
		{
			description: "should validate the topic name and return an error if the validation fails",
			args: args{
				topicName: vos.TopicName(""),
			},
			wantErr: vos.ErrEmptyTopicName,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			got, err := NewTopic(tt.args.topicName)
			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				assert.Equal(t, tt.wantTopicName, got.name)
				assert.Empty(t, got.subscribers)
			}
		})
	}
}

func TestDispatch(t *testing.T) {
	type testCase struct {
		description         string
		messages            []vos.Message
		numberOfSubscribers int
	}

	tests := []testCase{
		{
			description:         "should dispatch the messages to all subscribers successfully",
			numberOfSubscribers: 5,
			messages: []vos.Message{
				{
					TopicName:   vos.TopicName("test-topic"),
					PublishedBy: "first-publisher",
					Body:        []byte("{\"name\": \"Jorge\"}"),
				},
				{
					TopicName:   vos.TopicName("test-topic"),
					PublishedBy: "second-publisher",
					Body:        []byte("{\"book\": \"And then there were none\"}"),
				},
				{
					TopicName:   vos.TopicName("test-topic"),
					PublishedBy: "third-publisher",
					Body:        []byte("{\"age\": \"12\"}"),
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			topic, err := NewTopic(vos.TopicName("test-topic"))
			require.Nil(t, err)
			topic.Activate()

			var wg sync.WaitGroup
			for i := 0; i < tt.numberOfSubscribers; i++ {
				wg.Add(1)
				go assertDispatchedMessages(t, topic, tt.messages, &wg)
			}

			wg.Wait()
			for _, msg := range tt.messages {
				topic.Dispatch(msg)
			}
		})
	}
}

func assertDispatchedMessages(t *testing.T, topic Topic, wantMessages []vos.Message, wg *sync.WaitGroup) {
	defer wg.Done()
	sub := NewSubscriber(topic)
	subCh, _ := sub.Subscribe()
	msgIndex := 0

	go func() {
		defer close(subCh)
		for msg := range subCh {
			assert.Equal(t, sub.topic.name, msg.TopicName)
			// assert the message content and order of the messages as well
			assert.Equal(t, wantMessages[msgIndex], msg)
			msgIndex++

			if msgIndex == len(wantMessages) {
				// assert if subscriber got all messages
				assert.Equal(t, len(wantMessages), msgIndex)
				return
			}
		}
	}()
}
