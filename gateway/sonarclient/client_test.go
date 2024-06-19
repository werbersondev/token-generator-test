package sonarclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name           string
		params         TokenGenerationParams
		responseStatus int
		responseBody   string
		expectedToken  string
		expectError    bool
	}{
		{
			name: "successful token generation",
			params: TokenGenerationParams{
				Name:       "test-token",
				ProjectKey: "project-id",
				Type:       ProjectAnalysisTokenType,
			},
			responseStatus: http.StatusOK,
			responseBody:   `{"token": "generated-token"}`,
			expectedToken:  "generated-token",
			expectError:    false,
		},
		{
			name: "failed token generation",
			params: TokenGenerationParams{
				Name:       "test-token",
				ProjectKey: "project-id",
				Type:       ProjectAnalysisTokenType,
			},
			responseStatus: http.StatusBadRequest,
			responseBody:   `{"errors": [{"msg": "error message"}]}`,
			expectedToken:  "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify the Authorization header
				assert.Equal(t, "Bearer dummy-token", r.Header.Get("Authorization"))

				// Verify the Content-Type header
				assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

				// Verify the form data in the request body
				bodyBytes, err := io.ReadAll(r.Body)
				assert.NoError(t, err)
				bodyStr := string(bodyBytes)

				expectedValues := url.Values{
					"name":       {tt.params.Name},
					"projectKey": {tt.params.ProjectKey},
					"type":       {tt.params.Type},
				}

				for key, values := range expectedValues {
					for _, value := range values {
						assert.Contains(t, bodyStr, key+"="+value)
					}
				}

				w.WriteHeader(tt.responseStatus)
				_, _ = w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := New(Config{
				BaseURL:   server.URL,
				AuthToken: "dummy-token",
				Timeout:   5 * time.Second,
			})

			ctx := context.Background()
			token, err := client.GenerateToken(ctx, tt.params)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, token)
			}
		})
	}
}

func TestGenerateProjectAnalysisToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the Authorization header
		assert.Equal(t, "Bearer dummy-token", r.Header.Get("Authorization"))

		// Verify the Content-Type header
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		// Verify the form data in the request body
		expectedBody := url.Values{
			"name":       {"test-token"},
			"projectKey": {"project-id"},
			"type":       {ProjectAnalysisTokenType},
		}.Encode()
		bodyBytes, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		assert.Contains(t, string(bodyBytes), expectedBody)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"token": "generated-token"}`))
	}))
	defer server.Close()

	client := New(Config{
		BaseURL:   server.URL,
		AuthToken: "dummy-token",
		Timeout:   5 * time.Second,
	})

	ctx := context.Background()
	token, err := client.GenerateProjectAnalysisToken(ctx, "project-id", "test-token")

	assert.NoError(t, err)
	assert.Equal(t, "generated-token", token)
}
