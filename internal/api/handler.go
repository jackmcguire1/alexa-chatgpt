package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
)

type Handler struct {
	Logger         *slog.Logger
	ChatGptService chatmodels.Service
	lastResponse   *chatmodels.LastResponse
	ResponsesQueue queue.PullPoll
	RequestsQueue  queue.PullPoll
	PollDelay      int
	Model          chatmodels.ChatModel
}

func (h *Handler) randomFact(ctx context.Context) (string, error) {
	return h.ChatGptService.AutoComplete(ctx, "tell me a random fact", h.Model)
}

func (h *Handler) DispatchIntents(ctx context.Context, req alexa.Request) (res alexa.Response, err error) {
	h.Logger.With("intent", req.Body.Intent.Name).Debug("got intent")
	switch req.Body.Intent.Name {
	case alexa.PurgeIntent:
		err = h.ResponsesQueue.Purge(ctx)
		if err != nil {
			return res, err
		}
		res = alexa.NewResponse(
			"Purged",
			"successfully purged queue",
			false,
		)
	case alexa.ModelIntent:
		model := req.Body.Intent.Slots["chatModel"].Value
		h.Logger.With("model", model).Info("found model to use")

		switch strings.ToLower(model) {
		case chatmodels.CHAT_MODEL_GEMINI.String():
			h.Model = chatmodels.CHAT_MODEL_GEMINI
			res = alexa.NewResponse("Autocomplete", "ok", false)
			return
		case chatmodels.CHAT_MODEL_GPT.String():
			h.Model = chatmodels.CHAT_MODEL_GPT
			res = alexa.NewResponse("Autocomplete", "ok", false)
			return
		default:
			res = alexa.NewResponse("Autocomplete", fmt.Sprintf("I am using the model %s", h.Model.String()), false)
			return
		}

	case alexa.AutoCompleteIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		h.Logger.With("prompt", prompt).Info("found phrase to autocomplete")

		err = h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{Prompt: prompt, Model: h.Model})
		if err != nil {
			break
		}

		var data []byte
		data, err = h.ResponsesQueue.PullMessage(ctx, h.PollDelay)
		if err != nil && !errors.Is(err, queue.EmptyMessageErr) {
			break
		}

		if len(data) > 0 {
			var response *chatmodels.LastResponse
			err = json.Unmarshal(data, &response)
			if err != nil {
				h.Logger.
					With("error", err).
					With("data", string(data)).
					Error("failed to unmarshal chat model response")
				return
			}
			res = alexa.NewResponse(
				"Autocomplete",
				fmt.Sprintf("%s from the %s model", response.Response, response.Model),
				false,
			)
			h.lastResponse = response

			return
		}
		res = alexa.NewResponse("Autocomplete", "your response will be available shortly", false)
	case alexa.RandomFactIntent:
		h.Logger.Debug("random fact")
		var randomFact string

		execTime := time.Now().UTC()

		randomFact, err = h.randomFact(ctx)
		if err != nil {
			return
		}
		res = alexa.NewResponse("Random Fact", randomFact, false)
		h.lastResponse = &chatmodels.LastResponse{Response: randomFact, TimeDiff: fmt.Sprintf("%.0f", time.Since(execTime).Seconds()), Model: h.Model}

	case alexa.LastResponseIntent:
		h.Logger.Debug("fetching last response")

		var data []byte
		data, err = h.ResponsesQueue.PullMessage(ctx, h.PollDelay)
		if err != nil && !errors.Is(err, queue.EmptyMessageErr) {
			break
		}

		if len(data) > 0 {
			var response *chatmodels.LastResponse
			err = json.Unmarshal(data, &response)
			if err != nil {
				h.Logger.With("error", err).Error("failed to unmarshal response")
				return
			}
			res = alexa.NewResponse("Last Response",
				fmt.Sprintf(
					"%s, from the %s model, this took %s seconds to fetch the answer",
					response.Response,
					response.Model,
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
					"%s, from the %s model, this took %s seconds to fetch the answer",
					h.lastResponse.Response,
					h.lastResponse.Model,
					h.lastResponse.TimeDiff,
				),
				false,
			)
			return
		}

		h.Logger.Debug("no cached or polled last response")
		res = alexa.NewResponse("Last Response", "I do not have a answer to your last prompt", false)

	case alexa.HelpIntent:
		res = alexa.NewResponse(
			"Help",
			"simply repeat, question followed by a desired sentence, to change model simply say 'use' followed by 'gpt' or 'gemini'",
			false,
		)
	case alexa.CancelIntent:
		h.Logger.Debug("user has invoked cancelled intent")
		res = alexa.NewResponse(
			"Next Question",
			"okay, i'm listening",
			false,
		)
	case alexa.NoIntent, alexa.StopIntent:
		h.Logger.Debug("user has invoked no or stop intent")
		res = alexa.NewResponse(
			"Bye",
			"Good bye",
			true,
		)
	case alexa.FallbackIntent:
		h.Logger.Debug("user has invoked fallback intent")
		res = alexa.NewResponse(
			"Try again!",
			"Try again!",
			false,
		)
	default:
		h.Logger.Error("user has invoked unsupported intent")
		res = alexa.NewResponse(
			"unsupported intent",
			"unsupported intent!",
			false,
		)
	}

	return
}

func (h *Handler) Invoke(ctx context.Context, req alexa.Request) (resp alexa.Response, err error) {
	h.Logger.
		With("payload", utils.ToJSON(req)).
		Debug("lambda invoked")

	switch req.Body.Type {
	case alexa.LaunchRequestType:
		h.Logger.Debug("launch request type found")
		resp = alexa.NewResponse("chatGPT",
			"Hi, lets begin our conversation!",
			false,
		)
	default:
		resp, err = h.DispatchIntents(ctx, req)
	}

	if err != nil {
		h.Logger.
			With("error", err).
			Error("there as an error processing the request")

		err = nil
		resp = alexa.NewResponse("error",
			"an error occurred when processing your prompt",
			true,
		)
	}

	h.Logger.
		With("response", utils.ToJSON(resp)).
		Debug("returning response to alexa")

	return
}
