package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
)

func (h *Handler) GetResponse(ctx context.Context, delay int, lastResponse bool) (res alexa.Response, err error) {
	var data []byte
	data, err = h.ResponsesQueue.PullMessage(ctx, delay)
	if err != nil && !errors.Is(err, queue.EmptyMessageErr) {
		return
	}
	var response *chatmodels.LastResponse
	if len(data) == 0 && !lastResponse {
		res = alexa.NewResponse("Response", "your response will be available shortly", false)
		return
	}

	if len(data) == 0 && lastResponse {
		if h.lastResponse == nil {
			res = alexa.NewResponse("Response", "I do not have a answer to your last prompt", false)
			return
		}
		response = h.lastResponse
		goto response
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		h.Logger.
			With("error", err).
			With("data", string(data)).
			Error("failed to unmarshal chat model response")
		return
	}

response:
	if response.Error != "" {
		res = alexa.NewResponse(
			"Response",
			fmt.Sprintf("I encountered an error processing your prompt, %s", response.Error),
			false,
		)
		h.lastResponse = response
		return
	}

	switch response.Model {
	case chatmodels.IMAGE_MODEL_STABLE_DIFFUSION.String():
		res = alexa.NewImageResponse(
			"Response",
			fmt.Sprintf("your generated image took %s seconds to fetch", response.TimeDiff),
			response.ImagesResponse[0],
			response.ImagesResponse[1],
			false,
		)
		h.lastResponse = response
		return
	case chatmodels.CHAT_MODEL_TRANSLATIONS.String():
		res = alexa.NewResponse(
			"Response",
			fmt.Sprintf(
				"your translated prompt is %s, this took %s seconds to fetch the answer",
				response.Response,
				response.TimeDiff,
			),
			false,
		)
		h.lastResponse = response
		return
	default:
		res = alexa.NewResponse("Response",
			fmt.Sprintf(
				"%s, from the %s model, this took %s seconds to fetch the answer",
				response.Response,
				response.Model,
				response.TimeDiff,
			),
			false,
		)
		h.lastResponse = response
	}

	return

}
