package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/werbersondev/token-generator-test/domain/model"
	"github.com/werbersondev/token-generator-test/domain/service"
	"github.com/werbersondev/token-generator-test/domain/service/mocks"
)

func TestRequestTokenGenerationService_RequestTokenGeneration_Success(t *testing.T) {
	tests := []struct {
		name      string
		projectID string
		repoSetup func(*testing.T) service.RequestTokenGenerationRepository
	}{
		{
			name:      "Request Token Generation Success",
			projectID: "valid-project-id",
			repoSetup: func(t *testing.T) service.RequestTokenGenerationRepository {
				return &mocks.RequestTokenGenerationRepositoryMock{
					PublishRequestTokenGenerationFunc: func(ctx context.Context, request model.TokenGenerationRequest) error {
						// Simulate successful request token generation
						return nil
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.repoSetup(t)

			s := service.NewRequestTokenGenerationService(repository)
			err := s.RequestTokenGeneration(context.Background(), tt.projectID)

			assert.NoError(t, err)
		})
	}
}

func TestRequestTokenGenerationService_RequestTokenGeneration_Failure(t *testing.T) {
	tests := []struct {
		name        string
		projectID   string
		repoSetup   func(*testing.T) service.RequestTokenGenerationRepository
		expectedErr error
	}{
		{
			name:      "Empty Project ID",
			projectID: "",
			repoSetup: func(t *testing.T) service.RequestTokenGenerationRepository {
				return &mocks.RequestTokenGenerationRepositoryMock{}
			},
			expectedErr: errors.New("projectID cannot be blank"),
		},
		{
			name:      "Repository Error",
			projectID: "valid-project-id",
			repoSetup: func(t *testing.T) service.RequestTokenGenerationRepository {
				return &mocks.RequestTokenGenerationRepositoryMock{
					PublishRequestTokenGenerationFunc: func(ctx context.Context, request model.TokenGenerationRequest) error {
						return errors.New("failed to publish request token")
					},
				}
			},
			expectedErr: errors.New("publishing request token generation: failed to publish request token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.repoSetup(t)

			s := service.NewRequestTokenGenerationService(repository)
			err := s.RequestTokenGeneration(context.Background(), tt.projectID)

			assert.ErrorContains(t, err, tt.expectedErr.Error())
		})
	}
}
