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
	case CHAT_MODEL_META, CHAT_MODEL_SQL, CHAT_MODEL_OPEN, CHAT_MODEL_AWQ, CHAT_MODEL_QWEN:
		return client.CloudflareApiClient.GenerateText(ctx, prompt, CHAT_MODEL_TO_CF_MODEL[model])
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

func (client *Client) GenerateImage(ctx context.Context, prompt string, model ImageModel) ([]byte, error) {
	switch model {
	case IMAGE_MODEL_STABLE_DIFFUSION:
		return client.CloudflareApiClient.GenerateImage(ctx, prompt, IMAGE_MODEL_TO_CF_MODEL[model])
	case IMAGE_MODEL_DALL_E_2:
		return client.GPTApi.GenerateImage(ctx, prompt, IMAGE_MODEL_TO_CF_MODEL[model])
	default:
		return client.CloudflareApiClient.GenerateImage(ctx, prompt, IMAGE_MODEL_TO_CF_MODEL[model])
	}
	return nil, fmt.Errorf("unidentified image generation model")
}

func (client *Client) Translate(
	ctx context.Context,
	prompt string,
	sourceLang string,
	targetLang string,
	model ChatModel,
) (string, error) {
	if sourceLang == "" {
		sourceLang = "en"
	}
	if targetLang == "" {
		targetLang = "jp"
	}
	if model == "" {
		model = CHAT_MODEL_TRANSLATIONS
	}
	return client.CloudflareApiClient.GenerateTranslation(ctx, &GenerateTranslationRequest{
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
		Prompt:         prompt,
		Model:          CHAT_MODEL_TO_CF_MODEL[model],
	})
}
