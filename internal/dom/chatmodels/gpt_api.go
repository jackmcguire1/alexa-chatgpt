package chatmodels

import (
	"context"
	"encoding/base64"
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
	OpenAIClient *openai.Client
	LlmClient    *langchain_openai.LLM
}

func NewOpenAiApiClient(token string) *OpenAIApiClient {
	openAIClient := openai.NewClient(token)

	llm, err := langchain_openai.New(
		langchain_openai.WithToken(token),
	)
	if err != nil {
		log.Fatal(err)
	}

	return &OpenAIApiClient{
		Token:        token,
		LlmClient:    llm,
		OpenAIClient: openAIClient,
	}
}

func (api *OpenAIApiClient) GenerateTextWithModel(ctx context.Context, prompt string, model string) (string, error) {
	return llms.GenerateFromSinglePrompt(ctx, api.LlmClient, prompt, llms.WithModel(model))
}

func (api *OpenAIApiClient) GetModel() llms.Model {
	return api.LlmClient
}

func (api *OpenAIApiClient) GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error) {
	req := openai.ImageRequest{
		Model:          model,
		Prompt:         prompt,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		Quality:        "standard",
		N:              1,
	}

	respBase64, err := api.OpenAIClient.CreateImage(ctx, req)
	if err != nil {
		return nil, err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		return nil, err
	}

	return imgBytes, nil
}
