package sonarclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

const (
	defaultTimeout = 10 * time.Second

	UserTokenType            = "USER_TOKEN"
	GlobalAnalysisTokenType  = "GLOBAL_ANALYSIS_TOKEN"
	ProjectAnalysisTokenType = "PROJECT_ANALYSIS_TOKEN"
)

type Config struct {
	Timeout   time.Duration
	BaseURL   string
	AuthToken string
}

type HTTPClient struct {
	client    *http.Client
	baseURL   string
	authToken string
}

func New(config Config) *HTTPClient {
	if config.Timeout <= 0 {
		config.Timeout = defaultTimeout
	}

	retryableClient := retryablehttp.NewClient()
	retryableClient.HTTPClient = &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   config.Timeout,
	}

	retryableClient.Logger = hclog.NewNullLogger()

	return &HTTPClient{
		client:    retryableClient.StandardClient(),
		baseURL:   config.BaseURL,
		authToken: config.AuthToken,
	}
}

type TokenGenerationParams struct {
	Name           string
	ExpirationDate string
	Login          string
	ProjectKey     string
	Type           string
}

func (c *HTTPClient) GenerateProjectAnalysisToken(ctx context.Context, projectID, tokenName string) (string, error) {
	return c.GenerateToken(ctx, TokenGenerationParams{
		Name:       tokenName,
		ProjectKey: projectID,
		Type:       ProjectAnalysisTokenType,
	})
}

func (c *HTTPClient) GenerateToken(ctx context.Context, params TokenGenerationParams) (string, error) {
	urlTarget := fmt.Sprintf("%s/api/user_tokens/generate", c.baseURL)

	formData := url.Values{
		"name": {params.Name},
	}

	if params.ExpirationDate != "" {
		formData.Set("expirationDate", params.ExpirationDate)
	}
	if params.Login != "" {
		formData.Set("login", params.Login)
	}
	if params.ProjectKey != "" {
		formData.Set("projectKey", params.ProjectKey)
	}
	if params.Type != "" {
		formData.Set("type", params.Type)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlTarget, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("closing response body")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		dumpResponse, _ := httputil.DumpResponse(resp, true)
		return "", fmt.Errorf("unexpected status code: %d \n dump response: %s ", resp.StatusCode, dumpResponse)
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decoding response body: %w", err)
	}

	return response.Token, nil
}
