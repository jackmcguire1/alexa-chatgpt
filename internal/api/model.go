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
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_GEMINI) {
			res = alexa.NewResponse("Chat Models", "Gemini model is not available - Gemini API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_GEMINI
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_GPT.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_GPT) {
			res = alexa.NewResponse("Chat Models", "GPT model is not available - OpenAI API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_GPT
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_GPT_V4.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_GPT_V4) {
			res = alexa.NewResponse("Chat Models", "GPT-4 model is not available - OpenAI API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_GPT_V4
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_META.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_META) {
			res = alexa.NewResponse("Chat Models", "Meta model is not available - Cloudflare API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_META
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_SQL.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_SQL) {
			res = alexa.NewResponse("Chat Models", "SQL model is not available - Cloudflare API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_SQL
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_OPEN.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_OPEN) {
			res = alexa.NewResponse("Chat Models", "Open model is not available - Cloudflare API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_OPEN
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_AWQ.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_AWQ) {
			res = alexa.NewResponse("Chat Models", "AWQ model is not available - Cloudflare API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_AWQ
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_QWEN.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_QWEN) {
			res = alexa.NewResponse("Chat Models", "Qwen model is not available - Cloudflare API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_QWEN
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_OPUS.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_OPUS) {
			res = alexa.NewResponse("Chat Models", "Opus model is not available - Anthropic API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_OPUS
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_SONNET.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_SONNET) {
			res = alexa.NewResponse("Chat Models", "Sonnet model is not available - Anthropic API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_SONNET
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.CHAT_MODEL_GPT_OSS.String():
		if !chatmodels.IsModelAvailable(chatmodels.CHAT_MODEL_GPT_OSS) {
			res = alexa.NewResponse("Chat Models", "GPT-OSS model is not available - OpenAI API key not configured", false)
			return
		}
		h.Model = chatmodels.CHAT_MODEL_GPT_OSS
		res = alexa.NewResponse("Chat Models", "ok", false)
		return
	case chatmodels.IMAGE_MODEL_STABLE_DIFFUSION.String():
		if !chatmodels.IsImageModelAvailable(chatmodels.IMAGE_MODEL_STABLE_DIFFUSION) {
			res = alexa.NewResponse("Image Models", "Stable Diffusion model is not available - Cloudflare API key not configured", false)
			return
		}
		h.ImageModel = chatmodels.IMAGE_MODEL_STABLE_DIFFUSION
		res = alexa.NewResponse("Image Models", "ok", false)
		return
	case chatmodels.IMAGE_MODEL_DALL_E_3.String():
		if !chatmodels.IsImageModelAvailable(chatmodels.IMAGE_MODEL_DALL_E_3) {
			res = alexa.NewResponse("Image Models", "DALL-E 3 model is not available - OpenAI API key not configured", false)
			return
		}
		h.ImageModel = chatmodels.IMAGE_MODEL_DALL_E_3
		res = alexa.NewResponse("Image Models", "ok", false)
		return
	case chatmodels.IMAGE_MODEL_DALL_E_2.String():
		if !chatmodels.IsImageModelAvailable(chatmodels.IMAGE_MODEL_DALL_E_2) {
			res = alexa.NewResponse("Image Models", "DALL-E 2 model is not available - OpenAI API key not configured", false)
			return
		}
		h.ImageModel = chatmodels.IMAGE_MODEL_DALL_E_2
		res = alexa.NewResponse("Image Models", "ok", false)
		return
	case chatmodels.IMAGE_MODEL_GEMINI.String():
		if !chatmodels.IsImageModelAvailable(chatmodels.IMAGE_MODEL_GEMINI) {
			res = alexa.NewResponse("Image Models", "Gemini image model is not available - Gemini API key not configured", false)
			return
		}
		h.ImageModel = chatmodels.IMAGE_MODEL_GEMINI
		res = alexa.NewResponse("Image Models", "ok", false)
		return
	case "which":
		res = alexa.NewResponse("Chat Models", fmt.Sprintf("I am using the text-model %s and image-model %s", h.Model.String(), h.ImageModel.String()), false)
		return
	default:
		res = alexa.NewResponse(
			"Models",
			fmt.Sprintf("The available chat models are \n - %s and image models %s", strings.Join(chatmodels.AvailableModels, "\n - "), strings.Join(chatmodels.ImageModels, "\n - ")),
			false,
		)
		return
	}
}
