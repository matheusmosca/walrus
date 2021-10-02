package topics

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

func TestListTopics(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	type testCase struct {
		description string
		args        args
		beforeRun   func(storage map[vos.TopicName]entities.Topic)
		want        []vos.TopicName
		wantErr     error
	}

	tests := []testCase{
		{
			description: "should return a list of topics",
			args: args{
				context.Background(),
			},
			wantErr: nil,
			want: []vos.TopicName{
				vos.TopicName("topic1"),
				vos.TopicName("topic2"),
				vos.TopicName("topic3"),
			},
			beforeRun: func(storage map[vos.TopicName]entities.Topic) {
				t1, _ := entities.NewTopic(vos.TopicName("topic1"))
				t2, _ := entities.NewTopic(vos.TopicName("topic2"))
				t3, _ := entities.NewTopic(vos.TopicName("topic3"))
				storage[t1.GetName()] = t1
				storage[t2.GetName()] = t2
				storage[t3.GetName()] = t3
			},
		},
		{
			description: "should return a topics not found error when storage is empty",
			args: args{
				context.Background(),
			},
			wantErr:   entities.ErrTopicsNotFound,
			want:      nil,
			beforeRun: func(storage map[vos.TopicName]entities.Topic) {},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()

			storage := make(map[vos.TopicName]entities.Topic)
			tt.beforeRun(storage)
			repository := NewMemoryRepository(storage)

			topics, err := repository.ListTopics(tt.args.ctx)
			assert.ErrorIs(t, err, tt.wantErr)

			got := sortTopics(t, topics)
			assert.Equal(t, tt.want, got)
		})
	}
}

func sortTopics(t *testing.T, topics []entities.Topic) []vos.TopicName {
	t.Helper()

	if len(topics) == 0 {
		return nil
	}

	names := make([]vos.TopicName, 0, len(topics))
	sort.Slice(topics, func(i, j int) bool {
		return topics[i].GetName() < topics[j].GetName()
	})

	for _, topic := range topics {
		names = append(names, topic.GetName())
	}

	return names
}
