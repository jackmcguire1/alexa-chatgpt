package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	mockChatGptSvc := &chatmodels.MockClient{}
	mockChatGptSvc.On("AutoComplete", mock.Anything, "tell me a random fact", chatmodels.CHAT_MODEL_GPT).Return("The battle of zanzibar lasted 30 minutes.", nil)

	mockQueue := &queue.MockQueue{}
	mockQueue.On("PushMessage", mock.Anything, mock.Anything).Return(nil)

	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonH)

	h := &SqsHandler{
		ChatModelSvc:  mockChatGptSvc,
		ResponseQueue: mockQueue,
		Logger:        logger,
	}

	request := &chatmodels.Request{
		Prompt: "tell me a random fact",
		Model:  chatmodels.CHAT_MODEL_GPT,
	}

	err := h.ProcessSQS(context.Background(), events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{
				Body: utils.ToJSON(request),
			},
		},
	})
	assert.NoError(t, err)
}

func TestModelErrorResponse(t *testing.T) {
	mockChatGptSvc := &chatmodels.MockClient{}
	mockChatGptSvc.On("AutoComplete", mock.Anything, "tell me a random fact", chatmodels.CHAT_MODEL_GPT).Return("", fmt.Errorf("there was an issue with your API KEY"))

	mockQueue := &queue.MockQueue{}
	mockQueue.On("PushMessage", mock.Anything, mock.Anything).Return(nil)

	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonH)

	h := &SqsHandler{
		ChatModelSvc:  mockChatGptSvc,
		ResponseQueue: mockQueue,
		Logger:        logger,
	}

	request := &chatmodels.Request{
		Prompt: "tell me a random fact",
		Model:  chatmodels.CHAT_MODEL_GPT,
	}

	err := h.ProcessSQS(context.Background(), events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{
				Body: utils.ToJSON(request),
			},
		},
	})
	assert.NoError(t, err)
	mockQueue.AssertNumberOfCalls(t, "PushMessage", 1)
}

func TestSqsError(t *testing.T) {
	mockChatGptSvc := &chatmodels.MockClient{}
	mockChatGptSvc.On("AutoComplete", mock.Anything, "tell me a random fact", chatmodels.CHAT_MODEL_GPT).Return("", fmt.Errorf("there was an issue with your API KEY"))

	mockQueue := &queue.MockQueue{}
	mockQueue.On("PushMessage", mock.Anything, mock.Anything).Return(fmt.Errorf("500 internal server error"))

	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonH)

	h := &SqsHandler{
		ChatModelSvc:  mockChatGptSvc,
		ResponseQueue: mockQueue,
		Logger:        logger,
	}

	request := &chatmodels.Request{
		Prompt: "tell me a random fact",
		Model:  chatmodels.CHAT_MODEL_GPT,
	}

	err := h.ProcessSQS(context.Background(), events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{
				Body: utils.ToJSON(request),
			},
		},
	})
	assert.Error(t, err)
}

func TestBadJSONPayload(t *testing.T) {
	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonH)

	h := &SqsHandler{
		Logger: logger,
	}

	err := h.ProcessSQS(context.Background(), events.SQSEvent{
		Records: []events.SQSMessage{
			events.SQSMessage{
				Body: `{]`,
			},
		},
	})
	assert.Error(t, err)
}
