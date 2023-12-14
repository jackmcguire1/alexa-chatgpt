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

func (api *GeminiApiClient) GoogleChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(api.Token))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("")
	cs := model.StartChat()

	return cs.SendMessage(ctx, genai.Text(prompt))
}
