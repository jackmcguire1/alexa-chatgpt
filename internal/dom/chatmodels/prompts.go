package chatmodels

import (
	"context"

	"github.com/tmc/langchaingo/llms"
)

func (client *Client) TextGeneration(ctx context.Context, prompt string, model ChatModel) (string, error) {
	modelSvc, opts := client.GetLLmModel(model)
	return llms.GenerateFromSinglePrompt(ctx, modelSvc, prompt, opts...)
}

func (client *Client) GetLLmModel(model ChatModel) (llms.Model, []llms.CallOption) {
	switch model {
	case CHAT_MODEL_META, CHAT_MODEL_SQL, CHAT_MODEL_OPEN, CHAT_MODEL_AWQ, CHAT_MODEL_QWEN:
		return client.CloudflareApiClient.GetModel(), []llms.CallOption{llms.WithModel(CHAT_MODEL_TO_CF_MODEL[model])}
	case CHAT_MODEL_GEMINI:
		return client.GeminiAPI.GetModel(), []llms.CallOption{llms.WithModel(VERTEX_MODEL)}
	default:
		return client.GPTApi.GetModel(), []llms.CallOption{llms.WithModel(CHAT_MODEL_TO_OPENAI_MODEL[model])}
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
