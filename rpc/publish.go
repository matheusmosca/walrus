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
		Topic:       req.Message.Topic,
		Body:        req.Message.Body,
	}

	r.dispatcher.Dispatch(msg)
	return &emptypb.Empty{}, nil
}
