package chatmodels

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"

	"golang.org/x/oauth2/google"
	"google.golang.org/genai"
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
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Location:    "us-central1",
		Credentials: api.credentials,
		Backend:     genai.BackendVertexAI,
	})
	if err != nil {
		return "", err
	}

	res, err := client.Models.GenerateContent(ctx, MODEL, genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}

	// extract and return response
	if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprint(res.Candidates[0].Content.Parts[0]), nil
	}
	return "", fmt.Errorf("gemini missing content in response %w", MissingContentError)
}
