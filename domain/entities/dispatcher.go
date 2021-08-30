package entities

import "github.com/matheusmosca/walrus/domain/vos"

type dispatcher struct {
	subscribers []Subscriber
	newMessage  chan vos.Message
	newSub      chan Subscriber
}

type Dispatcher interface {
	AddSubscriber(Subscriber)
	// TODO RemoveSubscriber
	// RemoveSubscriber(Subscriber)
	Dispatch(vos.Message)
	Activate()
}

func NewDispatcher() Dispatcher {
	return dispatcher{
		subscribers: []Subscriber{},
		newMessage:  make(chan vos.Message),
		newSub:      make(chan Subscriber),
	}
}

func (d dispatcher) AddSubscriber(sub Subscriber) {
	d.newSub <- sub
}

func (d dispatcher) Dispatch(message vos.Message) {
	d.newMessage <- message
}

func (d dispatcher) Activate() {
	go d.listenForSubscriptions()
	go d.listenForMessages()
}

func (d *dispatcher) listenForMessages() {
	for msg := range d.newMessage {
		m := msg

		for _, sub := range d.subscribers {
			sub.ReceiveMessage(m)
		}
	}
}

func (d *dispatcher) listenForSubscriptions() {
	for newSub := range d.newSub {
		d.subscribers = append(d.subscribers, newSub)
	}
}
