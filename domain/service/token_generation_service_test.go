package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/werbersondev/token-generator-test/domain/model"
	"github.com/werbersondev/token-generator-test/domain/service"
	mock "github.com/werbersondev/token-generator-test/domain/service/mocks"
)

func TestTokenGenerationService_GenerateToken_Success(t *testing.T) {
	tests := []struct {
		name          string
		projectID     string
		repoSetup     func(*testing.T) service.TokenGenerationRepository
		expectedToken string
	}{
		{
			name:      "Generate Token Success",
			projectID: "valid-project-id",
			repoSetup: func(t *testing.T) service.TokenGenerationRepository {
				return &mock.TokenGenerationRepositoryMock{
					GenerateProjectAnalysisTokenFunc: func(ctx context.Context, projectID string, tokenName string) (string, error) {
						// Simulate successful token generation
						return "generated-token", nil
					},
				}
			},
			expectedToken: "generated-token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.repoSetup(t)

			request := model.TokenGenerationRequest{
				ProjectID: tt.projectID,
			}

			s := service.NewTokenGenerationService(repository)
			token, err := s.GenerateToken(context.Background(), request)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedToken, token)
		})
	}
}

func TestTokenGenerationService_GenerateToken_Failure(t *testing.T) {
	tests := []struct {
		name        string
		projectID   string
		repoSetup   func(*testing.T) service.TokenGenerationRepository
		expectedErr error
	}{
		{
			name:      "Empty Project ID",
			projectID: "",
			repoSetup: func(t *testing.T) service.TokenGenerationRepository {
				return &mock.TokenGenerationRepositoryMock{
					GenerateProjectAnalysisTokenFunc: func(ctx context.Context, projectID string, tokenName string) (string, error) {
						// it should not be called
						t.FailNow()
						return "", nil
					},
				}
			},
			expectedErr: errors.New("projectID cannot be blank"),
		},
		{
			name:      "Repository Error",
			projectID: "valid-project-id",
			repoSetup: func(t *testing.T) service.TokenGenerationRepository {
				return &mock.TokenGenerationRepositoryMock{
					GenerateProjectAnalysisTokenFunc: func(ctx context.Context, projectID string, tokenName string) (string, error) {
						return "", errors.New("failed to generate analysis token")
					},
				}
			},
			expectedErr: errors.New("failed to generate analysis token"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := tt.repoSetup(t)

			request := model.TokenGenerationRequest{
				ProjectID: tt.projectID,
			}

			s := service.NewTokenGenerationService(repository)
			token, err := s.GenerateToken(context.Background(), request)

			assert.ErrorContains(t, err, tt.expectedErr.Error())
			assert.Empty(t, token)
		})
	}
}
