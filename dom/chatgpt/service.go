package chatgpt

import "context"

type Resources struct {
	Api API
}

type Service interface {
	AutoComplete(context.Context, string) (string, error)
}

type Client struct {
	*Resources
}

func NewClient(r *Resources) *Client {
	return &Client{
		r,
	}
}
