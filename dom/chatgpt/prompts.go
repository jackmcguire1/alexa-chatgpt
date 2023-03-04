package chatgpt

import (
	"context"
	"fmt"
)

func (client *Client) AutoComplete(ctx context.Context, prompt string) (string, error) {
	resp, err := client.Api.AutoComplete(ctx, prompt)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		err = fmt.Errorf("missing choices")
		return "", err
	}

	choice := resp.Choices[0]
	return choice.Text, nil
}
