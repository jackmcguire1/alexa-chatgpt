package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/bucket"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

type SqsHandler struct {
	GenerationModelSvc chatmodels.Service
	ResponseQueue      queue.PullPoll
	Logger             *slog.Logger
	Bucket             bucket.FilePersistance
}

func (handler *SqsHandler) ProcessGenerationRequest(ctx context.Context, req *chatmodels.Request) error {
	handler.Logger.With("payload", utils.ToJSON(req)).Info("invoked with payload")
	execTime := time.Now().UTC()

	var errorMsg string
	var response string
	var imagesResponse []string
	var err error

	if req.ImageModel != nil {
		switch *req.ImageModel {
		case chatmodels.IMAGE_MODEL_STABLE_DIFFUSION,
			chatmodels.IMAGE_MODEL_DALL_E_2,
			chatmodels.IMAGE_MODEL_DALL_E_3,
			chatmodels.IMAGE_MODEL_GEMINI:
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

	switch req.Model {
	case chatmodels.CHAT_MODEL_TRANSLATIONS:
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

func main() {
	jsonH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(jsonH)

	h := &SqsHandler{
		GenerationModelSvc: chatmodels.NewClient(&chatmodels.Resources{
			GPTApi:              chatmodels.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
			GeminiAPI:           chatmodels.NewGeminiApiClient(os.Getenv("GEMINI_API_KEY")),
			CloudflareApiClient: chatmodels.NewCloudflareApiClient(os.Getenv("CLOUDFLARE_ACCOUNT_ID"), os.Getenv("CLOUDFLARE_API_KEY")),
		}),
		ResponseQueue: queue.NewQueue(os.Getenv("RESPONSES_QUEUE_URI")),
		Logger:        logger,
		Bucket: &bucket.Bucket{
			Name: os.Getenv("S3_BUCKET"),
		},
	}
	lambda.Start(h.ProcessSQS)
}
