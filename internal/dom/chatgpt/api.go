package chatgpt

import (
	"context"

	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"
)

type API interface {
	AutoComplete(context.Context, string) (openai.ChatCompletionResponse, error)
}

type mockAPI struct {
	mock.Mock
}

func (api *mockAPI) AutoComplete(ctx context.Context, prompt string) (res openai.ChatCompletionResponse, err error) {
	args := api.Called(ctx, prompt)
	res = args.Get(0).(openai.ChatCompletionResponse)
	return res, args.Error(1)
}
