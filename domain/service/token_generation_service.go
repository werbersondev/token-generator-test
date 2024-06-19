package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/werbersondev/token-generator-test/domain/model"
)

type TokenGenerationService struct {
	repository TokenGenerationRepository
}

//go:generate moq -stub -pkg mocks -out mocks/token_generation_repository.go . TokenGenerationRepository
type TokenGenerationRepository interface {
	GenerateProjectAnalysisToken(ctx context.Context, projectID, tokenName string) (string, error)
}

func NewTokenGenerationService(repo TokenGenerationRepository) *TokenGenerationService {
	return &TokenGenerationService{repository: repo}
}

func (r *TokenGenerationService) GenerateToken(ctx context.Context, request model.TokenGenerationRequest) (string, error) {
	if strings.TrimSpace(request.ProjectID) == "" {
		return "", errors.New("projectID cannot be blank")
	}

	tokenName := fmt.Sprintf("%s-analysis-%s", request.ProjectID, time.Now().String())

	token, err := r.repository.GenerateProjectAnalysisToken(ctx, request.ProjectID, tokenName)
	if err != nil {
		return "", fmt.Errorf("generating token on provider: %w", err)
	}

	return token, nil
}
