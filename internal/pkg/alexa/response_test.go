package alexa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewReponse(t *testing.T) {
	resp := NewResponse("test", "test", true)
	assert.True(t, resp.Body.ShouldEndSession)
	assert.Equal(t, "test", resp.Body.OutputSpeech.Text)
	assert.Equal(t, "test", resp.Body.OutputSpeech.Title)
	assert.Equal(t, "test", resp.Body.Card.Content)
	assert.Equal(t, "test", resp.Body.Card.Title)
}
