package entities

import "errors"

var (
	ErrTopicsNotFound     = errors.New("no topic was found")
	ErrTopicNotFound      = errors.New("there isn't any subscriber listening to this topic")
	ErrTopicAlreadyExists = errors.New("this topic already exists")
)
