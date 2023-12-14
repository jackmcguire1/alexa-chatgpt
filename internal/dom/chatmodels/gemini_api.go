package chatmodels

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

type GeminiApiClient struct {
	Token        string
	OpenAIClient *openai.Client
}

func NewGeminiApiClient(token string) *GeminiApiClient {
	return &GeminiApiClient{
		Token: token,
	}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(api.Token), option.WithQuotaProject(""), option.)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-pro")

	return model.GenerateContent(ctx, genai.Text(prompt))
}
