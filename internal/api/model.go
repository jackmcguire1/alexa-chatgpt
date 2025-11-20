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
		// Build formatted list with alias -> provider model mapping
		chatModelsMap := chatmodels.GetAvailableChatModelsWithProviderIDs()
		imageModelsMap := chatmodels.GetAvailableImageModelsWithProviderIDs()

		var chatModelsList []string
		for alias, providerModel := range chatModelsMap {
			chatModelsList = append(chatModelsList, fmt.Sprintf("%s (%s)", alias, providerModel))
		}

		var imageModelsList []string
		for alias, providerModel := range imageModelsMap {
			imageModelsList = append(imageModelsList, fmt.Sprintf("%s (%s)", alias, providerModel))
		}

		res = alexa.NewResponse(
			responseTitleModels,
			fmt.Sprintf("The available models are, TEXT MODELS: \n - %s\n\nIMAGE MODELS: \n - %s",
				strings.Join(chatModelsList, "\n - "),
				strings.Join(imageModelsList, "\n - ")),
			false,
		)
		return
	}
}
