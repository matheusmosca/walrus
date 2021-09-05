package entities

import "github.com/matheusmosca/walrus/domain/vos"

type publisher struct {
	topic Topic
}

type Publisher interface {
	Publish(vos.Message)
}

func NewPublisher(t Topic) Publisher {
	return publisher{
		topic: t,
	}
}

func (p publisher) Publish(message vos.Message) {
	p.topic.Dispatch(message)
}
