package chatgpt

import (
	"context"
	"testing"

	gogpt "github.com/sashabaranov/go-gpt3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAutoComplete(t *testing.T) {
	api := &mockAPI{}
	mockResponse := gogpt.CompletionResponse{
		ID:      "",
		Object:  "",
		Created: 0,
		Model:   "",
		Choices: []gogpt.CompletionChoice{
			gogpt.CompletionChoice{
				Text: "is the best",
			},
		},
		Usage: gogpt.Usage{},
	}
	api.On("AutoComplete", mock.Anything, "steve").Return(mockResponse, nil)
	c := Client{&Resources{Api: api}}
	resp, err := c.AutoComplete(context.Background(), "steve")
	assert.NoError(t, err)
	assert.EqualValues(t, resp, "is the best")
}
