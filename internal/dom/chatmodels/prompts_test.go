package chatmodels

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tmc/langchaingo/llms"
)

func TestTextGenerationWithGPT(t *testing.T) {
	mockResponse := "hello world"

	api := &mockAPI{}
	mockLlm := &mockLlmModel{}
	mockLlm.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(&llms.ContentResponse{
		Choices: []*llms.ContentChoice{{Content: mockResponse}},
	}, nil)

	api.On("GenerateTextWithModel", mock.Anything, "steve", mock.Anything).Return(mockResponse, nil)
	api.On("GetModel").Return(mockLlm, nil)
	c := Client{&Resources{GPTApi: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GPT)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTestTextGenerationWithGemini(t *testing.T) {
	mockResponse := "hello world"
	mockLlm := &mockLlmModel{}
	mockLlm.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(&llms.ContentResponse{
		Choices: []*llms.ContentChoice{{Content: mockResponse}},
	}, nil)

	api := &mockAPI{}
	api.On("GenerateTextWithModel", mock.Anything, "steve", mock.Anything).Return(mockResponse, nil)
	api.On("GetModel").Return(mockLlm, nil)

	c := Client{&Resources{GeminiAPI: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GEMINI)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTextGenerationWithMetaModel(t *testing.T) {
	mockResponse := "hello world"
	mockLlm := &mockLlmModel{}
	mockLlm.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(&llms.ContentResponse{
		Choices: []*llms.ContentChoice{{Content: mockResponse}},
	}, nil)

	api := &mockAPI{}
	api.On("GetModel").Return(mockLlm, nil)

	c := Client{&Resources{CloudflareApiClient: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_META)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTextGenerationWithMissingContent(t *testing.T) {
	api := &mockAPI{}
	mockLlm := &mockLlmModel{}

	mockError := fmt.Errorf("missing choice %w", MissingContentError)
	mockLlm.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).Return(&llms.ContentResponse{
		Choices: nil,
	}, mockError)

	api.On("GetModel").Return(mockLlm, nil)

	c := Client{&Resources{GPTApi: api}}
	_, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GPT)
	assert.Error(t, err)
}
