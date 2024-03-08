package chatmodels

type ChatModel string

const (
	CHAT_MODEL_GEMINI ChatModel = "gemini"
	CHAT_MODEL_GPT    ChatModel = "gpt"
	CHAT_MODEL_META   ChatModel = "meta"
	CHAT_MODEL_AWQ    ChatModel = "awsq"
	CHAT_MODEL_OPEN   ChatModel = "open chat"
	CHAT_MODEL_SQL    ChatModel = "sql"
)

var AvaliableModels = []string{
	CHAT_MODEL_GPT.String(),
	CHAT_MODEL_GEMINI.String(),
	CHAT_MODEL_META.String(),
	CHAT_MODEL_SQL.String(),
	CHAT_MODEL_OPEN.String(),
	CHAT_MODEL_AWQ.String(),
}

func (c ChatModel) String() string {
	return string(c)
}
