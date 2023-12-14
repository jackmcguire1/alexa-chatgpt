package chatmodels

import "context"

type Resources struct {
	GPTApi    GptAPI
	GoogleAPI GoogleApi
}

type Service interface {
	AutoComplete(context.Context, string, ChatModel) (string, error)
}

type Client struct {
	*Resources
}

func NewClient(r *Resources) *Client {
	return &Client{
		r,
	}
}
