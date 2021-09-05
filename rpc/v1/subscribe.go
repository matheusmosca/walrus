package rpc

import (
	"context"

	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
	"github.com/sirupsen/logrus"
)

func (r RPC) Subscribe(req *pb.SubscribeRequest, stream pb.Walrus_SubscribeServer) error {
	const operation = "RPC.Subscribe"

	log := r.log.WithFields(logrus.Fields{
		"subscriber_id": req.SubscriberId,
		"topic_name":    req.Topic,
		"operation":     operation,
	})
	log.Info("starting to handle subscription request")

	subscriptionCh, err := r.useCase.Subscribe(context.TODO(), req.SubscriberId, vos.TopicName(req.Topic))
	if err != nil {
		log.WithError(err).Error("could not handle subscribe request")
		return err
	}
	log.Info("subscription request has been completed successfully, starting to receive messages")

	for msg := range subscriptionCh {
		res := &pb.SubscribeResponse{
			Message: &pb.Message{
				Topic:       msg.TopicName.String(),
				Body:        msg.Body,
				PublishedBy: msg.PublishedBy,
			},
		}
		log := log.WithField("published_by", msg.PublishedBy)
		log.Info("new message published, starting to send it to subscriber")
		if err := stream.Send(res); err != nil {
			log.WithError(err).Error("could not send message to subscriber")
			return err
		}

		log.Info("message has been sent successfully")
	}

	log.Info("subscription stream has been closed")
	return nil
}
