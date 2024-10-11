package chatmodels

import (
	"context"
	"encoding/base64"
	"log/slog"

	"cloud.google.com/go/vertexai/genai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const MODEL string = "gemini-1.5-pro-latest"

type GeminiApiClient struct {
	credentials *google.Credentials
}

func NewGeminiApiClient(token string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(token)

	creds, err := google.CredentialsFromJSON(context.Background(), tkn, "https://www.googleapis.com/auth/generative-language")
	if err != nil {
		slog.With("error", err).Error("failed to init google creds")
	}

	return &GeminiApiClient{credentials: creds}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, api.credentials.ProjectID, "us-central1", option.WithCredentialsJSON(api.credentials.JSON))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel(MODEL)
	return model.GenerateContent(ctx, genai.Text(prompt))
}
