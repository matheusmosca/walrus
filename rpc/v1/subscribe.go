package rpc

import (
	"context"

	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
)

func (r RPC) Subscribe(req *pb.SubscribeRequest, stream pb.Walrus_SubscribeServer) error {
	subscriptionCh, err := r.useCase.Subscribe(context.TODO(), req.SubscriberId, vos.TopicName(req.Topic))
	if err != nil {
		return err
	}

	for msg := range subscriptionCh {
		res := &pb.SubscribeResponse{
			Message: &pb.Message{
				Topic:       msg.TopicName.String(),
				Body:        msg.Body,
				PublishedBy: msg.PublishedBy,
			},
		}
		if err := stream.Send(res); err != nil {
			return err
		}
	}

	return nil
}
