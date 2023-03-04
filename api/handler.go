package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackmcguire1/alexa-chatgpt/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/pkg/alexa"
	"github.com/jackmcguire1/alexa-chatgpt/pkg/utils"
)

type Handler struct {
	ChatGptService chatgpt.Service
}

func randomResponse() string {
	responses := []string{
		"did not compute, repeat did not compute, for futher assistance please ask for 'help'",
		"Alert! Our control point is being captured",
		"For The Horde!",
		"For The Alliance!",
		"Snake?! Snakeeeeee!",
	}
	rand.Seed(time.Now().Unix())
	i := rand.Intn(len(responses))
	if i == len(responses) {
		i = i - 1
	}

	return responses[i]
}

func (h *Handler) DispatchIntents(ctx context.Context, req alexa.Request) (res alexa.Response, err error) {
	switch req.Body.Intent.Name {
	case chatgpt.ChatGPTIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		log.Println("found name", prompt)

		var chatGptRes string
		chatGptRes, err = h.ChatGptService.GetPrompt(ctx, prompt)
		if err != nil {
			return
		}

		res = alexa.NewResponse("Prompt", chatGptRes, false)
	case alexa.HelpIntent:
		res = alexa.NewResponse(
			"Help",
			"Simply ask me a question",
			false,
		)
	case alexa.NoIntent, alexa.StopIntent:
		res = alexa.NewResponse(
			"Bye",
			"Good bye",
			true,
		)
	case alexa.CancelIntent:
		res = alexa.NewResponse(
			"Next Question",
			"okay, i'm listening",
			false,
		)
	default:
		res = alexa.NewResponse(
			"Random response",
			randomResponse(),
			true,
		)
	}

	return
}

func (h *Handler) Invoke(ctx context.Context, req alexa.Request) (resp alexa.Response, err error) {
	log.Println(utils.ToJSON(req))

	switch req.Body.Type {
	case alexa.LaunchRequestType:
		log.Println("launch request type")
		resp = alexa.NewResponse("chatGPT",
			"Hi, lets begin our convesation!",
			false,
		)
	case alexa.IntentRequestType:
		resp, err = h.DispatchIntents(ctx, req)
	default:
	}

	if err != nil {
		resp = alexa.NewResponse("error",
			fmt.Sprintf("there was an error, %s", err.Error()),
			false,
		)
	}

	log.Println(utils.ToJSON(resp))
	return
}
