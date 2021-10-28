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
		useCaseMock *usecases.UseCaseMock
	}

	testCases := []struct {
		name     string
		args     args
		fields   fields
		wantErr  error
		wantCode codes.Code
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
				useCaseMock: &usecases.UseCaseMock{
					PublishFunc: func(ctx context.Context, message vos.Message) error {
						return nil
					},
				},
			},
			wantErr:  nil,
			wantCode: codes.OK,
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
				useCaseMock: &usecases.UseCaseMock{
					PublishFunc: func(ctx context.Context, message vos.Message) error {
						return vos.ErrTopicNameTooShort
					},
				},
			},
			wantErr:  injectStatusCode(vos.ErrTopicNameTooShort),
			wantCode: codes.InvalidArgument,
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
				useCaseMock: &usecases.UseCaseMock{PublishFunc: func(ctx context.Context, message vos.Message) error {
					return vos.ErrEmptyTopicName
				}},
			},
			wantErr:  injectStatusCode(vos.ErrEmptyTopicName),
			wantCode: codes.InvalidArgument,
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
				useCaseMock: &usecases.UseCaseMock{
					PublishFunc: func(ctx context.Context, message vos.Message) error {
						return vos.ErrPublishedByTooShort
					},
				},
			},
			wantErr:  injectStatusCode(vos.ErrPublishedByTooShort),
			wantCode: codes.InvalidArgument,
		},
		{
			name: "empty published by message should error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "walrus",
					PublishedBy: "",
					Body:        []byte("hello world"),
				},
			},
			fields: fields{
				useCaseMock: &usecases.UseCaseMock{
					PublishFunc: func(ctx context.Context, message vos.Message) error {
						return vos.ErrEmptyPublishedBy
					},
				},
			},
			wantErr:  injectStatusCode(vos.ErrEmptyPublishedBy),
			wantCode: codes.InvalidArgument,
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
				useCaseMock: &usecases.UseCaseMock{
					PublishFunc: func(ctx context.Context, message vos.Message) error {
						return entities.ErrTopicNotFound
					},
				},
			},
			wantErr:  injectStatusCode(entities.ErrTopicNotFound),
			wantCode: codes.NotFound,
		},
		{
			name: "unexpected error should give internal error",
			args: args{
				ctx: context.Background(),
				message: vos.Message{
					TopicName:   "elephant",
					PublishedBy: "noone",
					Body:        []byte("bye world"),
				},
			},
			fields: fields{
				useCaseMock: &usecases.UseCaseMock{
					PublishFunc: func(ctx context.Context, message vos.Message) error {
						return fmt.Errorf("something bad happened")
					},
				},
			},
			wantErr:  injectStatusCode(fmt.Errorf("internal server error")),
			wantCode: codes.Internal,
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			// setup server
			buffer := 1024 * 1024
			lis := bufconn.Listen(buffer)
			srv := grpc.NewServer()

			conn, err := grpc.DialContext(tt.args.ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}), grpc.WithInsecure())
			require.NoError(t, err)
			defer conn.Close()

			// create usecase
			useCaseRPC := New(tt.fields.useCaseMock, logrus.NewEntry(&logrus.Logger{}))

			// serve listener
			pb.RegisterWalrusServer(srv, useCaseRPC)
			go func() {
				err := srv.Serve(lis)
				require.NoError(t, err)
			}()

			// init client
			client := pb.NewWalrusClient(conn)

			// publish to topic
			_, err = client.Publish(tt.args.ctx, &pb.PublishRequest{
				Message: &pb.Message{
					Topic:       tt.args.message.TopicName.String(),
					PublishedBy: tt.args.message.PublishedBy,
					Body:        tt.args.message.Body,
				},
			})
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr != nil {
				gotCode, _ := status.FromError(err)
				assert.Equal(t, tt.wantCode, gotCode.Code())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantCode, codes.OK)
		})
	}
}
