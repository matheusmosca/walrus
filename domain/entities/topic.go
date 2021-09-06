package entities

import (
	"sync"

	"github.com/matheusmosca/walrus/domain/vos"
)

type topic struct {
	name         vos.TopicName
	subscribers  sync.Map
	newMessageCh chan vos.Message
	killSubCh    chan string
	newSubCh     chan Subscriber
}

type Topic interface {
	Activate()
	Dispatch(vos.Message)
	AddSubscriber(Subscriber)
	RemoveSubscriber(subscriberID string)
}

func NewTopic(topicName vos.TopicName) (Topic, error) {
	if err := topicName.Validate(); err != nil {
		return nil, err
	}

	return topic{
		name:         topicName,
		subscribers:  sync.Map{},
		newMessageCh: make(chan vos.Message),
		newSubCh:     make(chan Subscriber),
		killSubCh:    make(chan string),
	}, nil
}

func (t topic) Activate() {
	go t.listenForSubscriptions()
	go t.listenForMessages()
	go t.listenForKills()
}

func (t topic) Dispatch(message vos.Message) {
	t.newMessageCh <- message
}

func (t topic) AddSubscriber(sub Subscriber) {
	t.newSubCh <- sub
}

func (t topic) RemoveSubscriber(subscriberID string) {
	t.killSubCh <- subscriberID
}

func (t *topic) listenForSubscriptions() {
	for newSubCh := range t.newSubCh {
		t.subscribers.Store(newSubCh.GetID(), newSubCh)
	}
}

func (t *topic) listenForKills() {
	for subscriberID := range t.killSubCh {
		t.subscribers.Delete(subscriberID)
	}
}

func (t *topic) listenForMessages() {
	for msg := range t.newMessageCh {
		m := msg

		t.subscribers.Range(func(key, value interface{}) bool {
			if key == nil || value == nil {
				return false
			}
			if subscriber, ok := value.(Subscriber); ok {
				subscriber.ReceiveMessage(m)
			}

			return true
		})
	}
}
