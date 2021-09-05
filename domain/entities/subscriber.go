package entities

import "github.com/matheusmosca/walrus/domain/vos"

type subscriber struct {
	id       string
	subChann chan vos.Message
	topic    Topic
}

type Subscriber interface {
	Subscribe() <-chan vos.Message
	ReceiveMessage(vos.Message)
}

func NewSubscriber(id string, topic Topic) Subscriber {
	sub := subscriber{
		id:       id,
		subChann: make(chan vos.Message),
		topic:    topic,
	}

	return sub
}

func (s subscriber) Subscribe() <-chan vos.Message {
	s.topic.AddSubscriber(s)
	return s.subChann
}

func (s subscriber) ReceiveMessage(msg vos.Message) {
	s.subChann <- msg
}
