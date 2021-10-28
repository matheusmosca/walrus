// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package usecases

import (
	"context"
	"sync"

	"github.com/matheusmosca/walrus/domain/entities"
	"github.com/matheusmosca/walrus/domain/vos"
)

// Ensure, that UseCaseMock does implement UseCase.
// If this is not the case, regenerate this file with moq.
var _ UseCase = &UseCaseMock{}

// UseCaseMock is a mock implementation of UseCase.
//
// 	func TestSomethingThatUsesUseCase(t *testing.T) {
//
// 		// make and configure a mocked UseCase
// 		mockedUseCase := &UseCaseMock{
// 			ListTopicsFunc: func(ctx context.Context) ([]entities.Topic, error) {
// 				panic("mock out the ListTopics method")
// 			},
// 			PublishFunc: func(ctx context.Context, message vos.Message) error {
// 				panic("mock out the Publish method")
// 			},
// 			SubscribeFunc: func(ctx context.Context, topicName vos.TopicName) (chan vos.Message, vos.SubscriberID, error) {
// 				panic("mock out the Subscribe method")
// 			},
// 			UnsubscribeFunc: func(ctx context.Context, subscriberID vos.SubscriberID, topicName vos.TopicName) error {
// 				panic("mock out the Unsubscribe method")
// 			},
// 		}
//
// 		// use mockedUseCase in code that requires UseCase
// 		// and then make assertions.
//
// 	}
type UseCaseMock struct {
	// ListTopicsFunc mocks the ListTopics method.
	ListTopicsFunc func(ctx context.Context) ([]entities.Topic, error)

	// PublishFunc mocks the Publish method.
	PublishFunc func(ctx context.Context, message vos.Message) error

	// SubscribeFunc mocks the Subscribe method.
	SubscribeFunc func(ctx context.Context, topicName vos.TopicName) (chan vos.Message, vos.SubscriberID, error)

	// UnsubscribeFunc mocks the Unsubscribe method.
	UnsubscribeFunc func(ctx context.Context, subscriberID vos.SubscriberID, topicName vos.TopicName) error

	// calls tracks calls to the methods.
	calls struct {
		// ListTopics holds details about calls to the ListTopics method.
		ListTopics []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Publish holds details about calls to the Publish method.
		Publish []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Message is the message argument value.
			Message vos.Message
		}
		// Subscribe holds details about calls to the Subscribe method.
		Subscribe []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// TopicName is the topicName argument value.
			TopicName vos.TopicName
		}
		// Unsubscribe holds details about calls to the Unsubscribe method.
		Unsubscribe []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// SubscriberID is the subscriberID argument value.
			SubscriberID vos.SubscriberID
			// TopicName is the topicName argument value.
			TopicName vos.TopicName
		}
	}
	lockListTopics  sync.RWMutex
	lockPublish     sync.RWMutex
	lockSubscribe   sync.RWMutex
	lockUnsubscribe sync.RWMutex
}

// ListTopics calls ListTopicsFunc.
func (mock *UseCaseMock) ListTopics(ctx context.Context) ([]entities.Topic, error) {
	if mock.ListTopicsFunc == nil {
		panic("UseCaseMock.ListTopicsFunc: method is nil but UseCase.ListTopics was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockListTopics.Lock()
	mock.calls.ListTopics = append(mock.calls.ListTopics, callInfo)
	mock.lockListTopics.Unlock()
	return mock.ListTopicsFunc(ctx)
}

// ListTopicsCalls gets all the calls that were made to ListTopics.
// Check the length with:
//     len(mockedUseCase.ListTopicsCalls())
func (mock *UseCaseMock) ListTopicsCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockListTopics.RLock()
	calls = mock.calls.ListTopics
	mock.lockListTopics.RUnlock()
	return calls
}

// Publish calls PublishFunc.
func (mock *UseCaseMock) Publish(ctx context.Context, message vos.Message) error {
	if mock.PublishFunc == nil {
		panic("UseCaseMock.PublishFunc: method is nil but UseCase.Publish was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Message vos.Message
	}{
		Ctx:     ctx,
		Message: message,
	}
	mock.lockPublish.Lock()
	mock.calls.Publish = append(mock.calls.Publish, callInfo)
	mock.lockPublish.Unlock()
	return mock.PublishFunc(ctx, message)
}

// PublishCalls gets all the calls that were made to Publish.
// Check the length with:
//     len(mockedUseCase.PublishCalls())
func (mock *UseCaseMock) PublishCalls() []struct {
	Ctx     context.Context
	Message vos.Message
} {
	var calls []struct {
		Ctx     context.Context
		Message vos.Message
	}
	mock.lockPublish.RLock()
	calls = mock.calls.Publish
	mock.lockPublish.RUnlock()
	return calls
}

// Subscribe calls SubscribeFunc.
func (mock *UseCaseMock) Subscribe(ctx context.Context, topicName vos.TopicName) (chan vos.Message, vos.SubscriberID, error) {
	if mock.SubscribeFunc == nil {
		panic("UseCaseMock.SubscribeFunc: method is nil but UseCase.Subscribe was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		TopicName vos.TopicName
	}{
		Ctx:       ctx,
		TopicName: topicName,
	}
	mock.lockSubscribe.Lock()
	mock.calls.Subscribe = append(mock.calls.Subscribe, callInfo)
	mock.lockSubscribe.Unlock()
	return mock.SubscribeFunc(ctx, topicName)
}

// SubscribeCalls gets all the calls that were made to Subscribe.
// Check the length with:
//     len(mockedUseCase.SubscribeCalls())
func (mock *UseCaseMock) SubscribeCalls() []struct {
	Ctx       context.Context
	TopicName vos.TopicName
} {
	var calls []struct {
		Ctx       context.Context
		TopicName vos.TopicName
	}
	mock.lockSubscribe.RLock()
	calls = mock.calls.Subscribe
	mock.lockSubscribe.RUnlock()
	return calls
}

// Unsubscribe calls UnsubscribeFunc.
func (mock *UseCaseMock) Unsubscribe(ctx context.Context, subscriberID vos.SubscriberID, topicName vos.TopicName) error {
	if mock.UnsubscribeFunc == nil {
		panic("UseCaseMock.UnsubscribeFunc: method is nil but UseCase.Unsubscribe was just called")
	}
	callInfo := struct {
		Ctx          context.Context
		SubscriberID vos.SubscriberID
		TopicName    vos.TopicName
	}{
		Ctx:          ctx,
		SubscriberID: subscriberID,
		TopicName:    topicName,
	}
	mock.lockUnsubscribe.Lock()
	mock.calls.Unsubscribe = append(mock.calls.Unsubscribe, callInfo)
	mock.lockUnsubscribe.Unlock()
	return mock.UnsubscribeFunc(ctx, subscriberID, topicName)
}

// UnsubscribeCalls gets all the calls that were made to Unsubscribe.
// Check the length with:
//     len(mockedUseCase.UnsubscribeCalls())
func (mock *UseCaseMock) UnsubscribeCalls() []struct {
	Ctx          context.Context
	SubscriberID vos.SubscriberID
	TopicName    vos.TopicName
} {
	var calls []struct {
		Ctx          context.Context
		SubscriberID vos.SubscriberID
		TopicName    vos.TopicName
	}
	mock.lockUnsubscribe.RLock()
	calls = mock.calls.Unsubscribe
	mock.lockUnsubscribe.RUnlock()
	return calls
}