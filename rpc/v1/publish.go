package rpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
)

func (r RPC) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	msg := vos.Message{
		PublishedBy: req.Message.PublishedBy,
		TopicName:   vos.TopicName(req.Message.Topic),
		Body:        req.Message.Body,
	}

	err := r.useCase.Publish(ctx, msg)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
