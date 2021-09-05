package rpc

import (
	"github.com/google/uuid"

	"github.com/matheusmosca/walrus/domain/entities"
	pb "github.com/matheusmosca/walrus/proto"
)

func (r RPC) Subscribe(req *pb.SubscribeRequest, stream pb.Walrus_SubscribeServer) error {
	subscriber := entities.NewSubscriber(uuid.New().String(), r.dispatcher)
	r.dispatcher.AddSubscriber(subscriber)

	newMessages := subscriber.Subscribe()
	for msg := range newMessages {

		res := &pb.SubscribeResponse{
			Message: &pb.Message{
				Topic:       msg.Topic,
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
