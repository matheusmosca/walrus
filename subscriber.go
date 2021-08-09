package walrus

import (
	"log"
)

type subscriber struct {
	id         int
	subChann   chan []byte
	dispatcher Dispatcher
}

type Subscriber interface {
	Subscribe()
	ReceiveMessage([]byte)
}

func NewSubscriber(id int, dispatcher Dispatcher) Subscriber {
	sub := subscriber{
		id:         id,
		subChann:   make(chan []byte),
		dispatcher: dispatcher,
	}

	return sub
}

func (s subscriber) Subscribe() {
	s.dispatcher.AddSubscriber(s)
	go s.listenForMessages()
}

func (s subscriber) listenForMessages() {
	for msg := range s.subChann {
		log.Printf("Subscriber %d got %v\n", s.id, string(msg))
	}
}

func (s subscriber) ReceiveMessage(message []byte) {
	s.subChann <- message
}
