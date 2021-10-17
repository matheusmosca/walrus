package rpc

import (
	"context"
	"fmt"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

func TestPublish(t *testing.T) {
	type args struct {
		ctx     context.Context
		message vos.Message
	}

	type fields struct {
		topicName vos.TopicName
	}

	type grpcError struct {
		code codes.Code
		msg  error
	}

	testCases := []struct {
		name    string
		args    args
		fields  fields
		wantErr *grpcError
	}{
		{
			name: "publish should succeed",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "walrus",
					PublishedBy: "walrus_test",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				topicName: "walrus",
			},
			wantErr: &grpcError{},
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
				topicName: "walrus",
			},
			wantErr: &grpcError{
				code: codes.InvalidArgument,
				msg:  vos.ErrTopicNameTooShort,
			},
		},
		{
			name: "empty topic name should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "",
					PublishedBy: "walrus_test",
					Body:        []byte("hello again"),
				},
			},
			fields: fields{
				topicName: "walrus",
			},
			wantErr: &grpcError{
				code: codes.InvalidArgument,
				msg:  vos.ErrEmptyTopicName,
			},
		},
		{
			name: "short published by message should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "walrus",
					PublishedBy: "fb",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				topicName: "walrus",
			},
			wantErr: &grpcError{
				code: codes.InvalidArgument,
				msg:  vos.ErrPublishedByTooShort,
			},
		},
		{
			name: "empty published by message should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "walrus",
					PublishedBy: "",
					Body:        []byte("helloorld"),
				},
			},
			fields: fields{
				topicName: "walrus",
			},
			wantErr: &grpcError{
				code: codes.InvalidArgument,
				msg:  vos.ErrEmptyPublishedBy,
			},
		},
		{
			name: "nonexistent topic should return error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "elephant",
					PublishedBy: "noone",
					Body:        []byte("bye world"),
				},
			},
			fields: fields{
				topicName: "walrus",
			},
			wantErr: &grpcError{
				code: codes.NotFound,
				msg:  entities.ErrTopicNotFound,
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// setup server
			buffer := 1024 * 1024
			ctx := context.Background()
			lis := bufconn.Listen(buffer)
			srv := grpc.NewServer()

			conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}), grpc.WithInsecure())
			require.NoError(t, err)
			defer conn.Close()

			// mock use case
			mockedUseCase := &usecases.UseCaseMock{
				PublishFunc: func(ctx context.Context, message vos.Message) error {
					if err := message.Validate(); err != nil {
						return err
					}
					if tt.fields.topicName != tt.args.message.TopicName {
						return entities.ErrTopicNotFound
					}
					return nil
				},
			}

			// create usecase
			useCaseRPC := New(mockedUseCase, logrus.NewEntry(&logrus.Logger{}))

			// serve listener
			pb.RegisterWalrusServer(srv, useCaseRPC)
			go func() {
				err := srv.Serve(lis)
				require.NoError(t, err)
			}()

			// init client
			client := pb.NewWalrusClient(conn)

			// publish to topic
			publishRequest := pb.PublishRequest{
				Message: &pb.Message{
					Topic:       tt.args.message.TopicName.String(),
					PublishedBy: tt.args.message.PublishedBy,
					Body:        tt.args.message.Body,
				},
			}
			_, err = client.Publish(tt.args.ctx, &publishRequest)
			if *tt.wantErr != (grpcError{}) {
				s, _ := status.FromError(err)
				assert.Equal(t, *tt.wantErr, grpcError{code: s.Code(), msg: fmt.Errorf("%s", s.Message())})
				return
			}
			assert.NoError(t, err)
		})
	}
}
