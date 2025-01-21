package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/api"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
)

func main() {
	jsonLogH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonLogH)
	svc := chatmodels.NewClient(&chatmodels.Resources{
		GPTApi:              chatmodels.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
		GeminiAPI:           chatmodels.NewGeminiApiClient(os.Getenv("GEMINI_API_KEY")),
		CloudflareApiClient: chatmodels.NewCloudflareApiClient(os.Getenv("CLOUDFLARE_ACCOUNT_ID"), os.Getenv("CLOUDFLARE_API_KEY")),
	})

	pollDelay, _ := strconv.Atoi(os.Getenv("POLL_DELAY"))

	h := api.Handler{
		UserCache:       &api.UserCache{Data: make(map[string]*chatmodels.LastResponse)},
		ChatGptService:  svc,
		RequestsQueue:   queue.NewQueue(os.Getenv("REQUESTS_QUEUE_URI")),
		ResponsesQueue:  queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		PollDelay:       pollDelay,
		Logger:          logger,
		Model:           chatmodels.CHAT_MODEL_GPT,
		ImageModel:      chatmodels.IMAGE_MODEL_STABLE_DIFFUSION,
		RandomNumberSvc: api.NewRandomNumberGame(100),
		BattleShips:     api.NewBattleShipSetup(),
	}
	lambda.Start(h.Invoke)
}
