package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatgpt"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

type Handler struct {
	ChatGptService chatgpt.Service
	lastResponse   *chatgpt.LastResponse
	ResponsesQueue queue.PullPoll
	RequestsQueue  queue.PullPoll
	PollDelay      int
}

func (h *Handler) randomFact(ctx context.Context) (string, error) {
	return h.ChatGptService.AutoComplete(ctx, "tell me a random fact")
}

func (h *Handler) DispatchIntents(ctx context.Context, req alexa.Request) (res alexa.Response, err error) {
	switch req.Body.Intent.Name {
	case alexa.AutoCompleteIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		log.Println("found phrase to autocomplete", prompt)

		err = h.RequestsQueue.PushMessage(ctx, &chatgpt.Request{Prompt: prompt})
		if err != nil {
			break
		}

		var data []byte
		data, err = h.ResponsesQueue.PullMessage(ctx, h.PollDelay)
		if err != nil && !errors.Is(err, queue.EmptyMessageErr) {
			break
		}

		if len(data) > 0 {
			var response *chatgpt.LastResponse
			err = json.Unmarshal(data, &response)
			if err != nil {
				log.Println("failed to unmarshal response", string(data), err)
				return
			}
			res = alexa.NewResponse("Autocomplete", response.Response, false)
			h.lastResponse = response

			return
		}
		res = alexa.NewResponse("Autocomplete", "your response will be available momentarily", false)
	case alexa.RandomFactIntent:
		var randomFact string

		execTime := time.Now().UTC()

		randomFact, err = h.randomFact(ctx)
		if err != nil {
			return
		}
		res = alexa.NewResponse("Random Fact", randomFact, false)
		h.lastResponse = &chatgpt.LastResponse{Response: randomFact, TimeDiff: time.Since(execTime).String()}

	case alexa.LastResponseIntent:
		log.Println("fetching last response")

		var data []byte
		data, err = h.ResponsesQueue.PullMessage(ctx, h.PollDelay)
		if err != nil && !errors.Is(err, queue.EmptyMessageErr) {
			break
		}

		if len(data) > 0 {
			var response *chatgpt.LastResponse
			err = json.Unmarshal(data, &response)
			if err != nil {
				log.Println("failed to unmarshal response", string(data), err)
				return
			}
			res = alexa.NewResponse("Last Response",
				fmt.Sprintf(
					"%s, this took %s to fetch the answer",
					response.Response,
					response.TimeDiff,
				),
				false,
			)
			h.lastResponse = response
			return
		}

		if h.lastResponse != nil {
			res = alexa.NewResponse("Last Response",
				fmt.Sprintf(
					"%s, this took %s to fetch the answer",
					h.lastResponse.Response,
					h.lastResponse.TimeDiff,
				),
				false,
			)
			return
		}

		res = alexa.NewResponse("Last Response", "I do not have a answer to your last prompt", false)

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
			"Hi, lets begin our conversation!",
			false,
		)
	default:
		resp, err = h.DispatchIntents(ctx, req)
	}

	if err != nil {
		resp = alexa.NewResponse("error",
			"an error occurred when processing your prompt",
			true,
		)
	}

	log.Println(utils.ToJSON(resp))
	return
}
