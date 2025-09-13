package init

import (
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
)

// InitializeResources creates and configures all AI provider clients based on environment variables
func InitializeResources() *chatmodels.Resources {
	resources := &chatmodels.Resources{}

	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey != "" {
		resources.GPTApi = chatmodels.NewOpenAiApiClient(openAIKey)
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	if geminiKey != "" {
		resources.GeminiAPI = chatmodels.NewGeminiApiClient(geminiKey)
	}

	cloudflareAccountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	cloudflareAPIKey := os.Getenv("CLOUDFLARE_API_KEY")
	if cloudflareAccountID != "" && cloudflareAPIKey != "" {
		resources.CloudflareApiClient = chatmodels.NewCloudflareApiClient(cloudflareAccountID, cloudflareAPIKey)
	}

	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicKey != "" {
		resources.AnthropicAPI = chatmodels.NewAnthropicApiClient(anthropicKey)
	}

	chatmodels.RegisterAvailableClients(
		resources.GPTApi != nil,
		resources.GeminiAPI != nil,
		resources.AnthropicAPI != nil,
		resources.CloudflareApiClient != nil,
	)

	return resources
}

// GetDefaultChatModel returns the default chat model based on available clients
func GetDefaultChatModel(resources *chatmodels.Resources) chatmodels.ChatModel {
	if len(chatmodels.AvailableModels) == 0 {
		return chatmodels.CHAT_MODEL_GPT
	}

	if resources.GPTApi != nil {
		return chatmodels.CHAT_MODEL_GPT
	}
	if resources.GeminiAPI != nil {
		return chatmodels.CHAT_MODEL_GEMINI
	}
	if resources.AnthropicAPI != nil {
		return chatmodels.CHAT_MODEL_OPUS
	}
	if resources.CloudflareApiClient != nil {
		return chatmodels.CHAT_MODEL_META
	}

	return chatmodels.CHAT_MODEL_GPT
}

// GetDefaultImageModel returns the default image model based on available clients
func GetDefaultImageModel(resources *chatmodels.Resources) chatmodels.ImageModel {
	if len(chatmodels.ImageModels) == 0 {
		return chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
	}

	if resources.CloudflareApiClient != nil {
		return chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
	}
	if resources.GPTApi != nil {
		return chatmodels.IMAGE_MODEL_DALL_E_3
	}
	if resources.GeminiAPI != nil {
		return chatmodels.IMAGE_MODEL_GEMINI
	}

	return chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
}
