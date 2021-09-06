package entities

import (
	"github.com/google/uuid"

	"github.com/matheusmosca/walrus/domain/vos"
)

type subscriber struct {
	id             string
	subscriptionCh chan vos.Message
	topic          Topic
}

type Subscriber interface {
	GetID() string
	ReceiveMessage(vos.Message)
	Subscribe() (chan vos.Message, string)
}

func NewSubscriber(topic Topic) Subscriber {
	sub := subscriber{
		id:             uuid.NewString(),
		subscriptionCh: make(chan vos.Message),
		topic:          topic,
	}

	return sub
}

func (s subscriber) Subscribe() (chan vos.Message, string) {
	s.topic.AddSubscriber(s)
	return s.subscriptionCh, s.GetID()
}

func (s subscriber) ReceiveMessage(msg vos.Message) {
	s.subscriptionCh <- msg
}

func (s subscriber) GetID() string {
	return s.id
}
