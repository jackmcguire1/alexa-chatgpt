package chatmodels

import (
	"context"
	"encoding/base64"
	"fmt"

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

func (api *OpenAIApiClient) GenerateText(ctx context.Context, prompt string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4o,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}
	resp, err := api.OpenAIClient.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 || resp.Choices[0].Message.Content == "" {
		return "", fmt.Errorf("missing choices %w", MissingContentError)
	}

	return resp.Choices[0].Message.Content, nil
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
