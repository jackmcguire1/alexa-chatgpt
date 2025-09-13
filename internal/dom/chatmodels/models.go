package chatmodels

import (
	"errors"
)

var MissingContentError = errors.New("Missing content")

var (
	hasOpenAI     bool
	hasGemini     bool
	hasAnthropic  bool
	hasCloudflare bool
)

type ChatModel string

type ImageModel string

const (
	CHAT_MODEL_GEMINI       ChatModel = "gemini"
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
	IMAGE_MODEL_STABLE_DIFFUSION ImageModel = "stable"
	IMAGE_MODEL_DALL_E_2         ImageModel = "dallas v2"
	IMAGE_MODEL_DALL_E_3         ImageModel = "dallas"
	IMAGE_MODEL_GEMINI           ImageModel = "gemini image"
)

var AvailableModels []string

var ImageModels []string

func RegisterAvailableClients(openAI, gemini, anthropic, cloudflare bool) {
	hasOpenAI = openAI
	hasGemini = gemini
	hasAnthropic = anthropic
	hasCloudflare = cloudflare

	AvailableModels = []string{}
	ImageModels = []string{}

	if hasOpenAI {
		AvailableModels = append(AvailableModels,
			CHAT_MODEL_GPT.String(),
			CHAT_MODEL_GPT_V4.String(),
		)
		ImageModels = append(ImageModels,
			IMAGE_MODEL_DALL_E_3.String(),
			IMAGE_MODEL_DALL_E_2.String(),
		)
	}

	if hasGemini {
		AvailableModels = append(AvailableModels, CHAT_MODEL_GEMINI.String())
		ImageModels = append(ImageModels, IMAGE_MODEL_GEMINI.String())
	}

	if hasAnthropic {
		AvailableModels = append(AvailableModels,
			CHAT_MODEL_OPUS.String(),
			CHAT_MODEL_SONNET.String(),
		)
	}

	if hasCloudflare {
		AvailableModels = append(AvailableModels,
			CHAT_MODEL_META.String(),
			CHAT_MODEL_QWEN.String(),
			CHAT_MODEL_GPT_OSS.String(),
		)
		ImageModels = append(ImageModels, IMAGE_MODEL_STABLE_DIFFUSION.String())
	}
}

func IsModelAvailable(model ChatModel) bool {
	switch model {
	case CHAT_MODEL_GPT, CHAT_MODEL_GPT_V4:
		return hasOpenAI
	case CHAT_MODEL_GEMINI:
		return hasGemini
	case CHAT_MODEL_OPUS, CHAT_MODEL_SONNET:
		return hasAnthropic
	case CHAT_MODEL_META, CHAT_MODEL_QWEN, CHAT_MODEL_TRANSLATIONS, CHAT_MODEL_GPT_OSS:
		return hasCloudflare
	default:
		return false
	}
}

func IsImageModelAvailable(model ImageModel) bool {
	switch model {
	case IMAGE_MODEL_DALL_E_2, IMAGE_MODEL_DALL_E_3:
		return hasOpenAI
	case IMAGE_MODEL_GEMINI:
		return hasGemini
	case IMAGE_MODEL_STABLE_DIFFUSION:
		return hasCloudflare
	default:
		return false
	}
}

var StrToImageModel = map[string]ImageModel{
	IMAGE_MODEL_STABLE_DIFFUSION.String(): IMAGE_MODEL_STABLE_DIFFUSION,
	IMAGE_MODEL_DALL_E_2.String():         IMAGE_MODEL_DALL_E_2,
	IMAGE_MODEL_DALL_E_3.String():         IMAGE_MODEL_DALL_E_3,
	IMAGE_MODEL_GEMINI.String():           IMAGE_MODEL_GEMINI,
}

func (c ChatModel) String() string {
	return string(c)
}

func (c ImageModel) String() string {
	return string(c)
}
