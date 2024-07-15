package chatmodels

type ChatModel string

const (
	CHAT_MODEL_GEMINI           ChatModel = "gemini"
	CHAT_MODEL_GPT              ChatModel = "gpt"
	CHAT_MODEL_META             ChatModel = "llama"
	CHAT_MODEL_AWQ              ChatModel = "awq"
	CHAT_MODEL_TRANSLATIONS     ChatModel = "translate"
	CHAT_MODEL_OPEN             ChatModel = "open chat"
	CHAT_MODEL_SQL              ChatModel = "sql"
	CHAT_MODEL_STABLE_DIFFUSION ChatModel = "stable"
	CHAT_MODEL_QWEN             ChatModel = "qwen"
)

var AvaliableModels = []string{
	CHAT_MODEL_GPT.String(),
	CHAT_MODEL_GEMINI.String(),
	CHAT_MODEL_META.String(),
	CHAT_MODEL_SQL.String(),
	CHAT_MODEL_OPEN.String(),
	CHAT_MODEL_AWQ.String(),
	CHAT_MODEL_STABLE_DIFFUSION.String(),
	CHAT_MODEL_QWEN.String(),
}

func (c ChatModel) String() string {
	return string(c)
}
