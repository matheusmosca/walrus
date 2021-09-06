package entities

import (
	"github.com/google/uuid"

	"github.com/matheusmosca/walrus/domain/vos"
)

type subscriber struct {
	id             vos.SubscriberID
	subscriptionCh chan vos.Message
	topic          Topic
}

type Subscriber interface {
	GetID() vos.SubscriberID
	ReceiveMessage(vos.Message)
	Subscribe() (chan vos.Message, vos.SubscriberID)
}

func NewSubscriber(topic Topic) Subscriber {
	sub := subscriber{
		id:             vos.SubscriberID(uuid.NewString()),
		subscriptionCh: make(chan vos.Message),
		topic:          topic,
	}

	return sub
}

func (s subscriber) Subscribe() (chan vos.Message, vos.SubscriberID) {
	s.topic.AddSubscriber(s)
	return s.subscriptionCh, s.GetID()
}

func (s subscriber) ReceiveMessage(msg vos.Message) {
	s.subscriptionCh <- msg
}

func (s subscriber) GetID() vos.SubscriberID {
	return s.id
}
