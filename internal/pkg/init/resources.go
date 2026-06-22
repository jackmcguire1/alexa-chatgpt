package init

import (
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
)

// InitializeResources creates and configures all AI provider clients based on environment variables.
func InitializeResources() *chatmodels.Resources {
	resources := &chatmodels.Resources{}

	resources.BedrockAPI = chatmodels.NewBedrockApiClient()
	resources.MantleAPI = chatmodels.NewMantleApiClient()

	if accountID, apiKey := os.Getenv("CLOUDFLARE_ACCOUNT_ID"), os.Getenv("CLOUDFLARE_API_KEY"); accountID != "" && apiKey != "" {
		resources.CloudflareAPI = chatmodels.NewCloudflareApiClient(accountID, apiKey)
	}

	return resources
}

// GetDefaultChatModel returns the default chat model.
func GetDefaultChatModel() chatmodels.ChatModel {
	return chatmodels.CHAT_MODEL_SONNET
}

// GetDefaultImageModel returns the default image model.
func GetDefaultImageModel() chatmodels.ImageModel {
	return chatmodels.IMAGE_MODEL_NOVA_CANVAS
}
