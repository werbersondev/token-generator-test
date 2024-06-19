// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/werbersondev/token-generator-test/domain/model"
	"github.com/werbersondev/token-generator-test/domain/service"
	"sync"
)

// Ensure, that RequestTokenGenerationRepositoryMock does implement service.RequestTokenGenerationRepository.
// If this is not the case, regenerate this file with moq.
var _ service.RequestTokenGenerationRepository = &RequestTokenGenerationRepositoryMock{}

// RequestTokenGenerationRepositoryMock is a mock implementation of service.RequestTokenGenerationRepository.
//
//	func TestSomethingThatUsesRequestTokenGenerationRepository(t *testing.T) {
//
//		// make and configure a mocked service.RequestTokenGenerationRepository
//		mockedRequestTokenGenerationRepository := &RequestTokenGenerationRepositoryMock{
//			PublishRequestTokenGenerationFunc: func(ctx context.Context, request model.TokenGenerationRequest) error {
//				panic("mock out the PublishRequestTokenGeneration method")
//			},
//		}
//
//		// use mockedRequestTokenGenerationRepository in code that requires service.RequestTokenGenerationRepository
//		// and then make assertions.
//
//	}
type RequestTokenGenerationRepositoryMock struct {
	// PublishRequestTokenGenerationFunc mocks the PublishRequestTokenGeneration method.
	PublishRequestTokenGenerationFunc func(ctx context.Context, request model.TokenGenerationRequest) error

	// calls tracks calls to the methods.
	calls struct {
		// PublishRequestTokenGeneration holds details about calls to the PublishRequestTokenGeneration method.
		PublishRequestTokenGeneration []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request model.TokenGenerationRequest
		}
	}
	lockPublishRequestTokenGeneration sync.RWMutex
}

// PublishRequestTokenGeneration calls PublishRequestTokenGenerationFunc.
func (mock *RequestTokenGenerationRepositoryMock) PublishRequestTokenGeneration(ctx context.Context, request model.TokenGenerationRequest) error {
	callInfo := struct {
		Ctx     context.Context
		Request model.TokenGenerationRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockPublishRequestTokenGeneration.Lock()
	mock.calls.PublishRequestTokenGeneration = append(mock.calls.PublishRequestTokenGeneration, callInfo)
	mock.lockPublishRequestTokenGeneration.Unlock()
	if mock.PublishRequestTokenGenerationFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.PublishRequestTokenGenerationFunc(ctx, request)
}

// PublishRequestTokenGenerationCalls gets all the calls that were made to PublishRequestTokenGeneration.
// Check the length with:
//
//	len(mockedRequestTokenGenerationRepository.PublishRequestTokenGenerationCalls())
func (mock *RequestTokenGenerationRepositoryMock) PublishRequestTokenGenerationCalls() []struct {
	Ctx     context.Context
	Request model.TokenGenerationRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request model.TokenGenerationRequest
	}
	mock.lockPublishRequestTokenGeneration.RLock()
	calls = mock.calls.PublishRequestTokenGeneration
	mock.lockPublishRequestTokenGeneration.RUnlock()
	return calls
}
