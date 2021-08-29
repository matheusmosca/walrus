package rpc

import (
	"context"

	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
)

func (r RPCServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	// TODO review nil pointer errors
	msg := vos.Message{
		CreatedBy: req.Message.CreatedBy,
		Topic:     req.Message.Topic,
		Body:      req.Message.Body,
	}

	r.dispatcher.Dispatch(msg)
	return &pb.PublishResponse{
		Success: true,
		Message: "message sended successfully",
	}, nil
}
