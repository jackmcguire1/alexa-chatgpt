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

	svc := chatmodels.NewClient(resources)

	pollDelay, _ := strconv.Atoi(os.Getenv("POLL_DELAY"))

	// Set default models based on available clients
	defaultChatModel := chatmodels.CHAT_MODEL_GPT
	defaultImageModel := chatmodels.IMAGE_MODEL_STABLE_DIFFUSION

	// Choose first available chat model as default
	if len(chatmodels.AvailableModels) > 0 {
		if resources.GPTApi != nil {
			defaultChatModel = chatmodels.CHAT_MODEL_GPT
		} else if resources.GeminiAPI != nil {
			defaultChatModel = chatmodels.CHAT_MODEL_GEMINI
		} else if resources.AnthropicAPI != nil {
			defaultChatModel = chatmodels.CHAT_MODEL_OPUS
		} else if resources.CloudflareApiClient != nil {
			defaultChatModel = chatmodels.CHAT_MODEL_META
		}
	}

	// Choose first available image model as default
	if len(chatmodels.ImageModels) > 0 {
		if resources.CloudflareApiClient != nil {
			defaultImageModel = chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
		} else if resources.GPTApi != nil {
			defaultImageModel = chatmodels.IMAGE_MODEL_DALL_E_3
		} else if resources.GeminiAPI != nil {
			defaultImageModel = chatmodels.IMAGE_MODEL_GEMINI
		}
	}

	h := api.Handler{
		ChatGptService:  svc,
		RequestsQueue:   queue.NewQueue(os.Getenv("REQUESTS_QUEUE_URI")),
		ResponsesQueue:  queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		PollDelay:       pollDelay,
		Logger:          logger,
		Model:           defaultChatModel,
		ImageModel:      defaultImageModel,
		RandomNumberSvc: api.NewRandomNumberGame(100),
		BattleShips:     api.NewBattleShipSetup(),
	}
	lambda.Start(otellambda.InstrumentHandler(h.Invoke, xrayconfig.WithRecommendedOptions(tracer)...))
}
