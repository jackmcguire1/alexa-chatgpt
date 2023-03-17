package api

import (
	"context"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
	"testing"
	"time"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLaunchIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}

	h := &Handler{
		ChatGptService: mockChatGptService,
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
	assert.EqualValues(t, "Hi, lets begin our conversation!", resp.Body.OutputSpeech.Text)
}

func TestFallbackIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		ChatGptService: mockChatGptService,
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
	assert.Contains(t, resp.Body.OutputSpeech.Text, "Try again")
}

func TestAutoCompleteIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	mockChatGptService.On("AutoComplete", mock.Anything, mock.Anything).Return("chimney", nil)

	mockRequestsQueue := &queue.MockQueue{}
	mockRequestsQueue.On("PushMessage", mock.Anything, mock.Anything).Return(nil)

	mockResponsesQueue := &queue.MockQueue{}
	queueResponse := chatgpt.LastResponse{Response: "chimney"}
	jsonResp := utils.ToJSON(queueResponse)
	mockResponsesQueue.On("PullMessage", mock.Anything, mock.Anything).Return([]byte(jsonResp), nil)

	h := &Handler{
		ChatGptService: mockChatGptService,
		ResponsesQueue: mockResponsesQueue,
		RequestsQueue:  mockRequestsQueue,
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
	assert.EqualValues(t, "chimney", resp.Body.OutputSpeech.Text)
}

func TestRandomIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}

	mockChatGptService.On("AutoComplete", mock.Anything, mock.Anything).Return("santa fell down the chimney", nil)
	h := &Handler{
		ChatGptService: mockChatGptService,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: alexa.RandomFactIntent,
				Slots: map[string]alexa.Slot{
					"prompt": {
						Name:        "prompt",
						Value:       "tell me a random fact",
						Resolutions: alexa.Resolutions{},
					},
				},
			},
			Type: alexa.RandomFactIntent,
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.EqualValues(t, "santa fell down the chimney", resp.Body.OutputSpeech.Text)
}

func TestLastResponseIntent(t *testing.T) {

	mockResponsesQueue := &queue.MockQueue{}
	queueResponse := chatgpt.LastResponse{Response: "chimney", TimeDiff: time.Since(time.Now().Add(-time.Second)).String()}
	jsonResp := utils.ToJSON(queueResponse)
	mockResponsesQueue.On("PullMessage", mock.Anything, mock.Anything).Return([]byte(jsonResp), nil)

	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		lastResponse:   &queueResponse,
		ChatGptService: mockChatGptService,
		ResponsesQueue: mockResponsesQueue,
	}

	req := alexa.Request{
		Version: "",
		Session: alexa.Session{},
		Body: alexa.ReqBody{
			Intent: alexa.Intent{
				Name: alexa.LastResponseIntent,
				Slots: map[string]alexa.Slot{
					"prompt": {
						Name:        "prompt",
						Value:       "hello",
						Resolutions: alexa.Resolutions{},
					},
				},
			},
			Type: alexa.LastResponseIntent,
		},
		Context: alexa.Context{},
	}

	resp, err := h.Invoke(context.Background(), req)
	assert.NoError(t, err)
	assert.Contains(t, resp.Body.OutputSpeech.Text, "chimney")
}

func TestStopIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		ChatGptService: mockChatGptService,
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
		ChatGptService: mockChatGptService,
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
		ChatGptService: mockChatGptService,
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
	assert.EqualValues(t, resp.Body.OutputSpeech.Text, "Simply repeat, question followed by a desired sentence")
	assert.False(t, resp.Body.ShouldEndSession)
}

func TestUnsupportedIntent(t *testing.T) {
	mockChatGptService := &chatgpt.MockClient{}
	h := &Handler{
		ChatGptService: mockChatGptService,
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
