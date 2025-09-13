package chatmodels

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

func (client *Client) TextGeneration(ctx context.Context, prompt string, model ChatModel) (string, error) {
	modelSvc, opts := client.GetLLmModel(model)
	if modelSvc == nil {
		return "", fmt.Errorf("model %s is not available: required client not configured", model)
	}
	return llms.GenerateFromSinglePrompt(ctx, modelSvc, prompt, opts...)
}

func (client *Client) TextGenerationWithSystem(ctx context.Context, system string, prompt string, model ChatModel) (result string, err error) {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, system),
		llms.TextParts(llms.ChatMessageTypeHuman, prompt),
	}
	llmModel, opts := client.GetLLmModel(model)
	if llmModel == nil {
		return "", fmt.Errorf("model %s is not available: required client not configured", model)
	}
	res, err := llmModel.GenerateContent(ctx, content, opts...)
	if err != nil {
		return "", err
	}
	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no choices in response from llm model %v", llmModel)
	}
	return res.Choices[0].Content, nil
}

func (client *Client) GetLLmModel(model ChatModel) (llms.Model, []llms.CallOption) {
	switch model {
	case CHAT_MODEL_OPUS, CHAT_MODEL_SONNET:
		if client.AnthropicAPI == nil {
			return nil, nil
		}
		return client.AnthropicAPI.GetModel(), []llms.CallOption{llms.WithModel(CHAT_MODEL_TO_ANTHROPIC_MODELS[model])}
	case CHAT_MODEL_META, CHAT_MODEL_QWEN, CHAT_MODEL_GPT_OSS:
		if client.CloudflareApiClient == nil {
			return nil, nil
		}
		return client.CloudflareApiClient.GetModel(), []llms.CallOption{llms.WithModel(CHAT_MODEL_TO_CF_MODEL[model])}
	case CHAT_MODEL_GEMINI:
		if client.GeminiAPI == nil {
			return nil, nil
		}
		return client.GeminiAPI.GetModel(), []llms.CallOption{llms.WithModel(VERTEX_MODEL)}
	default:
		if client.GPTApi == nil {
			return nil, nil
		}
		return client.GPTApi.GetModel(), []llms.CallOption{llms.WithModel(CHAT_MODEL_TO_OPENAI_MODEL[model]), llms.WithTemperature(1)}
	}
}

func (client *Client) GenerateImage(ctx context.Context, prompt string, model ImageModel) ([]byte, error) {
	switch model {
	case IMAGE_MODEL_DALL_E_2, IMAGE_MODEL_DALL_E_3:
		if client.GPTApi == nil {
			return nil, fmt.Errorf("image model %s is not available: OpenAI client not configured", model)
		}
		return client.GPTApi.GenerateImage(ctx, prompt, IMAGE_MODEL_TO_OPENAI_MODEL[model])
	case IMAGE_MODEL_GEMINI, IMAGE_MODEL_GEMINI_BANANA_NANO:
		if client.GeminiAPI == nil {
			return nil, fmt.Errorf("image model %s is not available: Gemini client not configured", model)
		}
		return client.GeminiAPI.GenerateImage(ctx, prompt, IMAGE_MODEL_TO_GEMINI_MODEL[model])
	case IMAGE_MODEL_STABLE_DIFFUSION:
		fallthrough
	default:
		// Both STABLE_DIFFUSION and default case use Cloudflare
		if client.CloudflareApiClient == nil {
			return nil, fmt.Errorf("image model %s is not available: Cloudflare client not configured", model)
		}
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
	if client.CloudflareApiClient == nil {
		return "", fmt.Errorf("translation model is not available: Cloudflare client not configured")
	}
	return client.CloudflareApiClient.GenerateTranslation(ctx, &GenerateTranslationRequest{
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
		Prompt:         prompt,
		Model:          CHAT_MODEL_TO_CF_MODEL[model],
	})
}
