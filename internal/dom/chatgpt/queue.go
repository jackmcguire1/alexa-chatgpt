package chatgpt

import "time"

type LastResponse struct {
	Prompt   string        `json:"prompt"`
	Response string        `json:"response"`
	TimeDiff time.Duration `json:"time_diff"`
}

type Request struct {
	Prompt string `json:"prompt"`
}
