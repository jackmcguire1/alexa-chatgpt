package chatmodels

import (
	"context"
	"fmt"

	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

func (client *Client) AutoComplete(ctx context.Context, prompt string, model ChatModel) (string, error) {
	switch model {
	case CHAT_MODEL_GEMINI:
		res, err := client.GeminiAPI.GeminiChat(ctx, prompt)
		if err != nil {
			return "", err
		}

		if len(res.Candidates) > 0 && len(res.Candidates[0].Content.Parts) > 0 {
			return fmt.Sprint(res.Candidates[0].Content.Parts[0]), nil
		}
		return "", fmt.Errorf("did not get enough info back from google %s", utils.ToJSON(res))

	default:
		resp, err := client.GPTApi.AutoComplete(ctx, prompt)
		if err != nil {
			return "", err
		}

		if len(resp.Choices) == 0 {
			err = fmt.Errorf("missing choices")
			return "", err
		}

		message := resp.Choices[0].Message.Content
		return message, nil
	}
}
