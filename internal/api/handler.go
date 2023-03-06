package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

type Handler struct {
	ChatGptService chatgpt.Service
}

func (h *Handler) randomResponse(ctx context.Context) (string, error) {
	responses := []string{
		"The cat sat on the",
		"The fox ran over the",
		"The pig rolled in",
		"Peppa pig jumped in a muddie puddle and",
		"ed, edd and eddy came up with a new hustle that involed",
		"Carter found a",
		"Chanel obtained a new rank in valorant called",
		"Robbie got a new rank in tekken called",
		"Perry wrote new lyrics about his",
		"Jack created a new application that solved",
		"George looted a new armor item called",
		"Nitish loved to play the computer game",
	}

	rand.Seed(time.Now().Unix())
	i := rand.Intn(len(responses))
	if i == len(responses) {
		i = i - 1
	}

	resp, err := h.ChatGptService.AutoComplete(ctx, responses[i])
	if err != nil {
		return "", err
	}
	phrase := fmt.Sprintf("%s %s", responses[i], resp)

	return phrase, nil
}

func (h *Handler) DispatchIntents(ctx context.Context, req alexa.Request) (res alexa.Response, err error) {
	switch req.Body.Intent.Name {
	case alexa.AutoCompleteIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		log.Println("found phrase to autocomplete", prompt)

		var chatGptRes string
		chatGptRes, err = h.ChatGptService.AutoComplete(ctx, prompt)
		if err != nil {
			break
		}

		res = alexa.NewResponse("Autocomplete", fmt.Sprintf("%s %s", prompt, chatGptRes), false)
	case alexa.HelpIntent:
		res = alexa.NewResponse(
			"Help",
			"Simply repeat, question followed by a desired sentence",
			false,
		)
	case alexa.CancelIntent:
		res = alexa.NewResponse(
			"Next Question",
			"okay, i'm listening",
			false,
		)
	case alexa.NoIntent, alexa.StopIntent:
		res = alexa.NewResponse(
			"Bye",
			"Good bye",
			true,
		)
	case alexa.FallbackIntent:
		res = alexa.NewResponse(
			"Try again!",
			"Try again!",
			false,
		)
	default:
		res = alexa.NewResponse(
			"unsupported intent",
			"unsupported intent!",
			false,
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
	default:
		resp, err = h.DispatchIntents(ctx, req)
	}

	if err != nil {
		resp = alexa.NewResponse("error",
			"an error occured when processing your prompt",
			true,
		)
	}

	log.Println(utils.ToJSON(resp))
	return
}
