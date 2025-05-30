package chatmodels

import (
	"context"
	"net/http"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const OPUS_MODEL string = "claude-opus-4-20250514"
const SONNET_MODEL string = "claude-sonnet-4-20250514"

var CHAT_MODEL_TO_ANTHROPIC_MODELS = map[ChatModel]string{
	CHAT_MODEL_OPUS:   OPUS_MODEL,
	CHAT_MODEL_SONNET: SONNET_MODEL,
}

type AnthropicApiClient struct {
	Token     string
	Model     string
	LlmClient *anthropic.LLM
}

func NewAnthropicApiClient(token string) *AnthropicApiClient {
	tokenOpt := anthropic.WithToken(token)
	httpClient := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	httpClientOpt := anthropic.WithHTTPClient(httpClient)
	llm, err := anthropic.New(tokenOpt, httpClientOpt)
	if err != nil {
		panic(err)
	}
	return &AnthropicApiClient{
		Token:     token,
		LlmClient: llm,
	}
}

func (api *AnthropicApiClient) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return api.LlmClient.GenerateContent(ctx, messages, options...)
}

func (api *AnthropicApiClient) GetModel() llms.Model {
	return api.LlmClient
}
