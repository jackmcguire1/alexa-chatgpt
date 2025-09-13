package api

import (
	"fmt"
	"strings"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
)

const (
	responseTitleChatModels  = "Chat Models"
	responseTitleImageModels = "Image Models"
	responseTitleModels      = "Models"
	responseOK               = "ok"
)

type modelConfig struct {
	model        chatmodels.ChatModel
	errorMessage string
}

type imageModelConfig struct {
	model        chatmodels.ImageModel
	errorMessage string
}

var chatModelConfigs = map[string]modelConfig{
	"gemini":                       {chatmodels.CHAT_MODEL_GEMINI, "Gemini model is not available - Gemini API key not configured"},
	"gpt":                          {chatmodels.CHAT_MODEL_GPT, "GPT model is not available - OpenAI API key not configured"},
	"g. p. t. version number four": {chatmodels.CHAT_MODEL_GPT_V4, "GPT-4 model is not available - OpenAI API key not configured"},
	"llama":                        {chatmodels.CHAT_MODEL_META, "Meta model is not available - Cloudflare API key not configured"},
	"qwen":                         {chatmodels.CHAT_MODEL_QWEN, "Qwen model is not available - Cloudflare API key not configured"},
	"opus":                         {chatmodels.CHAT_MODEL_OPUS, "Opus model is not available - Anthropic API key not configured"},
	"sonnet":                       {chatmodels.CHAT_MODEL_SONNET, "Sonnet model is not available - Anthropic API key not configured"},
	"apache":                       {chatmodels.CHAT_MODEL_GPT_OSS, "GPT-OSS model is not available - Cloudflare API key not configured"},
}

var imageModelConfigs = map[string]imageModelConfig{
	"stable":       {chatmodels.IMAGE_MODEL_STABLE_DIFFUSION, "Stable Diffusion model is not available - Cloudflare API key not configured"},
	"dallas":       {chatmodels.IMAGE_MODEL_DALL_E_3, "DALL-E 3 model is not available - OpenAI API key not configured"},
	"dallas v2":    {chatmodels.IMAGE_MODEL_DALL_E_2, "DALL-E 2 model is not available - OpenAI API key not configured"},
	"gemini image": {chatmodels.IMAGE_MODEL_GEMINI, "Gemini image model is not available - Gemini API key not configured"},
}

func (h *Handler) setChatModel(config modelConfig) alexa.Response {
	if !chatmodels.IsModelAvailable(config.model) {
		return alexa.NewResponse(responseTitleChatModels, config.errorMessage, false)
	}
	h.Model = config.model
	return alexa.NewResponse(responseTitleChatModels, responseOK, false)
}

func (h *Handler) setImageModel(config imageModelConfig) alexa.Response {
	if !chatmodels.IsImageModelAvailable(config.model) {
		return alexa.NewResponse(responseTitleImageModels, config.errorMessage, false)
	}
	h.ImageModel = config.model
	return alexa.NewResponse(responseTitleImageModels, responseOK, false)
}

func (h *Handler) getOrSetModel(model string) (res alexa.Response, err error) {
	lowerModel := strings.ToLower(model)

	// Check chat models
	if config, ok := chatModelConfigs[lowerModel]; ok {
		return h.setChatModel(config), nil
	}

	// Check image models
	if config, ok := imageModelConfigs[lowerModel]; ok {
		return h.setImageModel(config), nil
	}

	// Special cases
	switch lowerModel {
	case "which":
		res = alexa.NewResponse(responseTitleChatModels,
			fmt.Sprintf("I am using the text-model %s and image-model %s",
				h.Model.String(), h.ImageModel.String()), false)
		return
	default:
		res = alexa.NewResponse(
			responseTitleModels,
			fmt.Sprintf("The available chat models are \n - %s and image models %s",
				strings.Join(chatmodels.AvailableModels, "\n - "),
				strings.Join(chatmodels.ImageModels, "\n - ")),
			false,
		)
		return
	}
}
