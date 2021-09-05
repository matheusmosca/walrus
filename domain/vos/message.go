package vos

type Message struct {
	TopicName   TopicName
	PublishedBy string
	Body        []byte
}

func (m Message) Validate() error {
	if m.PublishedBy == "" {
		return ErrEmptyPublishedBy
	}
	if len(m.PublishedBy) < 3 {
		return ErrPublishedByTooShort
	}
	if err := m.TopicName.Validate(); err != nil {
		return err
	}

	return nil
}
