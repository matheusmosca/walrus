package rpc

import (
	"io"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/matheusmosca/walrus/domain/entities"
	pb "github.com/matheusmosca/walrus/proto"
)

func (r RPCServer) Subscribe(stream pb.Walrus_SubscribeServer) error {
	subscriber := entities.NewSubscriber(uuid.New().String(), r.dispatcher)
	r.dispatcher.AddSubscriber(subscriber)

	newMessages := subscriber.Subscribe()
	for msg := range newMessages {

		res := &pb.SubscribeResponse{
			Message: &pb.Message{
				Topic:     msg.Topic,
				CreatedBy: msg.CreatedBy,
				Body:      msg.Body,
			},
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		logrus.Info(in)
	}

	return nil
}
