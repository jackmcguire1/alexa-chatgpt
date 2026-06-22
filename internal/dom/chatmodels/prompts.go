package chatmodels

import (
	"context"
	"fmt"
)

// generateContent routes to the appropriate backend based on the model's provider.
func (client *Client) generateContent(ctx context.Context, messages []Message, model ChatModel) (string, error) {
	cfg, ok := GetChatModelConfig(model)
	if !ok {
		return "", fmt.Errorf("model %s is not configured", model)
	}
	opts := GenerateOptions{Model: cfg.ProviderModelID, MantleRegion: cfg.MantleRegion}

	switch cfg.Provider {
	case ProviderBedrockMantle:
		if client.MantleAPI == nil {
			return "", fmt.Errorf("model %s is not available: Mantle client not configured", model)
		}
		resp, err := client.MantleAPI.GenerateContent(ctx, messages, opts)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	case ProviderCloudflare:
		if client.CloudflareAPI == nil {
			return "", fmt.Errorf("model %s is not available: Cloudflare client not configured", model)
		}
		resp, err := client.CloudflareAPI.GenerateContent(ctx, messages, opts)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	default:
		if client.BedrockAPI == nil {
			return "", fmt.Errorf("model %s is not available: Bedrock client not configured", model)
		}
		resp, err := client.BedrockAPI.GenerateContent(ctx, messages, opts)
		if err != nil {
			return "", err
		}
		return resp.Content, nil
	}
}

func (client *Client) TextGeneration(ctx context.Context, prompt string, model ChatModel) (string, error) {
	return client.generateContent(ctx, []Message{
		{Role: RoleUser, Content: prompt},
	}, model)
}

func (client *Client) TextGenerationWithSystem(ctx context.Context, system string, prompt string, model ChatModel) (string, error) {
	return client.generateContent(ctx, []Message{
		{Role: RoleSystem, Content: system},
		{Role: RoleUser, Content: prompt},
	}, model)
}

func (client *Client) GenerateImage(ctx context.Context, prompt string, model ImageModel) ([]byte, error) {
	cfg, ok := GetImageModelConfig(model)
	if !ok {
		return nil, fmt.Errorf("image model %s is not configured", model)
	}

	switch cfg.Provider {
	case ProviderCloudflare:
		if client.CloudflareAPI == nil {
			return nil, fmt.Errorf("image model %s is not available: Cloudflare client not configured", model)
		}
		return client.CloudflareAPI.GenerateImage(ctx, prompt, cfg.ProviderModelID)
	default:
		if client.BedrockAPI == nil {
			return nil, fmt.Errorf("image model %s is not available: Bedrock client not configured", model)
		}
		return client.BedrockAPI.GenerateImage(ctx, prompt, cfg.ProviderModelID)
	}
}

func (client *Client) Translate(
	ctx context.Context,
	prompt string,
	sourceLang string,
	targetLang string,
	_ ChatModel,
) (string, error) {
	if sourceLang == "" {
		sourceLang = "en"
	}
	if targetLang == "" {
		targetLang = "jp"
	}
	systemPrompt := fmt.Sprintf(
		"You are a translator. Translate the following text from %s to %s. Output only the translated text, nothing else.",
		sourceLang, targetLang,
	)
	return client.TextGenerationWithSystem(ctx, systemPrompt, prompt, CHAT_MODEL_SONNET)
}
