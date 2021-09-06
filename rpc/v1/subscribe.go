package rpc

import (
	"github.com/matheusmosca/walrus/domain/vos"
	pb "github.com/matheusmosca/walrus/proto"
	"github.com/sirupsen/logrus"
)

func (r RPC) Subscribe(req *pb.SubscribeRequest, stream pb.Walrus_SubscribeServer) error {
	const operation = "RPC.Subscribe"

	log := r.log.WithContext(stream.Context()).WithFields(logrus.Fields{
		"topic_name": req.Topic,
		"operation":  operation,
	})
	log.Info("starting to handle subscription request")

	subscriptionCh, subscriberID, err := r.useCase.Subscribe(stream.Context(), vos.TopicName(req.Topic))
	if err != nil {
		log.WithError(err).Error("could not handle subscribe request")
		return err
	}
	defer close(subscriptionCh)
	log.Info("subscription request has been completed successfully, starting to receive messages")

	for {
		select {
		case msg := <-subscriptionCh:
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
		case <-stream.Context().Done():
			log.Info("subscription stream has been closed, starting to unsubscribe")
			err := r.useCase.Unsubscribe(stream.Context(), subscriberID, vos.TopicName(req.Topic))
			if err != nil {
				log.WithError(err).Error("something went wrong trying to unsubscribe")
				return err
			}

			log.Info("unsubscribed successfully")
			return nil
		}
	}
}
