package chatmodels

type ChatModel string

const (
	CHAT_MODEL_GEMINI ChatModel = "gemini"
	CHAT_MODEL_GPT    ChatModel = "gpt"
)

func (c ChatModel) String() string {
	return string(c)
}
