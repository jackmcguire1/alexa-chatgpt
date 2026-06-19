package chatmodels

import "errors"

var MissingContentError = errors.New("Missing content")

type ChatModel string

type ImageModel string

const (
	CHAT_MODEL_SONNET       ChatModel = "sonnet"
	CHAT_MODEL_OPUS         ChatModel = "opus"
	CHAT_MODEL_FABLE        ChatModel = "fable"
	CHAT_MODEL_TRANSLATIONS ChatModel = "translate"
	CHAT_MODEL_NOVA_LITE    ChatModel = "nova"
	CHAT_MODEL_NOVA_PRO     ChatModel = "nova pro"
	CHAT_MODEL_GROK         ChatModel = "grok"
	CHAT_MODEL_GPT          ChatModel = "gpt"
)

const (
	IMAGE_MODEL_NOVA_CANVAS ImageModel = "nova canvas"
	IMAGE_MODEL_TITAN       ImageModel = "titan"
)

func (c ChatModel) String() string {
	return string(c)
}

func (c ImageModel) String() string {
	return string(c)
}

// Provider identifies which AI service provides the model.
type Provider string

const (
	ProviderBedrock       Provider = "bedrock"
	ProviderBedrockMantle Provider = "bedrock-mantle"
)

// ModelType distinguishes between chat and image models.
type ModelType string

const (
	ModelTypeChat  ModelType = "chat"
	ModelTypeImage ModelType = "image"
)

// ModelConfig contains all configuration for a single model.
type ModelConfig struct {
	ChatModel  ChatModel
	ImageModel ImageModel

	Type     ModelType
	Provider Provider

	// Provider-specific model identifier used in the API call.
	ProviderModelID string

	// MantleRegion is the AWS region for ProviderBedrockMantle models.
	// Each mantle model may only be available in a specific region.
	MantleRegion string

	// Alexa voice command aliases (what users say to select this model).
	Aliases []string

	ErrorMessage string
}

// ModelRegistry holds all model configurations.
type ModelRegistry struct {
	configs    []ModelConfig
	hasBedrock bool

	chatModelByAlias  map[string]ModelConfig
	imageModelByAlias map[string]ModelConfig
}

var registry = &ModelRegistry{}

// allModelConfigs is the single source of truth for all models.
// All models use AWS Bedrock cross-region inference profile IDs.
var allModelConfigs = []ModelConfig{
	// Claude models
	{
		ChatModel:       CHAT_MODEL_SONNET,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrock,
		ProviderModelID: "us.anthropic.claude-sonnet-4-6",
		Aliases:         []string{string(CHAT_MODEL_SONNET)},
		ErrorMessage:    "Sonnet model is not available - Bedrock not configured",
	},
	{
		ChatModel:       CHAT_MODEL_OPUS,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrock,
		ProviderModelID: "us.anthropic.claude-opus-4-8",
		Aliases:         []string{string(CHAT_MODEL_OPUS)},
		ErrorMessage:    "Opus model is not available - Bedrock not configured",
	},
	{
		ChatModel:       CHAT_MODEL_FABLE,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrock,
		ProviderModelID: "us.anthropic.claude-fable-5",
		Aliases:         []string{string(CHAT_MODEL_FABLE)},
		ErrorMessage:    "Fable model is not available - Bedrock not configured",
	},

	// Amazon Nova
	{
		ChatModel:       CHAT_MODEL_NOVA_LITE,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrock,
		ProviderModelID: "us.amazon.nova-lite-v1:0",
		Aliases:         []string{string(CHAT_MODEL_NOVA_LITE)},
		ErrorMessage:    "Nova Lite model is not available - Bedrock not configured",
	},
	{
		ChatModel:       CHAT_MODEL_NOVA_PRO,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrock,
		ProviderModelID: "us.amazon.nova-pro-v1:0",
		Aliases:         []string{string(CHAT_MODEL_NOVA_PRO)},
		ErrorMessage:    "Nova Pro model is not available - Bedrock not configured",
	},

	// Translation routed through Sonnet with a system prompt.
	{
		ChatModel:       CHAT_MODEL_TRANSLATIONS,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrock,
		ProviderModelID: "us.anthropic.claude-sonnet-4-6",
		Aliases:         []string{string(CHAT_MODEL_TRANSLATIONS)},
		ErrorMessage:    "Translation model is not available - Bedrock not configured",
	},

	// Bedrock Mantle models (OpenAI-compatible endpoint for third-party providers)
	{
		ChatModel:       CHAT_MODEL_GROK,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrockMantle,
		ProviderModelID: "xai.grok-4.3",
		MantleRegion:    "us-west-2",
		Aliases:         []string{string(CHAT_MODEL_GROK)},
		ErrorMessage:    "Grok model is not available - Bedrock not configured",
	},
	{
		ChatModel:       CHAT_MODEL_GPT,
		Type:            ModelTypeChat,
		Provider:        ProviderBedrockMantle,
		ProviderModelID: "openai.gpt-5.5",
		MantleRegion:    "us-east-1",
		Aliases:         []string{string(CHAT_MODEL_GPT)},
		ErrorMessage:    "GPT model is not available - Bedrock not configured",
	},

	// Image models
	{
		ImageModel:      IMAGE_MODEL_NOVA_CANVAS,
		Type:            ModelTypeImage,
		Provider:        ProviderBedrock,
		ProviderModelID: "amazon.nova-canvas-v1:0",
		Aliases:         []string{string(IMAGE_MODEL_NOVA_CANVAS)},
		ErrorMessage:    "Nova Canvas model is not available - Bedrock not configured",
	},
	{
		ImageModel:      IMAGE_MODEL_TITAN,
		Type:            ModelTypeImage,
		Provider:        ProviderBedrock,
		ProviderModelID: "amazon.titan-image-generator-v2:0",
		Aliases:         []string{string(IMAGE_MODEL_TITAN)},
		ErrorMessage:    "Titan Image Generator model is not available - Bedrock not configured",
	},
}

