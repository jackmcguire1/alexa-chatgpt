package chatmodels

type LastResponse struct {
	Prompt         string    `json:"prompt"`
	Response       string    `json:"response"`
	TimeDiff       string    `json:"time_diff"`
	Model          ChatModel `json:"model"`
	ImagesResponse []string  `json:"images_responses"`
}

type Request struct {
	Prompt string    `json:"prompt"`
	Model  ChatModel `json:"model"`
}
