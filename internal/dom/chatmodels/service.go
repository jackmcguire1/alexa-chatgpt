package chatmodels

import "context"

type Resources struct {
	GPTApi              GptAPI
	GeminiAPI           GeminiAPI
	CloudflareApiClient CloudFlareAiWorkerAPI
}

type Service interface {
	AutoComplete(context.Context, string, ChatModel) (string, error)
	GenerateImage(context.Context, string, ChatModel) ([]byte, error)
}

type Client struct {
	*Resources
}

func NewClient(r *Resources) *Client {
	return &Client{
		r,
	}
}
