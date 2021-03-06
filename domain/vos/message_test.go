package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageValidate(t *testing.T) {
	type args struct {
		message Message
	}
	type testCase struct {
		description string
		args        args
		wantErr     error
	}

	tests := []testCase{
		{
			description: "should return no error, valid message",
			args: args{
				message: Message{
					PublishedBy: "any",
					TopicName:   "any",
					Body:        make([]byte, 0),
				},
			},
			wantErr: nil,
		},
		{
			description: "should return an error because published by is empty",
			args: args{
				message: Message{
					PublishedBy: "",
					TopicName:   "any",
					Body:        make([]byte, 0),
				},
			},
			wantErr: ErrEmptyPublishedBy,
		},
		{
			description: "should return an error because published by field has less than 3 characters",
			args: args{
				message: Message{
					PublishedBy: "an",
					TopicName:   "any",
					Body:        make([]byte, 0),
				},
			},
			wantErr: ErrPublishedByTooShort,
		},
		{
			description: "should return an error because topic name is empty",
			args: args{
				message: Message{
					PublishedBy: "publisher-name",
					TopicName:   "",
					Body:        make([]byte, 0),
				},
			},
			wantErr: ErrEmptyTopicName,
		},
		{
			description: "should return an error because topic name has less than 3 characters",
			args: args{
				message: Message{
					PublishedBy: "publisher-name",
					TopicName:   "an",
					Body:        make([]byte, 0),
				},
			},
			wantErr: ErrTopicNameTooShort,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			gotErr := tt.args.message.Validate()
			assert.ErrorIs(t, gotErr, tt.wantErr)
		})
	}
}
