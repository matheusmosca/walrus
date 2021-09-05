package vos

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicNameValidate(t *testing.T) {
	type args struct {
		topicName TopicName
	}
	type testCase struct {
		description string
		args        args
		wantErr     error
	}

	tests := []testCase{
		{
			description: "should return no error, valid topic name",
			args: args{
				topicName: TopicName("any"),
			},
			wantErr: nil,
		},
		{
			description: "should return an error because topic name is empty",
			args: args{
				topicName: TopicName(""),
			},
			wantErr: ErrEmptyTopicName,
		},
		{
			description: "should return an error because topic name has less than 3 characters",
			args: args{
				topicName: TopicName("an"),
			},
			wantErr: ErrTopicNameTooShort,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			gotErr := tt.args.topicName.Validate()
			assert.ErrorIs(t, gotErr, tt.wantErr)
		})
	}
}
