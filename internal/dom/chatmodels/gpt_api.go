package chatmodels

import (
	"context"
	"log"
	"net/http"

	"github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"github.com/tmc/langchaingo/llms"
	langchain_openai "github.com/tmc/langchaingo/llms/openai"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var IMAGE_MODEL_TO_OPENAI_MODEL = map[ImageModel]string{
	IMAGE_MODEL_DALL_E_3: "dall-e-3",
	IMAGE_MODEL_DALL_E_2: "dall-e-2",
	IMAGE_MODEL_GPT:      "gpt-image-1",
}

var CHAT_MODEL_TO_OPENAI_MODEL = map[ChatModel]string{
	CHAT_MODEL_GPT:    "gpt-5.1-2025-11-13",
	CHAT_MODEL_GPT_V4: "gpt-4o",
}

type OpenAIApiClient struct {
	Token     string
	LlmClient *langchain_openai.LLM
}

func NewOpenAiApiClient(token string) *OpenAIApiClient {

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport, otelhttp.WithSpanNameFormatter(otel.DefaultTransportFormatter))}

	llm, err := langchain_openai.New(
		langchain_openai.WithHTTPClient(&client),
		langchain_openai.WithToken(token),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &OpenAIApiClient{
		Token:     token,
		LlmClient: llm,
	}
}

func (api *OpenAIApiClient) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return api.LlmClient.GenerateContent(ctx, messages, options...)
}

func (api *OpenAIApiClient) GetModel() llms.Model {
	return api.LlmClient
}

func (api *OpenAIApiClient) GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error) {
	content, err := api.LlmClient.GenerateImage(ctx, prompt, llms.WithModel(model))
	if err != nil {
		return nil, err
	}
	return content.Data, nil
}
