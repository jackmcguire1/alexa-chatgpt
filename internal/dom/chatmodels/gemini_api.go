package chatmodels

import (
	"context"
	"encoding/base64"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type GeminiApiClient struct {
	credentials *google.Credentials
}

func NewGeminiApiClient(token string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(token)
	creds, _ := google.CredentialsFromJSON(context.Background(), tkn, "https://www.googleapis.com/auth/generative-language")

	return &GeminiApiClient{credentials: creds}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, option.WithCredentials(api.credentials))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	return model.GenerateContent(ctx, genai.Text(prompt))
}
