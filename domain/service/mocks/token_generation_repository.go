// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"context"
	"github.com/werbersondev/token-generator-test/domain/service"
	"sync"
)

// Ensure, that TokenGenerationRepositoryMock does implement service.TokenGenerationRepository.
// If this is not the case, regenerate this file with moq.
var _ service.TokenGenerationRepository = &TokenGenerationRepositoryMock{}

// TokenGenerationRepositoryMock is a mock implementation of service.TokenGenerationRepository.
//
//	func TestSomethingThatUsesTokenGenerationRepository(t *testing.T) {
//
//		// make and configure a mocked service.TokenGenerationRepository
//		mockedTokenGenerationRepository := &TokenGenerationRepositoryMock{
//			GenerateProjectAnalysisTokenFunc: func(ctx context.Context, projectID string, tokenName string) (string, error) {
//				panic("mock out the GenerateProjectAnalysisToken method")
//			},
//		}
//
//		// use mockedTokenGenerationRepository in code that requires service.TokenGenerationRepository
//		// and then make assertions.
//
//	}
type TokenGenerationRepositoryMock struct {
	// GenerateProjectAnalysisTokenFunc mocks the GenerateProjectAnalysisToken method.
	GenerateProjectAnalysisTokenFunc func(ctx context.Context, projectID string, tokenName string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// GenerateProjectAnalysisToken holds details about calls to the GenerateProjectAnalysisToken method.
		GenerateProjectAnalysisToken []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ProjectID is the projectID argument value.
			ProjectID string
			// TokenName is the tokenName argument value.
			TokenName string
		}
	}
	lockGenerateProjectAnalysisToken sync.RWMutex
}

// GenerateProjectAnalysisToken calls GenerateProjectAnalysisTokenFunc.
func (mock *TokenGenerationRepositoryMock) GenerateProjectAnalysisToken(ctx context.Context, projectID string, tokenName string) (string, error) {
	callInfo := struct {
		Ctx       context.Context
		ProjectID string
		TokenName string
	}{
		Ctx:       ctx,
		ProjectID: projectID,
		TokenName: tokenName,
	}
	mock.lockGenerateProjectAnalysisToken.Lock()
	mock.calls.GenerateProjectAnalysisToken = append(mock.calls.GenerateProjectAnalysisToken, callInfo)
	mock.lockGenerateProjectAnalysisToken.Unlock()
	if mock.GenerateProjectAnalysisTokenFunc == nil {
		var (
			sOut   string
			errOut error
		)
		return sOut, errOut
	}
	return mock.GenerateProjectAnalysisTokenFunc(ctx, projectID, tokenName)
}

// GenerateProjectAnalysisTokenCalls gets all the calls that were made to GenerateProjectAnalysisToken.
// Check the length with:
//
//	len(mockedTokenGenerationRepository.GenerateProjectAnalysisTokenCalls())
func (mock *TokenGenerationRepositoryMock) GenerateProjectAnalysisTokenCalls() []struct {
	Ctx       context.Context
	ProjectID string
	TokenName string
} {
	var calls []struct {
		Ctx       context.Context
		ProjectID string
		TokenName string
	}
	mock.lockGenerateProjectAnalysisToken.RLock()
	calls = mock.calls.GenerateProjectAnalysisToken
	mock.lockGenerateProjectAnalysisToken.RUnlock()
	return calls
}
