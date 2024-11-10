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
		res = alexa.NewResponse("Chat Models", "ok", false)
	case chatmodels.IMAGE_MODEL_DALL_E_2.String():
		h.ImageModel = chatmodels.IMAGE_MODEL_DALL_E_2
		res = alexa.NewResponse("Chat Models", "ok", false)
	case "which":
		res = alexa.NewResponse("Chat Models", fmt.Sprintf("I am using the model %s", h.Model.String()), false)
		return
	default:
		res = alexa.NewResponse(
			"Chat Models",
			fmt.Sprintf("The avaliable chat models are \n - %s", strings.Join(chatmodels.AvaliableModels, "\n - ")),
			false,
		)
		return
	}
	return
}
