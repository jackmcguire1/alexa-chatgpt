package main

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	mockChatGptSvc := &chatgpt.MockClient{}
	mockChatGptSvc.On("AutoComplete", mock.Anything, "tell me a random fact").Return("The battle of zanzibar lasted 30 minutes.", nil)

	mockQueue := &queue.MockQueue{}
	mockQueue.On("PushMessage", mock.Anything, mock.Anything).Return(nil)

	h := &SqsHandler{
		ChatGptSvc:    mockChatGptSvc,
		ResponseQueue: mockQueue,
	}

	request := &chatgpt.Request{
		Prompt: "tell me a random fact",
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
