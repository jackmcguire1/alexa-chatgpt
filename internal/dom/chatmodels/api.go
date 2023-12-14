package chatmodels

import (
	"context"

	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"
)

type GptAPI interface {
	AutoComplete(context.Context, string) (openai.ChatCompletionResponse, error)
}

type GeminiAPI interface {
	GeminiChat(context.Context, string) (*genai.GenerateContentResponse, error)
}

type mockAPI struct {
	mock.Mock
}

func (api *mockAPI) AutoComplete(ctx context.Context, prompt string) (res openai.ChatCompletionResponse, err error) {
	args := api.Called(ctx, prompt)
	res = args.Get(0).(openai.ChatCompletionResponse)
	return res, args.Error(1)
}

func (api *mockAPI) GeminiChat(ctx context.Context, prompt string) (res *genai.GenerateContentResponse, err error) {
	args := api.Called(ctx, prompt)
	res = args.Get(0).(*genai.GenerateContentResponse)
	return res, args.Error(1)
}
