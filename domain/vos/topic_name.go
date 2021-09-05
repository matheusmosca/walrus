package vos

type TopicName string

// TODO add more validations
func (t TopicName) Validate() error {
	if t == "" {
		return ErrEmptyTopicName
	}
	if len(t) < 3 {
		return ErrTopicNameTooShort
	}

	return nil
}

func (t TopicName) String() string {
	return string(t)
}
