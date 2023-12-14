package chatmodels

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/sashabaranov/go-openai"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type GeminiApiClient struct {
	Token        []byte
	OpenAIClient *openai.Client
}

func NewGeminiApiClient(token string) *GeminiApiClient {
	tkn, _ := base64.StdEncoding.DecodeString(token)

	return &GeminiApiClient{
		Token: tkn,
	}
}

func (api *GeminiApiClient) GeminiChat(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	creds, err := google.CredentialsFromJSON(ctx, api.Token)
	if err != nil {
		log.Printf("got token %s err:%s \n", string(api.Token), err)
		return nil, err
	}

	client, err := genai.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-pro")

	return model.GenerateContent(ctx, genai.Text(prompt))
}
