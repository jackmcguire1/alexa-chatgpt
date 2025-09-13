package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	otelsetup "github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/bucket"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var tracer = otel.Tracer("prompt-handler")

type SqsHandler struct {
	GenerationModelSvc chatmodels.Service
	ResponseQueue      queue.PullPoll
	Logger             *slog.Logger
	Bucket             bucket.FilePersistance
}

func (handler *SqsHandler) ProcessGenerationRequest(ctx context.Context, req *chatmodels.Request) error {
	defer func() {
		if r := recover(); r != nil {
			handler.Logger.
				With("payload", utils.ToJSON(req)).
				With("error", r).
				Error("panic occurred during request processing")

			handler.Recover(ctx, req)
		}
	}()

	handler.Logger.With("payload", utils.ToJSON(req)).Info("invoked with payload")
	execTime := time.Now().UTC()

	var errorMsg string
	var response string
	var imagesResponse []string
	var err error

	ctx, span := tracer.Start(ctx, "ProcessGenerationRequest")
	span.SetAttributes(
		attribute.String("prompt", req.Prompt),
		attribute.String("system-prompt", req.SystemPrompt),
	)

	defer span.End()

	if req.ImageModel != nil {
		switch *req.ImageModel {
		case chatmodels.IMAGE_MODEL_STABLE_DIFFUSION,
			chatmodels.IMAGE_MODEL_DALL_E_2,
			chatmodels.IMAGE_MODEL_DALL_E_3,
			chatmodels.IMAGE_MODEL_GEMINI:
			span.SetAttributes(
				attribute.String("image-model", string(*req.ImageModel)),
			)

			imageBody, err := handler.GenerationModelSvc.GenerateImage(ctx, req.Prompt, *req.ImageModel)
			if err != nil {
				handler.Logger.
					With("image-model", *req.ImageModel).
					With("prompt", req.Prompt).
					With("error", err).
					Error("failed to generate image from request")

				errorMsg = err.Error()
				goto respond
			}

			imagesResponse, err = handler.processImage(ctx, imageBody)
			if err != nil {
				handler.Logger.
					With("prompt", req.Prompt).
					With("error", err).
					Error("failed to persist image resolutions")

				errorMsg = err.Error()
			}
			goto respond
		}
	}

	span.SetAttributes(attribute.String("model", req.Model.String()))
	switch req.Model {
	case chatmodels.CHAT_MODEL_TRANSLATIONS:
		span.SetAttributes(
			attribute.String("source-language", req.SourceLanguage),
			attribute.String("target-language", req.TargetLanguage),
		)
		response, err = handler.GenerationModelSvc.Translate(ctx, req.Prompt, req.SourceLanguage, req.TargetLanguage, req.Model)
		if err != nil {
			handler.Logger.
				With("prompt", req.Prompt).
				With("error", err).
				Error("failed to process translation request")

			errorMsg = err.Error()
			break
		}
	default:
		if req.SystemPrompt != "" {

			span.SetAttributes(attribute.String("system-prompt", req.SystemPrompt))
			response, err = handler.GenerationModelSvc.TextGenerationWithSystem(ctx, req.SystemPrompt, req.Prompt, req.Model)
		} else {

			response, err = handler.GenerationModelSvc.TextGeneration(ctx, req.Prompt, req.Model)
		}
		if err != nil {
			handler.Logger.
				With("system-prompt", req.SystemPrompt).
				With("prompt", req.Prompt).
				With("error", err).
				Error("failed to process chat model request")

			errorMsg = err.Error()
			break
		}
	}
respond:
	since := time.Since(execTime)

	handler.Logger.
		With("response", response).
		With("since", since).
		Info("pushing response to queue")

	if errorMsg != "" {
		span.RecordError(errors.New(errorMsg))
	}
	span.SetAttributes(
		attribute.Int("response-bytes", len(response)),
		attribute.Int("image-response-count", len(imagesResponse)),
	)

	event := &chatmodels.LastResponse{
		Prompt:         req.Prompt,
		Response:       response,
		TimeDiff:       fmt.Sprintf("%.0f", since.Seconds()),
		Model:          req.Model.String(),
		ImagesResponse: imagesResponse,
		Error:          errorMsg,
		SystemPrompt:   req.SystemPrompt,
	}

	// override the model if image model was set
	if req.ImageModel != nil {
		event.Model = req.ImageModel.String()
	}

	err = handler.ResponseQueue.PushMessage(ctx, event)
	if err != nil {
		span.RecordError(err)
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

	return handler.ProcessGenerationRequest(ctx, request)
}

func (handler *SqsHandler) Recover(ctx context.Context, req *chatmodels.Request) {
	// Push a failure message to the queue
	event := &chatmodels.LastResponse{
		Prompt:       req.Prompt,
		Error:        "an error occured when processing the prompt",
		Model:        req.Model.String(),
		SystemPrompt: req.SystemPrompt,
	}

	if pushErr := handler.ResponseQueue.PushMessage(ctx, event); pushErr != nil {
		handler.Logger.
			With("event", utils.ToJSON(event)).
			With("error", pushErr).
			Error("failed to publish panic message to queue")
	}
}

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

func main() {
	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonH)

	ctx := context.Background()
	tp, err := otelsetup.SetupXrayOtel(ctx)
	if err != nil {
		logger.With("error", err).Error("failed to setup tracer")
		panic(err)
	}
	defer tp.Shutdown(ctx)

	resources := initializeResources()

	h := &SqsHandler{
		GenerationModelSvc: chatmodels.NewClient(resources),
		ResponseQueue:      queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		Logger:             logger,
		Bucket: &bucket.Bucket{
			Name: os.Getenv("S3_BUCKET"),
		},
	}
	lambda.Start(otellambda.InstrumentHandler(h.ProcessSQS, xrayconfig.WithRecommendedOptions(tp)...))
}
