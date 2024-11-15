package chatmodels

import (
	"context"

	"cloud.google.com/go/vertexai/genai"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/mock"
)

type GptAPI interface {
	GenerateText(context.Context, string) (string, error)
	GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error)
}

type GeminiAPI interface {
	GenerateText(context.Context, string) (string, error)
}

type CloudFlareAiWorkerAPI interface {
	GenerateTextWithModel(context.Context, string, string) (string, error)
	GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error)
	GenerateTranslation(ctx context.Context, req *GenerateTranslationRequest) (string, error)
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

func (api *mockAPI) GenerateText(ctx context.Context, prompt string) (string, error) {
	args := api.Called(ctx, prompt)
	return args.String(0), args.Error(1)
}

func (api *mockAPI) GenerateTextWithModel(ctx context.Context, prompt string, model string) (string, error) {
	args := api.Called(ctx, prompt, model)
	return args.String(0), args.Error(1)
}

func (api *mockAPI) GenerateImage(ctx context.Context, prompt string, model string) (res []byte, err error) {
	args := api.Called(ctx, prompt, model)
	if args.Get(0) != nil {
		res = args.Get(0).([]byte)
	}
	return res, args.Error(1)
}

func (api *mockAPI) GenerateTranslation(ctx context.Context, req *GenerateTranslationRequest) (string, error) {
	args := api.Called(ctx, req)
	return args.String(0), args.Error(1)
}
