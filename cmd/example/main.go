package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	resources := &chatmodels.Resources{
		BedrockAPI: chatmodels.NewBedrockApiClient(),
		MantleAPI:  chatmodels.NewMantleApiClient(),
	}
	logger.Info("Bedrock client initialized")

	logger.With("available_models", chatmodels.AvailableModels).Info("Available models")

	svc := chatmodels.NewClient(resources)

	if chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_SONNET) {
		resp, err := svc.TextGeneration(context.Background(), "monkey riding a unicorn", chatmodels.CHAT_MODEL_SONNET)
		if err != nil {
			logger.With("error", err).Error("failed to generate text")
			panic(err)
		}
		logger.With("text-response", resp).Info("got response from text generation model")
	} else {
		logger.Warn("Sonnet model not available - Bedrock not configured")
	}
}
