package chatmodels

type ChatModel string

type ImageModel string

const (
	CHAT_MODEL_GEMINI       ChatModel = "gemini"
	CHAT_MODEL_GPT          ChatModel = "gpt"
	CHAT_MODEL_META         ChatModel = "llama"
	CHAT_MODEL_AWQ          ChatModel = "awq"
	CHAT_MODEL_TRANSLATIONS ChatModel = "translate"
	CHAT_MODEL_OPEN         ChatModel = "open chat"
	CHAT_MODEL_SQL          ChatModel = "sql"
	CHAT_MODEL_QWEN         ChatModel = "qwen"
)

const (
	IMAGE_MODEL_STABLE_DIFFUSION ImageModel = "stable"
	IMAGE_MODEL_DALL_E_2         ImageModel = "dallas v2"
	IMAGE_MODEL_DALL_E_3         ImageModel = "dallas"
)

var AvaliableModels = []string{
	CHAT_MODEL_GPT.String(),
	CHAT_MODEL_GEMINI.String(),
	CHAT_MODEL_META.String(),
	CHAT_MODEL_SQL.String(),
	CHAT_MODEL_OPEN.String(),
	CHAT_MODEL_AWQ.String(),
	CHAT_MODEL_QWEN.String(),
}

var ImageModels = []string{
	IMAGE_MODEL_STABLE_DIFFUSION.String(),
	IMAGE_MODEL_DALL_E_2.String(),
}

var StrToImageModel = map[string]ImageModel{
	IMAGE_MODEL_STABLE_DIFFUSION.String(): IMAGE_MODEL_STABLE_DIFFUSION,
	IMAGE_MODEL_DALL_E_2.String():         IMAGE_MODEL_DALL_E_2,
}

func (c ChatModel) String() string {
	return string(c)
}

func (c ImageModel) String() string {
	return string(c)
}
