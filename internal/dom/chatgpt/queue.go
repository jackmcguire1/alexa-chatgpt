package chatgpt

import "time"

type LastResponse struct {
	Prompt   string
	Response string
	TimeDiff time.Duration
}

type Request struct {
	Prompt string
}
