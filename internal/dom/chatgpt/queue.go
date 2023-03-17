package chatgpt

type LastResponse struct {
	Prompt   string `json:"prompt"`
	Response string `json:"response"`
	TimeDiff string `json:"time_diff"`
}

type Request struct {
	Prompt string `json:"prompt"`
}
