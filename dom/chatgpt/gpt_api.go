package chatgpt

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type ChatGPTApiClient struct {
	Token  string
	Client *gogpt.Client
}

func NewChatGptClient(token string) *ChatGPTApiClient {
	c := gogpt.NewClient(token)
	return &ChatGPTApiClient{
		Token:  token,
		Client: c,
	}
}

func (api *ChatGPTApiClient) GetChatPrompt(ctx context.Context, prompt string) (gogpt.CompletionResponse, error) {
	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3Ada,
		MaxTokens: 5,
		Prompt:    prompt,
	}
	return api.Client.CreateCompletion(ctx, req)
}
