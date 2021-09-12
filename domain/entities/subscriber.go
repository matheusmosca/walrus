package entities

import (
	"github.com/google/uuid"

	"github.com/matheusmosca/walrus/domain/vos"
)

type Subscriber struct {
	id             vos.SubscriberID
	subscriptionCh chan vos.Message
	topic          Topic
}

func NewSubscriber(topic Topic) Subscriber {
	return Subscriber{
		id:             vos.SubscriberID(uuid.NewString()),
		subscriptionCh: make(chan vos.Message),
		topic:          topic,
	}
}

func (s Subscriber) Subscribe() (chan vos.Message, vos.SubscriberID) {
	s.topic.addSubscriber(s)
	return s.subscriptionCh, s.GetID()
}

func (s Subscriber) ReceiveMessage(msg vos.Message) {
	s.subscriptionCh <- msg
}

func (s Subscriber) GetID() vos.SubscriberID {
	return s.id
}
