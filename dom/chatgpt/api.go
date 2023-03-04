package chatgpt

import (
	"context"
	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/stretchr/testify/mock"
)

type API interface {
	GetChatPrompt(context.Context, string) (gogpt.CompletionResponse, error)
}

type MockAPI struct {
	mock.Mock
}

func (api *MockAPI) GetChatPrompt(ctx context.Context, prompt string) (res gogpt.CompletionResponse, err error) {
	args := api.Called(ctx, prompt)
	res = args.Get(0).(gogpt.CompletionResponse)
	return res, args.Error(1)
}
