package vos

import "errors"

var (
	ErrEmptyPublishedBy    = errors.New("published by field could not be empty")
	ErrPublishedByTooShort = errors.New("published by field must have at least 3 characters")
	ErrEmptyTopicName      = errors.New("topic name field could not be empty")
	ErrTopicNameTooShort   = errors.New("topic name must have at least 3 characters")
)
