package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/api"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
)

func main() {
	jsonLogH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonLogH)
	svc := chatgpt.NewClient(&chatgpt.Resources{
		Api: chatgpt.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
	})

	pollDelay, _ := strconv.Atoi(os.Getenv("POLL_DELAY"))

	h := api.Handler{
		ChatGptService: svc,
		RequestsQueue:  queue.NewQueue(os.Getenv("REQUESTS_QUEUE_URI")),
		ResponsesQueue: queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		PollDelay:      pollDelay,
		Logger:         logger,
	}
	lambda.Start(h.Invoke)
}
