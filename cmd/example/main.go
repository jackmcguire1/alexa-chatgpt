package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	svc := chatmodels.NewClient(&chatmodels.Resources{
		GPTApi:              chatmodels.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
		CloudflareApiClient: chatmodels.NewCloudflareApiClient(os.Getenv("CLOUDFLARE_ACCOUNT_ID"), os.Getenv("CLOUDFLARE_API_KEY")),
	})
	resp, err := svc.TextGeneration(context.Background(), "monkey riding a unicorn", chatmodels.CHAT_MODEL_META)
	if err != nil {
		logger.With("error", err).Error("failed to generate text")
		panic(err)
	}
	logger.With("text-response", resp).Info("got response from text generation model")
}
