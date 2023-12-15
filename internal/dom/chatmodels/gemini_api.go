package chatmodels

import (
	"context"
	"encoding/base64"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiApiClient struct {
	Token []byte
}

func NewGeminiApiClient(token string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(token)
	return &GeminiApiClient{Token: tkn}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, option.WithCredentialsJSON(api.Token))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	return model.GenerateContent(ctx, genai.Text(prompt))
}
