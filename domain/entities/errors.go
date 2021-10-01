package entities

import "errors"

var (
	ErrInvalidMessage     = errors.New("this message does not belong to the topic")
	ErrTopicsNotFound     = errors.New("no topic was found")
	ErrTopicNotFound      = errors.New("there isn't any subscriber listening to this topic")
	ErrTopicAlreadyExists = errors.New("this topic already exists")
)
