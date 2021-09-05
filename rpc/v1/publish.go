package rpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
	"github.com/sirupsen/logrus"
)

func (r RPC) Publish(ctx context.Context, req *pb.PublishRequest) (*emptypb.Empty, error) {
	const operation = "RPC.Publish"

	log := r.log.WithContext(ctx).WithFields(logrus.Fields{
		"publishes_by": req.Message.PublishedBy,
		"topic_name":   req.Message.Topic,
		"operation":    operation,
	})
	log.Info("starting to handle new publish request")

	msg := vos.Message{
		PublishedBy: req.Message.PublishedBy,
		TopicName:   vos.TopicName(req.Message.Topic),
		Body:        req.Message.Body,
	}

	err := r.useCase.Publish(ctx, msg)
	if err != nil {
		log.WithError(err).Error("could not handle publish request")
		return nil, err
	}
	log.Info("message has been published successfully")

	return &emptypb.Empty{}, nil
}
