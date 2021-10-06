package entities

import (
	"fmt"
	"sync"

	"github.com/matheusmosca/walrus/domain/vos"
)

type Topic struct {
	name         vos.TopicName
	subscribers  *sync.Map
	newMessageCh chan vos.Message
	killSubCh    chan vos.SubscriberID
	newSubCh     chan Subscriber
}

func NewTopic(topicName vos.TopicName) (Topic, error) {
	if err := topicName.Validate(); err != nil {
		return Topic{}, err
	}

	return Topic{
		name:         topicName,
		subscribers:  &sync.Map{},
		newMessageCh: make(chan vos.Message),
		newSubCh:     make(chan Subscriber),
		killSubCh:    make(chan vos.SubscriberID),
	}, nil
}

func (t Topic) Activate() {
	go t.listenForSubscriptions()
	go t.listenForMessages()
}

func (t Topic) Dispatch(message vos.Message) error {
	if message.TopicName != t.name {
		return ErrTopicNameDoesNotMatch
	}

	t.newMessageCh <- message

	return nil
}

func (t *Topic) RemoveSubscriber(subscriberID vos.SubscriberID) error {
	if _, ok := t.subscribers.LoadAndDelete(subscriberID); !ok {
		return ErrSubscriberNotFound
	}

	return nil
}

func (t Topic) addSubscriber(sub Subscriber) {
	t.newSubCh <- sub
}

func (t Topic) GetSubscriber(subscriberID vos.SubscriberID) (*Subscriber, error) {
	if value, ok := t.subscribers.Load(subscriberID); ok {
		subsInterface := value.(Subscriber)

		return &subsInterface, nil

	}

	return nil, ErrSubscriberNotFound
}

func (t Topic) UpdateTopic(topic Topic) (Topic, error) {
	fmt.Println("implement me")

	return Topic{}, nil
}

func (t *Topic) listenForSubscriptions() {
	for newSubCh := range t.newSubCh {
		t.subscribers.Store(newSubCh.GetID(), newSubCh)
	}
}

func (t *Topic) listenForMessages() {
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

func (t Topic) GetName() vos.TopicName {
	return t.name
}
