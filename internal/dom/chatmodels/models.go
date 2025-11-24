package chatmodels

import (
	"errors"
)

var MissingContentError = errors.New("Missing content")

type ChatModel string

type ImageModel string

const (
	CHAT_MODEL_GEMINI       ChatModel = "gemini"
	CHAT_MODEL_GEMINI_FLASH ChatModel = "gemini flash"
	CHAT_MODEL_GPT          ChatModel = "gpt"
	CHAT_MODEL_META         ChatModel = "llama"
	CHAT_MODEL_TRANSLATIONS ChatModel = "translate"
	CHAT_MODEL_QWEN         ChatModel = "qwen"
	CHAT_MODEL_OPUS         ChatModel = "opus"
	CHAT_MODEL_SONNET       ChatModel = "sonnet"
	CHAT_MODEL_GPT_V4       ChatModel = "g. p. t. version number four"
	CHAT_MODEL_GPT_OSS      ChatModel = "apache"
)

const (
	IMAGE_MODEL_STABLE_DIFFUSION   ImageModel = "stable"
	IMAGE_MODEL_DALL_E_2           ImageModel = "dallas v2"
	IMAGE_MODEL_DALL_E_3           ImageModel = "dallas"
	IMAGE_MODEL_GEMINI             ImageModel = "gemini image"
	IMAGE_MODEL_GEMINI_BANANA_NANO ImageModel = "banana nano"
	IMAGE_MODEL_GPT                ImageModel = "gpt-image"
)

func (c ChatModel) String() string {
	return string(c)
}

func (c ImageModel) String() string {
	return string(c)
}

// Provider identifies which AI service provides the model
type Provider string

const (
	ProviderOpenAI     Provider = "openai"
	ProviderGemini     Provider = "gemini"
	ProviderAnthropic  Provider = "anthropic"
	ProviderCloudflare Provider = "cloudflare"
)

// ModelType distinguishes between chat and image models
type ModelType string

const (
	ModelTypeChat  ModelType = "chat"
	ModelTypeImage ModelType = "image"
)

// ModelConfig contains all configuration for a single model
type ModelConfig struct {
	// The model constant identifier
	ChatModel  ChatModel
	ImageModel ImageModel

	// Type of model (chat or image)
	Type ModelType

	// Provider that owns this model
	Provider Provider

	// Provider-specific model identifier (e.g., "gpt-4o", "claude-sonnet-4-20250514")
	ProviderModelID string

	// Alexa voice command aliases (what users say to select this model)
	Aliases []string

	// Error message when model is unavailable
	ErrorMessage string
}

// ModelRegistry holds all model configurations
type ModelRegistry struct {
	configs []ModelConfig

	// Provider availability flags
	hasOpenAI     bool
	hasGemini     bool
	hasAnthropic  bool
	hasCloudflare bool

	// Cached lookups
	chatModelByAlias  map[string]ModelConfig
	imageModelByAlias map[string]ModelConfig
}

var registry = &ModelRegistry{}

