package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

//go:generate moq -stub -pkg mocks -out mocks/request_generation_uc.go . RequestTokenGenerationUseCase
type RequestTokenGenerationUseCase interface {
	RequestTokenGeneration(ctx context.Context, projectID string) error
}

type RequestTokenGenerationInput struct {
	ProjectID string `json:"project_id"`
}

func RequestTokenGenerationHandler(uc RequestTokenGenerationUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var body RequestTokenGenerationInput
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if body.ProjectID == "" {
			http.Error(w, "Missing required parameter: project_id", http.StatusUnprocessableEntity)
			return
		}

		err := uc.RequestTokenGeneration(ctx, body.ProjectID)
		if err != nil {
			http.Error(w, "Failed to publish message", http.StatusInternalServerError)
			return
		}

		log.Ctx(ctx).Info().Str("project_id", body.ProjectID).Msg("Token generation request sent")

		w.WriteHeader(http.StatusAccepted)
	}
}
