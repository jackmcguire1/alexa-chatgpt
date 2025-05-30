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

	svc := chatmodels.NewClient(&chatmodels.Resources{
		GPTApi:              chatmodels.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
		GeminiAPI:           chatmodels.NewGeminiApiClient(os.Getenv("GEMINI_API_KEY")),
		CloudflareApiClient: chatmodels.NewCloudflareApiClient(os.Getenv("CLOUDFLARE_ACCOUNT_ID"), os.Getenv("CLOUDFLARE_API_KEY")),
		AnthropicAPI:        chatmodels.NewAnthropicApiClient(os.Getenv("ANTHROPIC_API_KEY")),
	})

	pollDelay, _ := strconv.Atoi(os.Getenv("POLL_DELAY"))

	h := api.Handler{
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
	lambda.Start(otellambda.InstrumentHandler(h.Invoke, xrayconfig.WithRecommendedOptions(tracer)...))
}
