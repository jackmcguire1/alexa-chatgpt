package api

import (
	"context"
	"testing"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRandomNumberGameCorrectGuess(t *testing.T) {
	mockChatGptService := &chatmodels.MockClient{}
	mockChatGptService.On("TextGeneration", mock.Anything, mock.Anything, mock.Anything).Return("Congratulations! You guessed it right.", nil)

	h := &Handler{
		ChatGptService:  mockChatGptService,
		Logger:          logger,
		RandomNumberSvc: NewRandomNumberGame(10),
		Model:           chatmodels.CHAT_MODEL_GPT,
	}

	h.RandomNumberSvc.Number = 5

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: "RandomNumberGame",
				Slots: map[string]alexa.Slot{
					"number": {
						Name:  "number",
						Value: "5",
					},
				},
			},
			Type: alexa.IntentRequestType,
		},
		Context: alexa.Context{},
	}

	resp, err := h.RandomNumberGame(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, "Congratulations! You guessed it right.", resp.Body.OutputSpeech.Text)
}
