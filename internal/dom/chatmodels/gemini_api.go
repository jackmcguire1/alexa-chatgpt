package chatmodels

import (
	"context"
	"encoding/base64"
	"log/slog"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type GeminiApiClient struct {
	credentials *google.Credentials
	Token       string
}

func NewGeminiApiClient(token string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(token)
	creds, err := google.CredentialsFromJSON(context.Background(), tkn, "https://www.googleapis.com/auth/generative-language")
	if err != nil {
		slog.With("error", err).Error("failed to init google creds")
	}
	authToken, err := google.JWTAccessTokenSourceWithScope([]byte(tkn), "https://www.googleapis.com/auth/generative-language")
	authTokenObj, err := authToken.Token()
	if err != nil && authTokenObj != nil {
		slog.With("error", err).Error("failed to get google creds access token")
	}

	return &GeminiApiClient{credentials: creds, Token: authTokenObj.AccessToken}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(api.Token))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	return model.GenerateContent(ctx, genai.Text(prompt))
}
