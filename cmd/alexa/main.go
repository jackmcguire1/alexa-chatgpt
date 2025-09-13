package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/api"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	otelsetup "github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
)

func initializeResources() *chatmodels.Resources {
	resources := &chatmodels.Resources{}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey != "" {
		resources.GPTApi = chatmodels.NewOpenAiApiClient(openAIKey)
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey != "" {
		resources.GeminiAPI = chatmodels.NewGeminiApiClient(geminiKey)
	}

	cloudflareAccountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")
	if cloudflareAccountID != "" && cloudflareAPIKey != "" {
		resources.CloudflareApiClient = chatmodels.NewCloudflareApiClient(cloudflareAccountID, cloudflareAPIKey)
	}

	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicKey != "" {
		resources.AnthropicAPI = chatmodels.NewAnthropicApiClient(anthropicKey)
	}

	chatmodels.RegisterAvailableClients(
		resources.GPTApi != nil,
		resources.GeminiAPI != nil,
		resources.AnthropicAPI != nil,
		resources.CloudflareApiClient != nil,
	)

	return resources
}

func getDefaultChatModel(resources *chatmodels.Resources) chatmodels.ChatModel {
	if len(chatmodels.AvailableModels) == 0 {
		return chatmodels.CHAT_MODEL_GPT
	}
	
	if resources.GPTApi != nil {
		return chatmodels.CHAT_MODEL_GPT
	}
	if resources.GeminiAPI != nil {
		return chatmodels.CHAT_MODEL_GEMINI
	}
	if resources.AnthropicAPI != nil {
		return chatmodels.CHAT_MODEL_OPUS
	}
	if resources.CloudflareApiClient != nil {
		return chatmodels.CHAT_MODEL_META
	}
	
	return chatmodels.CHAT_MODEL_GPT
}

func getDefaultImageModel(resources *chatmodels.Resources) chatmodels.ImageModel {
	if len(chatmodels.ImageModels) == 0 {
		return chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
	}
	
	if resources.CloudflareApiClient != nil {
		return chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
	}
	if resources.GPTApi != nil {
		return chatmodels.IMAGE_MODEL_DALL_E_3
	}
	if resources.GeminiAPI != nil {
		return chatmodels.IMAGE_MODEL_GEMINI
	}
	
	return chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
}

func main() {
	jsonLogH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonLogH)

	ctx := context.Background()
	tracer, err := otelsetup.SetupXrayOtel(ctx)
	if err != nil {
		logger.With("error", err).Error("failed to setup tracer")
		panic(err)
	}
	defer tracer.Shutdown(ctx)

	resources := initializeResources()
	svc := chatmodels.NewClient(resources)
	pollDelay, _ := strconv.Atoi(os.Getenv("POLL_DELAY"))

	h := api.Handler{
		ChatGptService:  svc,
		RequestsQueue:   queue.NewQueue(os.Getenv("REQUESTS_QUEUE_URI")),
		ResponsesQueue:  queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		PollDelay:       pollDelay,
		Logger:          logger,
		Model:           getDefaultChatModel(resources),
		ImageModel:      getDefaultImageModel(resources),
		RandomNumberSvc: api.NewRandomNumberGame(100),
		BattleShips:     api.NewBattleShipSetup(),
	}
	lambda.Start(otellambda.InstrumentHandler(h.Invoke, xrayconfig.WithRecommendedOptions(tracer)...))
}
