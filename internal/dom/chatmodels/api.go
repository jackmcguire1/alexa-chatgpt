package chatmodels

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
)

type LlmModel interface {
	GetModel(options ...llms.CallOption) llms.Model
}

type LlmContentGenerator interface {
	GenerateContent(
		ctx context.Context,
		messages []llms.MessageContent,
		options ...llms.CallOption,
	) (*llms.ContentResponse, error)
}

type GptAPI interface {
	LlmModel
	LlmContentGenerator
	GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error)
}

type GeminiAPI interface {
	LlmModel
	LlmContentGenerator
	GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error)
}

type AnthropicAPI interface {
	LlmModel
	LlmContentGenerator
}

type CloudFlareAiWorkerAPI interface {
	LlmModel
	LlmContentGenerator
	GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error)
	GenerateTranslation(ctx context.Context, req *GenerateTranslationRequest) (string, error)
}

type mockLlmModel struct {
	llms.Model
	mock.Mock
}

func (api *mockLlmModel) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (content *llms.ContentResponse, err error) {
	args := api.Called(ctx, messages, options)
	if args.Get(0) != nil {
		content = args.Get(0).(*llms.ContentResponse)
	}
	return content, args.Error(1)
}

type mockAPI struct {
	mock.Mock
}

func (api *mockAPI) GenerateContent(ctx context.Context, messages []llms.MessageContent, options ...llms.CallOption) (content *llms.ContentResponse, err error) {
	args := api.Called(ctx, messages, options)
	if args.Get(0) != nil {
		content = args.Get(0).(*llms.ContentResponse)
	}
	return content, args.Error(1)
}

func (api *mockAPI) GenerateText(ctx context.Context, prompt string) (string, error) {
	args := api.Called(ctx, prompt)
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

func (api *mockAPI) GetModel(options ...llms.CallOption) llms.Model {
	args := api.Called(options)
	return args.Get(0).(llms.Model)
}
