package chatgpt

import (
	"context"

	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/stretchr/testify/mock"
)

type API interface {
	AutoComplete(context.Context, string) (gogpt.CompletionResponse, error)
}

type mockAPI struct {
	mock.Mock
}

func (api *mockAPI) AutoComplete(ctx context.Context, prompt string) (res gogpt.CompletionResponse, err error) {
	args := api.Called(ctx, prompt)
	res = args.Get(0).(gogpt.CompletionResponse)
	return res, args.Error(1)
}
