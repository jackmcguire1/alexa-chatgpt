package chatgpt

import "context"

const ChatGPTIntent = "ChatGPTIntent"

type Resources struct {
	Api API
}

type Service interface {
	GetPrompt(context.Context, string) (string, error)
}

type Client struct {
	*Resources
}

func NewClient(r *Resources) *Client {
	return &Client{
		r,
	}
}
