package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
)

type SqsHandler struct {
	ChatGptSvc    chatgpt.Service
	ResponseQueue queue.PullPoll
}

func (handler *SqsHandler) ProcessChatGPTRequest(ctx context.Context, req *chatgpt.Request) error {
	execTime := time.Now().UTC()

	response, err := handler.ChatGptSvc.AutoComplete(ctx, req.Prompt)
	if err != nil {
		log.Println("failed to process chatgpt request", req.Prompt, err)
		return err
	}

	since := time.Since(execTime)
	log.Printf("pushing response to queue %q %s", response, since)

	err = handler.ResponseQueue.PushMessage(ctx, &chatgpt.LastResponse{
		Prompt:   req.Prompt,
		Response: response,
		TimeDiff: since.String(),
	})
	if err != nil {
		log.Println("failed to publish message to queue", err)
	}

	return err
}

func (handler *SqsHandler) ProcessSQS(ctx context.Context, event events.SQSEvent) error {
	rawData := event.Records[0].Body

	var request *chatgpt.Request

	err := json.Unmarshal([]byte(rawData), &request)
	if err != nil {
		log.Println("failed to unmarshal event", rawData, err)
		return err
	}

	return handler.ProcessChatGPTRequest(ctx, request)
}

func main() {
	h := &SqsHandler{
		ChatGptSvc: chatgpt.NewClient(&chatgpt.Resources{
			Api: chatgpt.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
		}),
		ResponseQueue: queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
	}
	lambda.Start(h.ProcessSQS)
}
