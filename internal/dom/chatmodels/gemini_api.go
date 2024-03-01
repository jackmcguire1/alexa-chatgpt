package chatmodels

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"encoding/base64"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log/slog"
)

type GeminiApiClient struct {
	credentials *google.Credentials
	projectID   string
}

func NewGeminiApiClient(token, projectID string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(token)
	creds, err := google.CredentialsFromJSON(context.Background(), tkn, "https://www.googleapis.com/auth/generative-language")
	if err != nil {
		slog.With("error", err).Error("failed to init google creds")
	}

	return &GeminiApiClient{credentials: creds, projectID: projectID}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, api.projectID, "us-central1", option.WithCredentials(api.credentials))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.0-pro")
	return model.GenerateContent(ctx, genai.Text(prompt))
}
