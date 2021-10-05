package rpc

import (
	"context"
	"net"
	"testing"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/usecases"
	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const buffer = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
func TestPublish(t *testing.T) {
	type args struct {
		ctx     context.Context
		message pb.Message
	}

	type fields struct {
		storage   usecases.Repository
		topicName vos.TopicName
	}

	testCases := []struct {
		name      string
		args      args
		fields    fields
		beforeRun func(topic entities.Topic) chan vos.Message
		want      []byte
		wantErr   error
	}{
		{
			name: "publish should succeed",
			args: args{
				ctx: context.Background(),
				message: pb.Message{
					Topic:       "walrus",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &usecases.RepositoryMock{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
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
				message: pb.Message{
					Topic:       "xd",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &usecases.RepositoryMock{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    []byte{},
			wantErr: injectStatusCode(vos.ErrTopicNameTooShort),
		},
		{
			name: "empty topic name should error",
			args: args{
				ctx: context.Background(),
				message: pb.Message{
					Topic:       "",
					PublishedBy: "walrus_test",
					Body:        []byte("hello again"),
				},
			},
			fields: fields{
				storage:   &usecases.RepositoryMock{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    []byte{},
			wantErr: injectStatusCode(vos.ErrEmptyTopicName),
		},
		{
			name: "short published by message should error",
			args: args{
				ctx: context.Background(),
				message: pb.Message{
					Topic:       "walrus",
					PublishedBy: "fb",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				storage:   &usecases.RepositoryMock{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    []byte{},
			wantErr: injectStatusCode(vos.ErrPublishedByTooShort),
		},
		{
			name: "empty published by message should error",
			args: args{
				ctx: context.Background(),
				message: pb.Message{
					Topic:       "walrus",
					PublishedBy: "",
					Body:        []byte("helloorld"),
				},
			},
			fields: fields{
				storage:   &usecases.RepositoryMock{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    []byte{},
			wantErr: injectStatusCode(vos.ErrEmptyPublishedBy),
		},
		{
			name: "nonexistent topic should return error",
			args: args{
				ctx: context.Background(),
				message: pb.Message{
					Topic:       "elephant",
					PublishedBy: "noone",
					Body:        []byte("bye world"),
				},
			},
			fields: fields{
				storage:   &usecases.RepositoryMock{},
				topicName: "walrus",
			},
			beforeRun: func(topic entities.Topic) chan vos.Message {
				subscriber := entities.NewSubscriber(topic)
				ch, _ := subscriber.Subscribe()
				return ch
			},
			want:    []byte{},
			wantErr: injectStatusCode(entities.ErrTopicNotFound),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// setup server
			ctx := context.Background()
			lis = bufconn.Listen(buffer)
			srv := grpc.NewServer()

			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
			require.NoError(t, err)
			defer conn.Close()

			// init topic
			topic, err := entities.NewTopic(tt.fields.topicName)
			require.NoError(t, err)

			// activate topic
			topic.Activate()

			// subscribe to topic
			subCh := tt.beforeRun(topic)
			defer close(subCh)

			// mock repository
			tt.fields.storage = &usecases.RepositoryMock{
				GetTopicFunc: func(ctx context.Context, topicName vos.TopicName) (entities.Topic, error) {
					if tt.fields.topicName != vos.TopicName(tt.args.message.Topic) {
						return entities.Topic{}, entities.ErrTopicNotFound
					}
					return topic, nil
				},
			}

			// create usecase
			useCase := usecases.New(tt.fields.storage)
			useCaseRPC := New(useCase, logrus.NewEntry(&logrus.Logger{}))

			// serve listener
			pb.RegisterWalrusServer(srv, useCaseRPC)
			go func() {
				err := srv.Serve(lis)
				require.NoError(t, err)
			}()

			// init client
			client := pb.NewWalrusClient(conn)

			// publish to topic
			_, err = client.Publish(tt.args.ctx, &pb.PublishRequest{Message: &tt.args.message})
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)

			// check published message
			actualMsg := <-subCh
			require.Equal(t, tt.want, actualMsg.Body)
		})

	}
}
