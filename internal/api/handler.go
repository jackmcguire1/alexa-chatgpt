package api

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jackmcguire1/alexa-chatgpt/internal/dom/chatmodels"
	otelsetup "github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/queue"
	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var trace = otel.Tracer("prompt-requester")

type Handler struct {
	Logger          *slog.Logger
	ChatGptService  chatmodels.Service
	lastResponse    *chatmodels.LastResponse
	ResponsesQueue  queue.PullPoll
	RequestsQueue   queue.PullPoll
	PollDelay       int
	Model           chatmodels.ChatModel
	ImageModel      chatmodels.ImageModel
	RandomNumberSvc *RandomNumberGame
	BattleShips     *Battleships
	LastIntent      alexa.Request
	SystemMessage   string
}

func (h *Handler) randomFact(ctx context.Context) (string, error) {
	ctx, span := trace.Start(ctx, "randomFact")
	defer span.End()

	return h.ChatGptService.TextGeneration(ctx, "tell me a random fact", h.Model)
}

func (h *Handler) DispatchIntents(ctx context.Context, req alexa.Request) (res alexa.Response, err error) {
	h.Logger.With("intent", utils.ToJSON(req)).Info("got intent")
	ctx, span := trace.Start(ctx, "DispatchIntents")
	defer span.End()
	xrayID := otelsetup.GetXRayTraceID(ctx)
	span.SetAttributes(attribute.String("intent", req.Body.Intent.Name))

	switch req.Body.Intent.Name {
	case alexa.PurgeIntent:
		err = h.ResponsesQueue.Purge(ctx)
		if err != nil {
			break
		}
		res = alexa.NewResponse(
			"Purged",
			"successfully purged queue",
			false,
		)
	case alexa.ModelIntent:
		model := req.Body.Intent.Slots["chatModel"].Value
		res, err = h.getOrSetModel(model)

	case alexa.ImageIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		h.Logger.With("prompt", prompt).Info("found phrase to autocomplete")

		err = h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{Prompt: prompt, ImageModel: &h.ImageModel, TraceID: xrayID})
		if err != nil {
			break
		}

		res, err = h.GetResponse(ctx, h.PollDelay, false)
	case alexa.SystemMessageIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		h.Logger.With("prompt", prompt).Info("settings system message")
		h.SystemMessage = prompt
		res = alexa.NewResponse(
			"System Message Set",
			"Ok",
			false,
		)
	case alexa.SystemAutoCompleteIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		h.Logger.With("prompt", prompt).Info("found phrase to autocomplete")

		err = h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{Prompt: prompt, Model: h.Model, SystemPrompt: h.SystemMessage, TraceID: xrayID})
		if err != nil {
			break
		}
		return h.GetResponse(ctx, h.PollDelay, false)
	case alexa.TranslateIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		spaces := strings.Split(prompt, " ")
		sourceLanguage := spaces[0]
		targetLanguage := spaces[2]
		promptPhrase := strings.SplitN(prompt, targetLanguage, 2)
		promptToTranslate := promptPhrase[1]

		err = h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{
			Prompt:         promptToTranslate,
			TargetLanguage: targetLanguage,
			SourceLanguage: sourceLanguage,
			Model:          chatmodels.CHAT_MODEL_TRANSLATIONS,
			TraceID:        xrayID,
		})
		if err != nil {
			break
		}

		res, err = h.GetResponse(ctx, h.PollDelay, false)
	case alexa.AutoCompleteIntent:
		prompt := req.Body.Intent.Slots["prompt"].Value
		h.Logger.With("prompt", prompt).Info("found phrase to autocomplete")

		err = h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{Prompt: prompt, Model: h.Model, TraceID: xrayID})
		if err != nil {
			break
		}

		return h.GetResponse(ctx, h.PollDelay, false)

	case alexa.RandomFactIntent:
		h.Logger.Debug("random fact")
		var randomFact string

		execTime := time.Now().UTC()

		randomFact, err = h.randomFact(ctx)
		if err != nil {
			return
		}
		res = alexa.NewResponse("Random Fact", randomFact, false)
		h.lastResponse = &chatmodels.LastResponse{Response: randomFact, TimeDiff: fmt.Sprintf("%.0f", time.Since(execTime).Seconds()), Model: h.Model.String()}
	case alexa.BattleshipStatusIntent:
		alive, killed := h.BattleShips.ShipsTotals()
		hits, misses := h.BattleShips.TotalHitsAndMisses()

		statusStr := "the user is playing a game of battleships, tell the status update of their game, ther are %d boats still alive, %d boats have been killed. Their total hits are %d, their total misses are %d."
		statement, _ := h.ChatGptService.TextGeneration(ctx, fmt.Sprintf(statusStr, alive, killed, hits, misses), h.Model)
		res = alexa.NewResponse("BattleShips", statement, false)

	case alexa.BattleShipsIntent:
		execTime := time.Now().UTC()

		x, ok := req.Body.Intent.Slots["x"]
		if !ok || x.Value == "" || x.Value == "?" {
			res = alexa.NewResponse(
				"Try again!",
				"Try again!",
				false,
			)
			return
		}

		y, ok := req.Body.Intent.Slots["y"]
		if !ok || y.Value == "" || y.Value == "?" {
			res = alexa.NewResponse(
				"Try again!",
				"Try again!",
				false,
			)
			return
		}

		x_cord, _ := strconv.Atoi(x.Value)
		y_cord, _ := strconv.Atoi(y.Value)

		var statement string
		switch h.BattleShips.Attack(x_cord, y_cord) {
		case Hit:
			statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they hit a ship", h.Model)
			res = alexa.NewResponse("BattleShips", statement, false)
		case Miss:
			statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they missed a ship", h.Model)
			res = alexa.NewResponse("BattleShips", statement, false)
		case Sink:
			statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they sunk a ship", h.Model)
			res = alexa.NewResponse("BattleShips", statement, false)
		case GameOver:
			statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they won the game", h.Model)
			res = alexa.NewResponse("BattleShips", statement, false)
			h.BattleShips = NewBattleShipSetup()
		case Invalid:
			statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they made an invalid move", h.Model)
			res = alexa.NewResponse("BattleShips", statement, false)
		}
		h.lastResponse = &chatmodels.LastResponse{Response: statement, TimeDiff: fmt.Sprintf("%.0f", time.Since(execTime).Seconds()), Model: h.Model.String()}
	case alexa.RandomNumberIntent:
		if req.Body.Intent.Slots["number"].Value == "cheat" {
			res = alexa.NewResponse("Random Fact", fmt.Sprintf("You gave up so easily, your number is %d", h.RandomNumberSvc.Number), false)
			h.RandomNumberSvc.ShuffleRandomNumber()
			return
		}
		return h.RandomNumberGame(ctx, req)
	case alexa.LastResponseIntent:
		h.Logger.Debug("fetching last response")

		return h.GetResponse(ctx, h.PollDelay, true)
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
	ctx, span := trace.Start(ctx, "Invoke")
	defer span.End()
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
	h.LastIntent = req

	if err != nil {
		span.RecordError(err)
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