// Model registry - single source of truth for all models
var allModelConfigs = []ModelConfig{
	// OpenAI Chat Models
	{
		ChatModel:       CHAT_MODEL_GPT,
		Type:            ModelTypeChat,
		Provider:        ProviderOpenAI,
		ProviderModelID: "gpt-5.1-2025-11-13",
		Aliases:         []string{string(CHAT_MODEL_GPT)},
		ErrorMessage:    "GPT model is not available - OpenAI API key not configured",
	},
	{
		ChatModel:       CHAT_MODEL_GPT_V4,
		Type:            ModelTypeChat,
		Provider:        ProviderOpenAI,
		ProviderModelID: "gpt-4o",
		Aliases:         []string{string(CHAT_MODEL_GPT_V4)},
		ErrorMessage:    "GPT-4 model is not available - OpenAI API key not configured",
	},

	// OpenAI Image Models
	{
		ImageModel:      IMAGE_MODEL_DALL_E_3,
		Type:            ModelTypeImage,
		Provider:        ProviderOpenAI,
		ProviderModelID: "dall-e-3",
		Aliases:         []string{string(IMAGE_MODEL_DALL_E_3)},
		ErrorMessage:    "DALL-E 3 model is not available - OpenAI API key not configured",
	},
	{
		ImageModel:      IMAGE_MODEL_DALL_E_2,
		Type:            ModelTypeImage,
		Provider:        ProviderOpenAI,
		ProviderModelID: "dall-e-2",
		Aliases:         []string{string(IMAGE_MODEL_DALL_E_2)},
		ErrorMessage:    "DALL-E 2 model is not available - OpenAI API key not configured",
	},
	{
		ImageModel:      IMAGE_MODEL_GPT,
		Type:            ModelTypeImage,
		Provider:        ProviderOpenAI,
		ProviderModelID: "gpt-image-1",
		Aliases:         []string{string(IMAGE_MODEL_GPT)},
		ErrorMessage:    "GPT-Image model is not available - OpenAI API key not configured",
	},

	// Gemini Chat Models
	{
		ChatModel:       CHAT_MODEL_GEMINI,
		Type:            ModelTypeChat,
		Provider:        ProviderGemini,
		ProviderModelID: "gemini-3-pro-preview",
		Aliases:         []string{string(CHAT_MODEL_GEMINI)},
		ErrorMessage:    "Gemini model is not available - Gemini API key not configured",
	},
	{
		ChatModel:       CHAT_MODEL_GEMINI_FLASH,
		Type:            ModelTypeChat,
		Provider:        ProviderGemini,
		ProviderModelID: "gemini-2.5-flash",
		Aliases:         []string{string(CHAT_MODEL_GEMINI_FLASH)},
		ErrorMessage:    "Gemini Flash model is not available - Gemini API key not configured",
	},

	// Gemini Image Models
	{
		ImageModel:      IMAGE_MODEL_GEMINI,
		Type:            ModelTypeImage,
		Provider:        ProviderGemini,
		ProviderModelID: "imagen-4.0-generate-001",
		Aliases:         []string{string(IMAGE_MODEL_GEMINI)},
		ErrorMessage:    "Gemini imagen model is not available - Gemini API key not configured",
	},
	{
		ImageModel:      IMAGE_MODEL_GEMINI_BANANA_NANO,
		Type:            ModelTypeImage,
		Provider:        ProviderGemini,
		ProviderModelID: "gemini-2.5-flash-image-preview",
		Aliases:         []string{string(IMAGE_MODEL_GEMINI_BANANA_NANO)},
		ErrorMessage:    "Gemini banana nano image model is not available - Gemini API key not configured",
	},

	// Anthropic Chat Models
	{
		ChatModel:       CHAT_MODEL_OPUS,
		Type:            ModelTypeChat,
		Provider:        ProviderAnthropic,
		ProviderModelID: "claude-opus-4-20250514",
		Aliases:         []string{string(CHAT_MODEL_OPUS)},
		ErrorMessage:    "Opus model is not available - Anthropic API key not configured",
	},
	{
		ChatModel:       CHAT_MODEL_SONNET,
		Type:            ModelTypeChat,
		Provider:        ProviderAnthropic,
		ProviderModelID: "claude-sonnet-4-20250514",
		Aliases:         []string{string(CHAT_MODEL_SONNET)},
		ErrorMessage:    "Sonnet model is not available - Anthropic API key not configured",
	},

	// Cloudflare Chat Models
	{
		ChatModel:       CHAT_MODEL_META,
		Type:            ModelTypeChat,
		Provider:        ProviderCloudflare,
		ProviderModelID: "@cf/meta/llama-4-scout-17b-16e-instruct",
		Aliases:         []string{string(CHAT_MODEL_META)},
		ErrorMessage:    "Meta model is not available - Cloudflare API key not configured",
	},
	{
		ChatModel:       CHAT_MODEL_QWEN,
		Type:            ModelTypeChat,
		Provider:        ProviderCloudflare,
		ProviderModelID: "@cf/deepseek-ai/deepseek-r1-distill-qwen-32b",
		Aliases:         []string{string(CHAT_MODEL_QWEN)},
		ErrorMessage:    "Qwen model is not available - Cloudflare API key not configured",
	},
	{
		ChatModel:       CHAT_MODEL_GPT_OSS,
		Type:            ModelTypeChat,
		Provider:        ProviderCloudflare,
		ProviderModelID: "@cf/openai/gpt-oss-120b",
		Aliases:         []string{string(CHAT_MODEL_GPT_OSS)},
		ErrorMessage:    "GPT-OSS model is not available - Cloudflare API key not configured",
	},
	{
		ChatModel:       CHAT_MODEL_TRANSLATIONS,
		Type:            ModelTypeChat,
		Provider:        ProviderCloudflare,
		ProviderModelID: "@cf/meta/m2m100-1.2b",
		Aliases:         []string{string(CHAT_MODEL_TRANSLATIONS)},
		ErrorMessage:    "Translation model is not available - Cloudflare API key not configured",
	},

	// Cloudflare Image Models
	{
		ImageModel:      IMAGE_MODEL_STABLE_DIFFUSION,
		Type:            ModelTypeImage,
		Provider:        ProviderCloudflare,
		ProviderModelID: "@cf/stabilityai/stable-diffusion-xl-base-1.0",
		Aliases:         []string{string(IMAGE_MODEL_STABLE_DIFFUSION)},
		ErrorMessage:    "Stable Diffusion model is not available - Cloudflare API key not configured",
	},
}

