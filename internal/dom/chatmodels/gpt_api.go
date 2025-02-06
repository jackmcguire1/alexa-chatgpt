package chatmodels

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/tmc/langchaingo/llms"
	langchain_openai "github.com/tmc/langchaingo/llms/openai"
)

var IMAGE_MODEL_TO_OPENAI_MODEL = map[ImageModel]string{
	IMAGE_MODEL_DALL_E_3: openai.CreateImageModelDallE3,
	IMAGE_MODEL_DALL_E_2: openai.CreateImageModelDallE2,
}

var CHAT_MODEL_TO_OPENAI_MODEL = map[ChatModel]string{
	CHAT_MODEL_GPT:    openai.O1Mini,
	CHAT_MODEL_GPT_V4: openai.GPT4o,
}

type OpenAIApiClient struct {
	Token        string
	OpenAIClient *langchain_openai.LLM
}

func NewOpenAiApiClient(token string) *OpenAIApiClient {

	llm, err := langchain_openai.New(
		langchain_openai.WithToken(token),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &OpenAIApiClient{
		Token:        token,
		OpenAIClient: llm,
	}
}

func (api *OpenAIApiClient) GenerateTextWithModel(ctx context.Context, prompt string, model string) (string, error) {

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Answer as if the user is asking a question"),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	completion, err := api.OpenAIClient.GenerateContent(ctx, content, llms.WithModel(model))
	if err != nil {
		return "", nil
	}
	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no choices in response from openai")
	}

	return completion.Choices[0].Content, nil
}

func (api *OpenAIApiClient) GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "Answer as if the user is asking a question"),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	completion, err := api.OpenAIClient.GenerateContent(ctx, content, llms.WithModel(model))
	if err != nil {
		return nil, err
	}
	if len(completion.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response from openai")
	}

	imgBytes, err := base64.StdEncoding.DecodeString(completion.Choices[0].Content)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil
}
