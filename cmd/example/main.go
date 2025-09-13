package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	// Initialize resources with optional clients
	resources := &chatmodels.Resources{}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey != "" {
		resources.GPTApi = chatmodels.NewOpenAiApiClient(openAIKey)
		logger.Info("OpenAI client initialized")
	}

	cloudflareAccountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")
	if cloudflareAccountID != "" && cloudflareAPIKey != "" {
		resources.CloudflareApiClient = chatmodels.NewCloudflareApiClient(cloudflareAccountID, cloudflareAPIKey)
		logger.Info("Cloudflare client initialized")
	}

	// Register available clients
	chatmodels.RegisterAvailableClients(
		resources.GPTApi != nil,
		false, // Gemini not used in this example
		false, // Anthropic not used in this example
		resources.CloudflareApiClient != nil,
	)

	logger.With("available_models", chatmodels.AvailableModels).Info("Available models")

	svc := chatmodels.NewClient(resources)

	// Use META model only if Cloudflare is available
	if chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_META) {
		resp, err := svc.TextGeneration(context.Background(), "monkey riding a unicorn", chatmodels.CHAT_MODEL_META)
		if err != nil {
			logger.With("error", err).Error("failed to generate text")
			panic(err)
		}
		logger.With("text-response", resp).Info("got response from text generation model")
	} else {
		logger.Warn("META model not available - Cloudflare client not configured")
	}
}
