package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
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
					Topic:       "any",
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
					Topic:       "any",
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
					Topic:       "any",
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
					Topic:       "",
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
					Topic:       "an",
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
