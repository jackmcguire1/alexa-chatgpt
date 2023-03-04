package chatgpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTApiClient struct {
	Token        string
	OpenAIClient *openai.Client
}

func NewChatGptClient(token string) *ChatGPTApiClient {
	openAIClient := openai.NewClient(token)
	return &ChatGPTApiClient{
		Token:        token,
		OpenAIClient: openAIClient,
	}
}

func (api *ChatGPTApiClient) AutoComplete(ctx context.Context, prompt string) (openai.ChatCompletionResponse, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	return api.OpenAIClient.CreateChatCompletion(ctx, req)
}
