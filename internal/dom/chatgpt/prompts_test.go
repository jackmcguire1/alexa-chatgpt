package chatgpt

import (
	"context"
	"testing"

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
	c := Client{&Resources{Api: api}}
	resp, err := c.AutoComplete(context.Background(), "steve")
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
	c := Client{&Resources{Api: api}}
	_, err := c.AutoComplete(context.Background(), "steve")
	assert.Error(t, err)
	assert.Equal(t, "missing choices", err.Error())
}
