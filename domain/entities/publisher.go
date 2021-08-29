package entities

import "github.com/matheusmosca/walrus/domain/vos"

type publisher struct {
	dispatcher Dispatcher
}

type Publisher interface {
	Publish(vos.Message)
}

func NewPublisher(d Dispatcher) Publisher {
	return publisher{
		dispatcher: d,
	}
}

func (p publisher) Publish(message vos.Message) {
	p.dispatcher.Dispatch(message)
}
