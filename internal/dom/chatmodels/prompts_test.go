package chatmodels

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTextGenerationWithBedrock(t *testing.T) {
	mockResponse := "hello world"
	mockBedrock := &mockBedrockAPI{}
	mockBedrock.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).
		Return(&GenerateResponse{Content: mockResponse}, nil)

	c := Client{&Resources{BedrockAPI: mockBedrock}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_SONNET)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTextGenerationWithSystem(t *testing.T) {
	mockResponse := "hello world"
	mockBedrock := &mockBedrockAPI{}
	mockBedrock.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).
		Return(&GenerateResponse{Content: mockResponse}, nil)

	c := Client{&Resources{BedrockAPI: mockBedrock}}
	resp, err := c.TextGenerationWithSystem(context.Background(), "you are helpful", "steve", CHAT_MODEL_OPUS)
	assert.NoError(t, err)
	assert.EqualValues(t, mockResponse, resp)
}

func TestTextGenerationWithMissingContent(t *testing.T) {
	mockBedrock := &mockBedrockAPI{}
	mockBedrock.On("GenerateContent", mock.Anything, mock.Anything, mock.Anything).
		Return((*GenerateResponse)(nil), MissingContentError)

	c := Client{&Resources{BedrockAPI: mockBedrock}}
	_, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_SONNET)
	assert.Error(t, err)
}

func TestTextGenerationNilBedrock(t *testing.T) {
	c := Client{&Resources{}}
	_, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_SONNET)
	assert.Error(t, err)
}
