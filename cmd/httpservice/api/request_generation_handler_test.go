package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/werbersondev/token-generator-test/cmd/httpservice/api"
	"github.com/werbersondev/token-generator-test/cmd/httpservice/api/mocks"
)

func setupAPITest(t *testing.T, httpApi *api.API) (*httptest.Server, func()) {
	t.Helper()

	router := chi.NewRouter()
	httpApi.Routes(router)
	server := httptest.NewServer(router)

	return server, server.Close
}

func TestRequestTokenGenerationHandler_Success(t *testing.T) {
	tests := []struct {
		name           string
		projectID      string
		setupUseCase   func(*testing.T) api.RequestTokenGenerationUseCase
		expectedStatus int
	}{
		{
			name:      "Successful Request",
			projectID: "test-project-id",
			setupUseCase: func(t *testing.T) api.RequestTokenGenerationUseCase {
				return &mocks.RequestTokenGenerationUseCaseMock{
					RequestTokenGenerationFunc: func(ctx context.Context, projectID string) error {
						return nil
					},
				}
			},
			expectedStatus: http.StatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpApi := api.New(tt.setupUseCase(t))

			server, tearDownFn := setupAPITest(t, httpApi)
			defer tearDownFn()

			// Prepare request body
			requestBody := api.RequestTokenGenerationInput{
				ProjectID: tt.projectID,
			}
			bodyBytes, err := json.Marshal(requestBody)
			assert.NoError(t, err)

			// Create HTTP request
			req, err := http.NewRequest(http.MethodPost, server.URL+"/generate_token", bytes.NewReader(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check response status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestRequestTokenGenerationHandler_Failure(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		setupUseCase   func(*testing.T) api.RequestTokenGenerationUseCase
		expectedStatus int
	}{
		{
			name:        "Empty Project ID",
			requestBody: `{"project_id": ""}`,
			setupUseCase: func(t *testing.T) api.RequestTokenGenerationUseCase {
				return &mocks.RequestTokenGenerationUseCaseMock{}
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:        "Non valid JSON",
			requestBody: `{"project_id": "`,
			setupUseCase: func(t *testing.T) api.RequestTokenGenerationUseCase {
				return &mocks.RequestTokenGenerationUseCaseMock{}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "UseCase Error",
			requestBody: `{"project_id": "error-project-id"}`,
			setupUseCase: func(t *testing.T) api.RequestTokenGenerationUseCase {
				return &mocks.RequestTokenGenerationUseCaseMock{
					RequestTokenGenerationFunc: func(ctx context.Context, projectID string) error {
						return errors.New("mocked error from use case")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock use case
			mockUseCase := tt.setupUseCase(t)
			httpAPI := api.New(mockUseCase)

			// Setup HTTP server
			server, tearDownFn := setupAPITest(t, httpAPI)
			defer tearDownFn()

			// Create HTTP request
			req, err := http.NewRequest(http.MethodPost, server.URL+"/generate_token", bytes.NewReader([]byte(tt.requestBody)))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Perform the request
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check response status code
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
