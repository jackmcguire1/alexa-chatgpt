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

	case CHAT_MODEL_META:
		return client.CloudflareApiClient.GenerateText(ctx, prompt, CF_LLAMA_2_7B_CHAT_INT8_MODEL)
	case CHAT_MODEL_SQL:
		return client.CloudflareApiClient.GenerateText(ctx, prompt, CF_SQL_MODEL)
	case CHAT_MODEL_OPEN:
		return client.CloudflareApiClient.GenerateText(ctx, prompt, CF_OPEN_CHAT_MODEL)
	case CF_AWQ_MODEL:
		return client.CloudflareApiClient.GenerateText(ctx, prompt, CF_AWQ_MODEL)
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
