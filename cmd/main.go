package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jackmcguire1/alexa-chatgpt/internal/api"
	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
)

func main() {
	svc := chatgpt.NewClient(&chatgpt.Resources{
		Api: chatgpt.NewChatGptClient(os.Getenv("OPENAI_API_KEY")),
	})

	h := api.Handler{ChatGptService: svc}
	lambda.Start(h.Invoke)
}
