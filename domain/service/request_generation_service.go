package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/werbersondev/token-generator-test/domain/model"
)

type RequestTokenGenerationService struct {
	repository RequestTokenGenerationRepository
}

//go:generate moq -stub -pkg mocks -out mocks/request_generation_repository.go . RequestTokenGenerationRepository
type RequestTokenGenerationRepository interface {
	PublishRequestTokenGeneration(ctx context.Context, request model.TokenGenerationRequest) error
}

func NewRequestTokenGenerationService(repo RequestTokenGenerationRepository) *RequestTokenGenerationService {
	return &RequestTokenGenerationService{repository: repo}
}

func (r *RequestTokenGenerationService) RequestTokenGeneration(ctx context.Context, projectID string) error {
	if strings.TrimSpace(projectID) == "" {
		return errors.New("projectID cannot be blank")
	}
	err := r.repository.PublishRequestTokenGeneration(ctx, model.TokenGenerationRequest{ProjectID: projectID})
	if err != nil {
		return fmt.Errorf("publishing request token generation: %w", err)
	}

	return nil
}
