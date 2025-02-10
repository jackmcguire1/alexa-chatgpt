package chatmodels

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai/googlegenai"
	"golang.org/x/oauth2/google"
)

const (
	VERTEX_MODEL        string = "gemini-2.0-flash-exp"
	VERTEX_API_LOCATION string = "us-central1"
	IMAGE_IMAGEN_MODEL         = "imagen-3.0-generate-002"
)

type GeminiApiClient struct {
	credentials *google.Credentials
	LlmClient   *googlegenai.GoogleAI
}

func NewGeminiApiClient(credsToken string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(credsToken)

	creds, err := google.CredentialsFromJSON(context.Background(), tkn, "https://www.googleapis.com/auth/generative-language", "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		slog.
			With("error", err).
			Error("failed to init google creds")
		panic(err)
	}

	vertexClient, err := googlegenai.New(
		context.Background(),
		googlegenai.WithCloudProject(creds.ProjectID),
		googlegenai.WithCloudLocation(VERTEX_API_LOCATION),
		googlegenai.WithCredentialsJSON(tkn, nil),
		googlegenai.WithApiBackend(googlegenai.VERTEX_BACKEND),
	)
	if err != nil {
		slog.With("error", err).Error("failed to init vertex client")
		panic(err)
	}

	return &GeminiApiClient{
		credentials: creds,
		LlmClient:   vertexClient,
	}
}

func (api *GeminiApiClient) GetModel() llms.Model {
	return api.LlmClient
}

func (api *GeminiApiClient) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return api.LlmClient.GenerateContent(ctx, messages, options...)
}

func (api *GeminiApiClient) GenerateImage(ctx context.Context, prompt string, model string) (res []byte, err error) {
	resp, err := api.LlmClient.GenerateImage(
		ctx,
		llms.TextContent{Text: prompt},
		llms.WithModel(model),
		llms.WithResponseMIMEType("image/jpeg"),
	)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, fmt.Errorf("empty image response")
	}

	return resp[0].Data, nil
}
