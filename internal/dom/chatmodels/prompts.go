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
	// Get provider and model ID from centralized registry
	provider, ok := GetChatModelProvider(model)
	if !ok {
		return nil, nil
	}

	providerModelID, ok := GetProviderModelID(model)
	if !ok {
		return nil, nil
	}

	// Route to correct API client based on provider
	switch provider {
	case ProviderAnthropic:
		if client.AnthropicAPI == nil {
			return nil, nil
		}
		opts := []llms.CallOption{llms.WithModel(providerModelID)}
		return client.AnthropicAPI.GetModel(opts...), opts
	case ProviderCloudflare:
		if client.CloudflareApiClient == nil {
			return nil, nil
		}
		opts := []llms.CallOption{llms.WithModel(providerModelID)}
		return client.CloudflareApiClient.GetModel(opts...), opts
	case ProviderGemini:
		if client.GeminiAPI == nil {
			return nil, nil
		}
		opts := []llms.CallOption{llms.WithModel(providerModelID)}
		return client.GeminiAPI.GetModel(opts...), opts
	case ProviderOpenAI:
		if client.GPTApi == nil {
			return nil, nil
		}
		opts := []llms.CallOption{llms.WithModel(providerModelID), llms.WithTemperature(1)}
		return client.GPTApi.GetModel(opts...), opts
	default:
		return nil, nil
	}
}

func (client *Client) GenerateImage(ctx context.Context, prompt string, model ImageModel) ([]byte, error) {
	// Get provider and model ID from centralized registry
	provider, ok := GetImageModelProvider(model)
	if !ok {
		return nil, fmt.Errorf("image model %s is not configured", model)
	}

	providerModelID, ok := GetImageProviderModelID(model)
	if !ok {
		return nil, fmt.Errorf("image model %s is not configured", model)
	}

	// Route to correct API client based on provider
	switch provider {
	case ProviderOpenAI:
		if client.GPTApi == nil {
			return nil, fmt.Errorf("image model %s is not available: OpenAI client not configured", model)
		}
		return client.GPTApi.GenerateImage(ctx, prompt, providerModelID)
	case ProviderGemini:
		if client.GeminiAPI == nil {
			return nil, fmt.Errorf("image model %s is not available: Gemini client not configured", model)
		}
		return client.GeminiAPI.GenerateImage(ctx, prompt, providerModelID)
	case ProviderCloudflare:
		if client.CloudflareApiClient == nil {
			return nil, fmt.Errorf("image model %s is not available: Cloudflare client not configured", model)
		}
		return client.CloudflareApiClient.GenerateImage(ctx, prompt, providerModelID)
	default:
		return nil, fmt.Errorf("image model %s has unsupported provider %s", model, provider)
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

	// Get provider model ID from centralized registry
	providerModelID, ok := GetProviderModelID(model)
	if !ok {
		return "", fmt.Errorf("translation model %s is not configured", model)
	}

	return client.CloudflareApiClient.GenerateTranslation(ctx, &GenerateTranslationRequest{
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
		Prompt:         prompt,
		Model:          providerModelID,
	})
}
