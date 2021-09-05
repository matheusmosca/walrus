package vos

type Message struct {
	PublishedBy string
	Topic       string
	Body        []byte
}

func (m Message) Validate() error {
	if m.PublishedBy == "" {
		return ErrEmptyPublishedBy
	}
	if len(m.PublishedBy) < 3 {
		return ErrPublishedByTooShort
	}
	if m.Topic == "" {
		return ErrEmptyTopicName
	}
	if len(m.Topic) < 3 {
		return ErrTopicNameTooShort
	}

	return nil
}
