package chatmodels

import (
	"context"
	"fmt"
)

func (client *Client) TextGeneration(ctx context.Context, prompt string, model ChatModel) (string, error) {
	switch model {
	case CHAT_MODEL_GEMINI:
		return client.GeminiAPI.GenerateText(ctx, prompt)
	case CHAT_MODEL_META, CHAT_MODEL_SQL, CHAT_MODEL_OPEN, CHAT_MODEL_AWQ, CHAT_MODEL_QWEN:
		return client.CloudflareApiClient.GenerateTextWithModel(ctx, prompt, CHAT_MODEL_TO_CF_MODEL[model])
	case CHAT_MODEL_GPT:
		fallthrough
	default:
		return client.GPTApi.GenerateTextWithModel(ctx, prompt, CHAT_MODEL_TO_OPENAI_MODEL[model])
	}
}

func (client *Client) GenerateImage(ctx context.Context, prompt string, model ImageModel) ([]byte, error) {
	switch model {
	case IMAGE_MODEL_STABLE_DIFFUSION:
		return client.CloudflareApiClient.GenerateImage(ctx, prompt, IMAGE_MODEL_TO_CF_MODEL[model])
	case IMAGE_MODEL_DALL_E_2, IMAGE_MODEL_DALL_E_3:
		return client.GPTApi.GenerateImage(ctx, prompt, IMAGE_MODEL_TO_OPENAI_MODEL[model])
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
