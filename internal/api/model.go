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

func (h *Handler) setChatModel(config chatmodels.ModelConfig) alexa.Response {
	if !chatmodels.IsModelAvailable(config.ChatModel) {
		return alexa.NewResponse(responseTitleChatModels, config.ErrorMessage, false)
	}
	h.Model = config.ChatModel
	return alexa.NewResponse(responseTitleChatModels, responseOK, false)
}

func (h *Handler) setImageModel(config chatmodels.ModelConfig) alexa.Response {
	if !chatmodels.IsImageModelAvailable(config.ImageModel) {
		return alexa.NewResponse(responseTitleImageModels, config.ErrorMessage, false)
	}
	h.ImageModel = config.ImageModel
	return alexa.NewResponse(responseTitleImageModels, responseOK, false)
}

func (h *Handler) getOrSetModel(model string) (res alexa.Response, err error) {
	lowerModel := strings.ToLower(model)

	// Check chat models
	if config, ok := chatmodels.GetChatModelByAlias(lowerModel); ok {
		return h.setChatModel(config), nil
	}

	// Check image models
	if config, ok := chatmodels.GetImageModelByAlias(lowerModel); ok {
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
		chatModels := chatmodels.GetAvailableChatModels()
		imageModels := chatmodels.GetAvailableImageModels()
		res = alexa.NewResponse(
			responseTitleModels,
			fmt.Sprintf("The available models are, TEXT MODELS: \n - %s IMAGE MODELS %s",
				strings.Join(chatModels, "\n - "),
				strings.Join(imageModels, "\n - ")),
			false,
		)
		return
	}
}
