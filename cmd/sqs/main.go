package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

type SqsHandler struct {
	ChatModelSvc  chatmodels.Service
	ResponseQueue queue.PullPoll
	Logger        *slog.Logger
}

func (handler *SqsHandler) ProcessChatGPTRequest(ctx context.Context, req *chatmodels.Request) error {
	execTime := time.Now().UTC()

	response, err := handler.ChatModelSvc.AutoComplete(ctx, req.Prompt, req.Model)
	if err != nil {
		handler.Logger.
			With("prompt", req.Prompt).
			With("error", err).
			Error("failed to process chat model request")
		return err
	}

	since := time.Since(execTime)

	handler.Logger.
		With("response", response).
		With("since", since).
		Info("pushing response to queue")

	event := &chatmodels.LastResponse{
		Prompt:   req.Prompt,
		Response: response,
		TimeDiff: since.String(),
		Model:    req.Model,
	}
	err = handler.ResponseQueue.PushMessage(ctx, event)
	if err != nil {
		handler.Logger.
			With("event", utils.ToJSON(event)).
			With("error", err).
			Error("failed to publish message to queue")
	}

	return err
}

func (handler *SqsHandler) ProcessSQS(ctx context.Context, event events.SQSEvent) error {
	rawData := event.Records[0].Body

	var request *chatmodels.Request

	err := json.Unmarshal([]byte(rawData), &request)
	if err != nil {
		handler.Logger.
			With("data", string(rawData)).
			With("error", err).
			Error("failed to unmarshal event")

		return err
	}

	return handler.ProcessChatGPTRequest(ctx, request)
}

func main() {
	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonH)

	h := &SqsHandler{
		ChatModelSvc: chatmodels.NewClient(&chatmodels.Resources{
			GPTApi:    chatmodels.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
			GeminiAPI: chatmodels.NewGeminiApiClient(os.Getenv("GEMINI_API_KEY")),
		}),
		ResponseQueue: queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		Logger:        logger,
	}
	lambda.Start(h.ProcessSQS)
}
