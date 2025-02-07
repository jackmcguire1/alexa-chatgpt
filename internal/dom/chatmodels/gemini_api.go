package chatmodels

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"

	"golang.org/x/oauth2/google"
	"google.golang.org/genai"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/googleai/vertex"
)

const (
	VERTEX_MODEL        string = "gemini-2.0-flash-exp"
	VERTEX_API_LOCATION string = "us-central1"
)

type GeminiApiClient struct {
	credentials *google.Credentials
	GenAIClient *genai.Client
	LlmClient   *vertex.Vertex
}

func NewGeminiApiClient(credsToken string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(credsToken)

	creds, err := google.CredentialsFromJSON(context.Background(), tkn, "https://www.googleapis.com/auth/generative-language", "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		slog.With("error", err).Error("failed to init google creds")
	}

	genAiClient, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		Location:    VERTEX_API_LOCATION,
		Credentials: creds,
		Backend:     genai.BackendVertexAI,
		Project:     creds.ProjectID,
	})
	if err != nil {
		slog.With("error", err).Error("failed to init genAiClient")
		panic(err)
	}

	token, err := creds.TokenSource.Token()
	if err != nil {
		slog.With("error", err).Error("failed to get a google token from credentials")
		panic(err)
	}

	vertexClient, err := vertex.New(context.Background(), googleai.WithCloudProject(creds.ProjectID), googleai.WithAPIKey(token.AccessToken), googleai.WithCloudLocation(VERTEX_API_LOCATION))
	if err != nil {
		slog.With("error", err).Error("failed to init vertex client")
		panic(err)
	}

	return &GeminiApiClient{
		credentials: creds,
		GenAIClient: genAiClient,
		LlmClient:   vertexClient,
	}
}

func (api *GeminiApiClient) GetModel() llms.Model {
	return api.LlmClient
}

func (api *GeminiApiClient) GenerateText(ctx context.Context, prompt string) (string, error) {
	res, err := api.GenAIClient.Models.GenerateContent(ctx, VERTEX_MODEL, genai.Text(prompt), nil)
	if err != nil {
		return "", err
	}

	// extract and return response
	if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
		return fmt.Sprint(res.Candidates[0].Content.Parts[0]), nil
	}
	return "", fmt.Errorf("gemini missing content in response %w", MissingContentError)
}

func (api *GeminiApiClient) GenerateTextWithSystemCommand(ctx context.Context, system string, prompt string) (string, error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}

	resp, err := api.GenerateContent(ctx, content, llms.WithModel(VERTEX_MODEL))
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response from vertex")
	}

	return resp.Choices[0].Content, nil
}

func (api *GeminiApiClient) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return api.LlmClient.GenerateContent(ctx, messages, options...)
}
