package chatmodels

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type OpenAIApiClient struct {
	Token        string
	OpenAIClient *openai.Client
}

func NewOpenAiApiClient(token string) *OpenAIApiClient {
	openAIClient := openai.NewClient(token)
	return &OpenAIApiClient{
		Token:        token,
		OpenAIClient: openAIClient,
	}
}

func (api *OpenAIApiClient) AutoComplete(ctx context.Context, prompt string) (openai.ChatCompletionResponse, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4TurboPreview,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	return api.OpenAIClient.CreateChatCompletion(ctx, req)
}