// RegisterAvailableClients initialises the model registry.
func RegisterAvailableClients() {
	registry.hasBedrock = true
	registry.configs = allModelConfigs

	registry.chatModelByAlias = make(map[string]ModelConfig)
	registry.imageModelByAlias = make(map[string]ModelConfig)

	AvailableModels = []string{}
	ImageModels = []string{}

	for _, cfg := range registry.configs {
		if !registry.isProviderAvailable(cfg.Provider) {
			continue
		}
		for _, alias := range cfg.Aliases {
			if cfg.Type == ModelTypeChat {
				registry.chatModelByAlias[alias] = cfg
				AvailableModels = append(AvailableModels, alias)
			} else {
				registry.imageModelByAlias[alias] = cfg
				ImageModels = append(ImageModels, alias)
			}
		}
	}
}

// IsModelAvailable checks if a chat model is available.
func IsModelAvailable(model ChatModel) bool {
	for _, cfg := range registry.configs {
		if cfg.ChatModel == model && cfg.Type == ModelTypeChat {
			return registry.isProviderAvailable(cfg.Provider)
		}
	}
	return false
}

// IsImageModelAvailable checks if an image model is available.
func IsImageModelAvailable(model ImageModel) bool {
	for _, cfg := range registry.configs {
		if cfg.ImageModel == model && cfg.Type == ModelTypeImage {
			return registry.isProviderAvailable(cfg.Provider)
		}
	}
	return false
}

func (r *ModelRegistry) isProviderAvailable(provider Provider) bool {
	switch provider {
	case ProviderBedrock, ProviderBedrockMantle:
		return r.hasBedrock
	default:
		return false
	}
}

// GetChatModelByAlias returns a chat model config by its alias.
func GetChatModelByAlias(alias string) (ModelConfig, bool) {
	cfg, ok := registry.chatModelByAlias[alias]
	return cfg, ok
}

// GetImageModelByAlias returns an image model config by its alias.
func GetImageModelByAlias(alias string) (ModelConfig, bool) {
	cfg, ok := registry.imageModelByAlias[alias]
	return cfg, ok
}

// GetAvailableChatModels returns all available chat model aliases.
func GetAvailableChatModels() []string {
	var models []string
	for _, cfg := range registry.configs {
		if cfg.Type == ModelTypeChat && registry.isProviderAvailable(cfg.Provider) {
			models = append(models, cfg.Aliases...)
		}
	}
	return models
}

// GetAvailableImageModels returns all available image model aliases.
func GetAvailableImageModels() []string {
	var models []string
	for _, cfg := range registry.configs {
		if cfg.Type == ModelTypeImage && registry.isProviderAvailable(cfg.Provider) {
			models = append(models, cfg.Aliases...)
		}
	}
	return models
}

// GetAvailableChatModelsWithProviderIDs returns a map of alias -> provider model ID.
func GetAvailableChatModelsWithProviderIDs() map[string]string {
	models := make(map[string]string)
	for _, cfg := range registry.configs {
		if cfg.Type == ModelTypeChat && registry.isProviderAvailable(cfg.Provider) {
			for _, alias := range cfg.Aliases {
				models[alias] = cfg.ProviderModelID
			}
		}
	}
	return models
}

// GetAvailableImageModelsWithProviderIDs returns a map of alias -> provider model ID.
func GetAvailableImageModelsWithProviderIDs() map[string]string {
	models := make(map[string]string)
	for _, cfg := range registry.configs {
		if cfg.Type == ModelTypeImage && registry.isProviderAvailable(cfg.Provider) {
			for _, alias := range cfg.Aliases {
				models[alias] = cfg.ProviderModelID
			}
		}
	}
	return models
}

// GetChatModelConfig returns the full config for a chat model.
func GetChatModelConfig(model ChatModel) (ModelConfig, bool) {
	for _, cfg := range allModelConfigs {
		if cfg.ChatModel == model && cfg.Type == ModelTypeChat {
			return cfg, true
		}
	}
	return ModelConfig{}, false
}

// GetProviderModelID returns the provider-specific model ID for a chat model.
func GetProviderModelID(model ChatModel) (string, bool) {
	for _, cfg := range allModelConfigs {
		if cfg.ChatModel == model && cfg.Type == ModelTypeChat {
			return cfg.ProviderModelID, true
		}
	}
	return "", false
}

// GetImageProviderModelID returns the provider-specific model ID for an image model.
func GetImageProviderModelID(model ImageModel) (string, bool) {
	for _, cfg := range allModelConfigs {
		if cfg.ImageModel == model && cfg.Type == ModelTypeImage {
			return cfg.ProviderModelID, true
		}
	}
	return "", false
}

// Legacy compatibility exports populated by RegisterAvailableClients.
var AvailableModels []string
var ImageModels []string
