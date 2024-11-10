package main

import (
	"context"
	"log"
	"os"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
)

func main() {
	svc := chatmodels.NewClient(&chatmodels.Resources{
		GPTApi:              chatmodels.NewOpenAiApiClient(os.Getenv("OPENAI_API_KEY")),
		CloudflareApiClient: chatmodels.NewCloudflareApiClient("1c5e8bd244a9566794bcffc3cafe27fc", "Von5nrK3hhWQ-iarP1fQoi4_5624oIMR_Q7rfznP"),
	})
	resp, err := svc.AutoComplete(context.Background(), "monkey riding a unicorn", chatmodels.CHAT_MODEL_GPT)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(resp)
}
