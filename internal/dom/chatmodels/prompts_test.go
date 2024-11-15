package chatmodels

import (
	"context"
	"testing"

	"cloud.google.com/go/vertexai/genai"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAutoComplete(t *testing.T) {
	api := &mockAPI{}
	mockResponse := openai.ChatCompletionResponse{
		ID:      "",
		Object:  "",
		Created: 0,
		Model:   "",
		Choices: []openai.ChatCompletionChoice{
			openai.ChatCompletionChoice{
				Message: openai.ChatCompletionMessage{Content: "is the best"},
			},
		},
	}
	api.On("TextGeneration", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{GPTApi: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GPT)
	assert.NoError(t, err)
	assert.EqualValues(t, "is the best", resp)
}

func TestAutoCompleteMissingChoices(t *testing.T) {
	api := &mockAPI{}
	mockResponse := openai.ChatCompletionResponse{
		ID:      "",
		Object:  "",
		Created: 0,
		Model:   "",
		Choices: []openai.ChatCompletionChoice{},
	}
	api.On("TextGeneration", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{GPTApi: api}}
	_, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GPT)
	assert.Error(t, err)
	assert.Equal(t, "missing choices", err.Error())
}

func TestGeminiChat(t *testing.T) {
	api := &mockAPI{}
	part := genai.Text("is the best")
	mockResponse := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			&genai.Candidate{
				Index: 0,
				Content: &genai.Content{
					Parts: []genai.Part{
						part,
					},
					Role: "",
				},
				FinishReason:     0,
				SafetyRatings:    nil,
				CitationMetadata: nil,
			},
		},
		PromptFeedback: nil,
	}
	api.On("GeminiChat", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{GeminiAPI: api}}
	resp, err := c.TextGeneration(context.Background(), "steve", CHAT_MODEL_GEMINI)
	assert.NoError(t, err)
	assert.EqualValues(t, "is the best", resp)
}
