package api

import (
	"context"
	"testing"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLaunchIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}

	mockChatGptService.On("AutoComplete", mock.Anything, mock.Anything).Return("chimney", nil)
	h := &Handler{
		mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Type: alexa.LaunchRequestType,
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.EqualValues(t, "Hi, lets begin our convesation!", resp.Body.OutputSpeech.Text)
}

func TestFallbackIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}

	mockChatGptService.On("AutoComplete", mock.Anything, mock.Anything).Return("chimney", nil)
	h := &Handler{
		mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: alexa.FallbackIntent,
				Slots: map[string]alexa.Slot{
					"prompt": {
						Name:        "prompt",
						Value:       "the boy fell down the",
						Resolutions: alexa.Resolutions{},
					},
				},
			},
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, resp.Body.OutputSpeech.Text, "chimney")
}

func TestAutoCompleteIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}

	mockChatGptService.On("AutoComplete", mock.Anything, mock.Anything).Return("chimney", nil)
	h := &Handler{
		mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: "AutoCompleteIntent",
				Slots: map[string]alexa.Slot{
					"prompt": {
						Name:        "prompt",
						Value:       "the boy fell down the",
						Resolutions: alexa.Resolutions{},
					},
				},
			},
			Type: alexa.IntentRequestType,
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.EqualValues(t, "the boy fell down the chimney", resp.Body.OutputSpeech.Text)
}

func TestStopIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: alexa.StopIntent,
			},
			Type: alexa.IntentRequestType,
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.EqualValues(t, resp.Body.OutputSpeech.Text, "Good bye")
	assert.True(t, resp.Body.ShouldEndSession)
}

func TestCancelIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: alexa.CancelIntent,
			},
			Type: alexa.IntentRequestType,
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.EqualValues(t, resp.Body.OutputSpeech.Text, "okay, i'm listening")
	assert.False(t, resp.Body.ShouldEndSession)
}

func TestHelpIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: alexa.HelpIntent,
			},
			Type: alexa.IntentRequestType,
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.EqualValues(t, resp.Body.OutputSpeech.Text, "Simply repeat, complete the sentence followed by a desired sentence")
	assert.False(t, resp.Body.ShouldEndSession)
}

func TestUnsupportedIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: "AMAZON.random",
			},
			Type: "AMAZON.random",
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.EqualValues(t, resp.Body.OutputSpeech.Text, "unsupported intent!")
	assert.False(t, resp.Body.ShouldEndSession)
}
