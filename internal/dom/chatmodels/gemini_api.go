package chatmodels

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"

	"cloud.google.com/go/vertexai/genai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const MODEL string = "gemini-1.0-pro"

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

func (api *GeminiApiClient) GenerateText(ctx context.Context, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, api.credentials.ProjectID, "us-central1", option.WithCredentialsJSON(api.credentials.JSON))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel(MODEL)
	res, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	// extract and return response
	if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprint(res.Candidates[0].Content.Parts[0]), nil
	}
	return "", fmt.Errorf("gemini missing content in response %w", MissingContentError)
}
