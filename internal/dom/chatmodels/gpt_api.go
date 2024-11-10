package chatmodels

import (
	"context"
	"encoding/base64"

	"github.com/sashabaranov/go-openai"
)

const (
	OPENAI_IMAGE_MODEL_DALL_E_3 string = "dall-e-3"
	OPENAI_IMAGE_MODEL_DALL_E_2 string = "dall-e-2"
)

var IMAGE_MODEL_TO_OPENAI_MODEL = map[ImageModel]string{
	IMAGE_MODEL_DALL_E_3: OPENAI_IMAGE_MODEL_DALL_E_3,
	IMAGE_MODEL_DALL_E_2: OPENAI_IMAGE_MODEL_DALL_E_2,
}

type OpenAIApiClient struct {
	Token        string
	OpenAIClient *openai.Client
}

func NewOpenAiApiClient(token string) *OpenAIApiClient {
	openAIClient := openai.NewClient(token)
	return &OpenAIApiClient{
		Token:        token,
		OpenAIClient: openAIClient,
	}
}

func (api *OpenAIApiClient) AutoComplete(ctx context.Context, prompt string) (openai.ChatCompletionResponse, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	return api.OpenAIClient.CreateChatCompletion(ctx, req)
}

func (api *OpenAIApiClient) GenerateImage(ctx context.Context, prompt string, model string) ([]byte, error) {
	req := openai.ImageRequest{
		Model:          model,
		Prompt:         prompt,
		Size:           openai.CreateImageSize256x256,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
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
