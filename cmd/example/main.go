package main

import (
	"context"
	"log"
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
)

func main() {
	svc := chatmodels.NewClient(&chatmodels.Resources{
		GPTApi: chatmodels.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
	})
	resp, err := svc.AutoComplete(context.Background(), "tell me a random fact", chatmodels.CHAT_MODEL_GPT)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
