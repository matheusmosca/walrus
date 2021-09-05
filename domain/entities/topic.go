package entities

import "github.com/matheusmosca/walrus/domain/vos"

type topic struct {
	name        vos.TopicName
	subscribers []Subscriber
	newMessage  chan vos.Message
	newSub      chan Subscriber
}

type Topic interface {
	AddSubscriber(Subscriber)
	// TODO RemoveSubscriber
	// RemoveSubscriber(Subscriber)
	Dispatch(vos.Message)
	Activate()
}

func NewTopic(topicName vos.TopicName) (Topic, error) {
	if err := topicName.Validate(); err != nil {
		return nil, err
	}

	return topic{
		name:        topicName,
		subscribers: []Subscriber{},
		newMessage:  make(chan vos.Message),
		newSub:      make(chan Subscriber),
	}, nil
}

func (t topic) AddSubscriber(sub Subscriber) {
	t.newSub <- sub
}

func (t topic) Dispatch(message vos.Message) {
	t.newMessage <- message
}

func (t topic) Activate() {
	go t.listenForSubscriptions()
	go t.listenForMessages()
}

func (t *topic) listenForMessages() {
	for msg := range t.newMessage {
		m := msg

		for _, sub := range t.subscribers {
			sub.ReceiveMessage(m)
		}
	}
}

func (t *topic) listenForSubscriptions() {
	for newSub := range t.newSub {
		t.subscribers = append(t.subscribers, newSub)
	}
}
