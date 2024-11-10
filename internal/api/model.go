package api

import (
	"fmt"
	"strings"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
)

func (h *Handler) getOrSetModel(model string) (res alexa.Response, err error) {
	switch strings.ToLower(model) {
	case chatmodels.CHAT_MODEL_GEMINI.String():
		h.Model = chatmodels.CHAT_MODEL_GEMINI
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_GPT.String():
		h.Model = chatmodels.CHAT_MODEL_GPT
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_META.String():
		h.Model = chatmodels.CHAT_MODEL_META
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_SQL.String():
		h.Model = chatmodels.CHAT_MODEL_SQL
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_OPEN.String():
		h.Model = chatmodels.CHAT_MODEL_OPEN
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_AWQ.String():
		h.Model = chatmodels.CHAT_MODEL_AWQ
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_QWEN.String():
		h.Model = chatmodels.CHAT_MODEL_QWEN
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.IMAGE_MODEL_STABLE_DIFFUSION.String():
		h.ImageModel = chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
		res = alexa.NewResponse("Image Models", "ok", false)
		return
	case chatmodels.IMAGE_MODEL_DALL_E_3.String():
		h.ImageModel = chatmodels.IMAGE_MODEL_DALL_E_3
		res = alexa.NewResponse("Image Models", "ok", false)
		return
	case chatmodels.IMAGE_MODEL_DALL_E_2.String():
		h.ImageModel = chatmodels.IMAGE_MODEL_DALL_E_2
		res = alexa.NewResponse("Image Models", "ok", false)
		return
	case "which":
		res = alexa.NewResponse("Chat Models", fmt.Sprintf("I am using the text-model %s and image-model %s", h.Model.String(), h.ImageModel.String()), false)
		return
	default:
		res = alexa.NewResponse(
			"Models",
			fmt.Sprintf("The avaliable chat models are \n - %s and image models %s", strings.Join(chatmodels.AvaliableModels, "\n - "), strings.Join(chatmodels.ImageModels, "\n - ")),
			false,
		)
		return
	}
	return
}
