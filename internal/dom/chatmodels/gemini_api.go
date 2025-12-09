package chatmodels

import (
	"context"
	"encoding/base64"
	"log/slog"
	"net/http"

	"cloud.google.com/go/auth"
	"cloud.google.com/go/auth/httptransport"
	"cloud.google.com/go/auth/oauth2adapt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai/googlegenai"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2/google"
)

const (
	VERTEX_API_LOCATION string = "us-central1"
)

type GeminiApiClient struct {
	apiKey          string
	credentials     *google.Credentials
	VertexLlmClient *googlegenai.GoogleAI
	GeminiLlmClient *googlegenai.GoogleAI
}

func NewGeminiApiClient(credsToken string, apiKey string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(credsToken)

	creds, err := google.CredentialsFromJSON(context.Background(), tkn, "https://www.googleapis.com/auth/generative-language", "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		slog.
			With("error", err).
			Error("failed to init google creds")
		panic(err)
	}

	credentialOptions := &auth.CredentialsOptions{
		TokenProvider: oauth2adapt.TokenProviderFromTokenSource(creds.TokenSource),
		JSON:          creds.JSON,
	}

	httpClient, err := httptransport.NewClient(
		&httptransport.Options{
			Credentials: auth.NewCredentials(credentialOptions),
			Headers: http.Header{
				"X-Goog-User-Project": []string{creds.ProjectID},
			},
		})
	if err != nil {
		panic(err)
	}

	client := http.Client{Transport: otelhttp.NewTransport(httpClient.Transport, otelhttp.WithSpanNameFormatter(otel.DefaultTransportFormatter))}
	vertexClient, err := googlegenai.New(
		context.Background(),
		googlegenai.WithHTTPClient(&client),
		googlegenai.WithCloudProject(creds.ProjectID),
		googlegenai.WithCloudLocation(VERTEX_API_LOCATION),
		googlegenai.WithCredentialsJSON(tkn, []string{"https://www.googleapis.com/auth/generative-language", "https://www.googleapis.com/auth/cloud-platform"}),
		googlegenai.WithAPIBackend(googlegenai.APIVertexBackend),
	)
	if err != nil {
		slog.With("error", err).Error("failed to init vertex client")
		panic(err)
	}

	geminiClient := NewGeminiApiKeyClient(apiKey)
	return &GeminiApiClient{
		apiKey:          apiKey,
		credentials:     creds,
		VertexLlmClient: vertexClient,
		GeminiLlmClient: geminiClient,
	}
}

func NewGeminiApiKeyClient(apiKey string) *googlegenai.GoogleAI {
	geminiClient, err := googlegenai.NewFromApiKey(
		context.Background(),
		apiKey,
		googlegenai.WithAPIBackend(googlegenai.APIGeminiBackend),
	)
	if err != nil {
		slog.With("error", err).Error("failed to init vertex client")
		panic(err)
	}

	return geminiClient
}

func (api *GeminiApiClient) GetModel() llms.Model {
	return api.VertexLlmClient
}

func (api *GeminiApiClient) GenerateContent(
	ctx context.Context,
	messages []llms.MessageContent,
	options ...llms.CallOption,
) (*llms.ContentResponse, error) {
	return api.GeminiLlmClient.GenerateContent(ctx, messages, options...)
}

func (api *GeminiApiClient) GenerateImage(ctx context.Context, prompt string, model string) (res []byte, err error) {
	resp, err := api.VertexLlmClient.GenerateImage(
		ctx,
		prompt,
		llms.WithModel(model),
		llms.WithResponseMIMEType("image/jpeg"),
	)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}
