package main

import (
	"context"
	"log"
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
)

func main() {
	svc := chatgpt.NewClient(&chatgpt.Resources{
		Api: chatgpt.NewOpenAIClient(os.Getenv("OPENAI_API_KEY")),
	})
	resp, err := svc.AutoComplete(context.Background(), "tell me a random fact")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