// RegisterAvailableClients initializes the registry with provider availability
func RegisterAvailableClients(openAI, gemini, anthropic, cloudflare bool) {
	registry.hasOpenAI = openAI
	registry.hasGemini = gemini
	registry.hasAnthropic = anthropic
	registry.hasCloudflare = cloudflare
	registry.configs = allModelConfigs

	// Build alias lookup maps
	registry.chatModelByAlias = make(map[string]ModelConfig)
	registry.imageModelByAlias = make(map[string]ModelConfig)

	// Populate legacy arrays
	AvailableModels = []string{}
	ImageModels = []string{}

	for _, config := range registry.configs {
		// Only include available models
		if !registry.isProviderAvailable(config.Provider) {
			continue
		}

		for _, alias := range config.Aliases {
			if config.Type == ModelTypeChat {
				registry.chatModelByAlias[alias] = config
				AvailableModels = append(AvailableModels, alias)
			} else {
				registry.imageModelByAlias[alias] = config
				ImageModels = append(ImageModels, alias)
			}
		}
	}
}

// IsModelAvailable checks if a chat model is available
func IsModelAvailable(model ChatModel) bool {
	for _, config := range registry.configs {
		if config.ChatModel == model && config.Type == ModelTypeChat {
			return registry.isProviderAvailable(config.Provider)
		}
	}
	return false
}

// IsImageModelAvailable checks if an image model is available
func IsImageModelAvailable(model ImageModel) bool {
	for _, config := range registry.configs {
		if config.ImageModel == model && config.Type == ModelTypeImage {
			return registry.isProviderAvailable(config.Provider)
		}
	}
	return false
}

// isProviderAvailable checks if a provider is configured
func (r *ModelRegistry) isProviderAvailable(provider Provider) bool {
	switch provider {
	case ProviderOpenAI:
		return r.hasOpenAI
	case ProviderGemini:
		return r.hasGemini
	case ProviderAnthropic:
		return r.hasAnthropic
	case ProviderCloudflare:
		return r.hasCloudflare
	default:
		return false
	}
}

// GetChatModelByAlias returns a chat model config by its alias
func GetChatModelByAlias(alias string) (ModelConfig, bool) {
	config, ok := registry.chatModelByAlias[alias]
	return config, ok
}

// GetImageModelByAlias returns an image model config by its alias
func GetImageModelByAlias(alias string) (ModelConfig, bool) {
	config, ok := registry.imageModelByAlias[alias]
	return config, ok
}

// GetAvailableChatModels returns all available chat model aliases
func GetAvailableChatModels() []string {
	var models []string
	for _, config := range registry.configs {
		if config.Type == ModelTypeChat && registry.isProviderAvailable(config.Provider) {
			models = append(models, config.Aliases...)
		}
	}
	return models
}

// GetAvailableImageModels returns all available image model aliases
func GetAvailableImageModels() []string {
	var models []string
	for _, config := range registry.configs {
		if config.Type == ModelTypeImage && registry.isProviderAvailable(config.Provider) {
			models = append(models, config.Aliases...)
		}
	}
	return models
}

// GetAvailableChatModelsWithProviderIDs returns all available chat models with their provider model IDs
// Returns a map of alias -> provider model ID
func GetAvailableChatModelsWithProviderIDs() map[string]string {
	models := make(map[string]string)
	for _, config := range registry.configs {
		if config.Type == ModelTypeChat && registry.isProviderAvailable(config.Provider) {
			for _, alias := range config.Aliases {
				models[alias] = config.ProviderModelID
			}
		}
	}
	return models
}

// GetAvailableImageModelsWithProviderIDs returns all available image models with their provider model IDs
// Returns a map of alias -> provider model ID
func GetAvailableImageModelsWithProviderIDs() map[string]string {
	models := make(map[string]string)
	for _, config := range registry.configs {
		if config.Type == ModelTypeImage && registry.isProviderAvailable(config.Provider) {
			for _, alias := range config.Aliases {
				models[alias] = config.ProviderModelID
			}
		}
	}
	return models
}

// GetProviderModelID returns the provider-specific model ID for a chat model
func GetProviderModelID(model ChatModel) (string, bool) {
	// Use allModelConfigs directly (works even before RegisterAvailableClients is called)
	for _, config := range allModelConfigs {
		if config.ChatModel == model && config.Type == ModelTypeChat {
			return config.ProviderModelID, true
		}
	}
	return "", false
}

// GetImageProviderModelID returns the provider-specific model ID for an image model
func GetImageProviderModelID(model ImageModel) (string, bool) {
	// Use allModelConfigs directly (works even before RegisterAvailableClients is called)
	for _, config := range allModelConfigs {
		if config.ImageModel == model && config.Type == ModelTypeImage {
			return config.ProviderModelID, true
		}
	}
	return "", false
}

// GetChatModelProvider returns the provider for a chat model
func GetChatModelProvider(model ChatModel) (Provider, bool) {
	for _, config := range allModelConfigs {
		if config.ChatModel == model && config.Type == ModelTypeChat {
			return config.Provider, true
		}
	}
	return "", false
}

// GetImageModelProvider returns the provider for an image model
func GetImageModelProvider(model ImageModel) (Provider, bool) {
	for _, config := range allModelConfigs {
		if config.ImageModel == model && config.Type == ModelTypeImage {
			return config.Provider, true
		}
	}
	return "", false
}

// Legacy compatibility exports - populated by RegisterAvailableClients
var AvailableModels []string
var ImageModels []string
