package chatmodels

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTextGenerationWithGPT(t *testing.T) {
	api := &mockAPI{}
	mockResponse := "hello world"
	api.On("GenerateTextWithModel", mock.Anything, "steve", mock.Anything).Return(mockResponse, nil)
	c := Client{&Resources{GPTApi: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GPT)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTestTextGenerationWithGemini(t *testing.T) {
	api := &mockAPI{}
	mockResponse := "hello world"
	api.On("GenerateText", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{GeminiAPI: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GEMINI)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTextGenerationWithMetaModel(t *testing.T) {
	api := &mockAPI{}
	mockResponse := "hello world"
	api.On("GenerateTextWithModel", mock.Anything, "steve", mock.Anything).Return(mockResponse, nil)
	c := Client{&Resources{CloudflareApiClient: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_META)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTextGenerationWithMissingContent(t *testing.T) {
	api := &mockAPI{}
	mockError := fmt.Errorf("missing choice %w", MissingContentError)
	api.On("GenerateTextWithModel", mock.Anything, "steve", mock.Anything).Return("", mockError)
	c := Client{&Resources{GPTApi: api}}
	_, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GPT)
	assert.Error(t, err)
}
