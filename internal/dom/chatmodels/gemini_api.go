package chatmodels

import (
	"context"
	"encoding/base64"
	"log/slog"

	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type GeminiApiClient struct {
	Token        []byte
	OpenAIClient *openai.Client
}

func NewGeminiApiClient(token string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(token)
	return &GeminiApiClient{
		Token: tkn,
	}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	creds, err := google.CredentialsFromJSON(ctx, api.Token, "https://www.googleapis.com/auth/generative-language")
	client, err := genai.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		slog.
			With("token-json", string(api.Token)).
			Error("failed to process gemini generative request")

		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		slog.
			With("prompt", prompt).
			With("token-json", string(api.Token)).
			Error("failed to process gemini generative req")
	}
	return res, err
}
