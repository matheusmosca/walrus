package rpc

import (
	"errors"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func injectStatusCode(err error) error {
	if isAlreadyExistsError(err) {
		return status.Error(codes.AlreadyExists, err.Error())
	}
	if isNotFoundError(err) {
		return status.Error(codes.NotFound, err.Error())
	}
	if isInvalidArgumentError(err) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return status.Error(codes.Internal, "internal server error")
}

func isNotFoundError(err error) bool {
	return errors.Is(err, entities.ErrTopicNotFound)
}

func isAlreadyExistsError(err error) bool {
	return errors.Is(err, entities.ErrTopicAlreadyExists)
}

func isInvalidArgumentError(err error) bool {
	if errors.Is(err, vos.ErrEmptyPublishedBy) {
		return true
	}
	if errors.Is(err, vos.ErrEmptyTopicName) {
		return true
	}
	if errors.Is(err, vos.ErrPublishedByTooShort) {
		return true
	}
	if errors.Is(err, vos.ErrTopicNameTooShort) {
		return true
	}

	return false
}
