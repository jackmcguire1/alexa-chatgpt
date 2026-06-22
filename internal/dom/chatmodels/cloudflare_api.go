package chatmodels

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	localOtel "github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// CloudflareApiClient calls Cloudflare Workers AI using the OpenAI-compatible
// Chat Completions endpoint for text and a direct REST call for image generation.
type CloudflareApiClient struct {
	accountID  string
	apiKey     string
	chatClient openai.Client
	httpClient *http.Client
}

func NewCloudflareApiClient(accountID, apiKey string) *CloudflareApiClient {
	instrumented := &http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithSpanNameFormatter(localOtel.DefaultTransportFormatter),
		),
	}

	baseURL := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/v1", accountID)
	chatClient := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(instrumented),
	)

	return &CloudflareApiClient{
		accountID:  accountID,
		apiKey:     apiKey,
		chatClient: chatClient,
		httpClient: instrumented,
	}
}

func (api *CloudflareApiClient) GenerateContent(ctx context.Context, messages []Message, opts GenerateOptions) (*GenerateResponse, error) {
	var params openai.ChatCompletionNewParams
	params.Model = opts.Model

	for _, msg := range messages {
		switch msg.Role {
		case RoleSystem:
			params.Messages = append(params.Messages, openai.SystemMessage(msg.Content))
		case RoleUser:
			params.Messages = append(params.Messages, openai.UserMessage(msg.Content))
		case RoleAssistant:
			params.Messages = append(params.Messages, openai.AssistantMessage(msg.Content))
		}
	}

	if opts.MaxTokens > 0 {
		params.MaxTokens = openai.Int(int64(opts.MaxTokens))
	}

	resp, err := api.chatClient.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("cloudflare chat error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("cloudflare chat: no choices in response")
	}

	return &GenerateResponse{Content: resp.Choices[0].Message.Content}, nil
}

// GenerateImage calls the Cloudflare Workers AI image generation endpoint.
// The REST API returns JSON: {"result":{"image":"<base64_jpeg>"},"success":true,...}
func (api *CloudflareApiClient) GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/ai/run/%s", api.accountID, model)

	body, err := json.Marshal(map[string]string{"prompt": prompt})
	if err != nil {
		return nil, fmt.Errorf("cloudflare image: failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("cloudflare image: failed to build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+api.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloudflare image: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cloudflare image: failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cloudflare image: HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp struct {
		Result struct {
			Image string `json:"image"`
		} `json:"result"`
		Success bool     `json:"success"`
		Errors  []string `json:"errors"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("cloudflare image: failed to unmarshal response: %w", err)
	}
	if !apiResp.Success || apiResp.Result.Image == "" {
		return nil, fmt.Errorf("cloudflare image: no image in response (errors: %v)", apiResp.Errors)
	}

	imgBytes, err := base64.StdEncoding.DecodeString(apiResp.Result.Image)
	if err != nil {
		return nil, fmt.Errorf("cloudflare image: failed to decode base64 image: %w", err)
	}
	return imgBytes, nil
}
