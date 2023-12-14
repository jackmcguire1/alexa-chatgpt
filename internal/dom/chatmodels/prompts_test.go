package chatmodels

import (
	"context"
	"testing"

	"github.com/google/generative-ai-go/genai"
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
	api.On("AutoComplete", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{GPTApi: api}}
	resp, err := c.AutoComplete(context.Background(), "steve", CHAT_MODEL_GPT)
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
	api.On("AutoComplete", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{GPTApi: api}}
	_, err := c.AutoComplete(context.Background(), "steve", CHAT_MODEL_GPT)
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
				TokenCount:       0,
			},
		},
		PromptFeedback: nil,
	}
	api.On("GoogleChat", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{GoogleAPI: api}}
	resp, err := c.AutoComplete(context.Background(), "steve", CHAT_MODEL_GOOGLE)
	assert.NoError(t, err)
	assert.EqualValues(t, "is the best", resp)
}
