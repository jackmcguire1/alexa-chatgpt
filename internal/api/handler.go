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

type intentHandler func(ctx context.Context, req alexa.Request, xrayID string) (alexa.Response, error)

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
	AnimalGame      *AnimalGame
	LastIntent      alexa.Request
	SystemMessage   string
	routes          map[string]intentHandler
}

func NewHandler(
	logger *slog.Logger,
	chatGptService chatmodels.Service,
	responsesQueue queue.PullPoll,
	requestsQueue queue.PullPoll,
	pollDelay int,
	model chatmodels.ChatModel,
	imageModel chatmodels.ImageModel,
	randomNumberSvc *RandomNumberGame,
	battleShips *Battleships,
	animalGame *AnimalGame,
) *Handler {
	h := &Handler{
		Logger:          logger,
		ChatGptService:  chatGptService,
		ResponsesQueue:  responsesQueue,
		RequestsQueue:   requestsQueue,
		PollDelay:       pollDelay,
		Model:           model,
		ImageModel:      imageModel,
		RandomNumberSvc: randomNumberSvc,
		BattleShips:     battleShips,
		AnimalGame:      animalGame,
	}
	h.initRoutes()
	return h
}

func (h *Handler) initRoutes() {
	h.routes = map[string]intentHandler{
		alexa.PurgeIntent:              h.handlePurge,
		alexa.ModelIntent:              h.handleModel,
		alexa.ImageIntent:              h.handleImage,
		alexa.SystemMessageIntent:      h.handleSystemMessage,
		alexa.SystemAutoCompleteIntent: h.handleSystemAutoComplete,
		alexa.TranslateIntent:          h.handleTranslate,
		alexa.AutoCompleteIntent:       h.handleAutoComplete,
		alexa.RandomFactIntent:         h.handleRandomFact,
		alexa.BattleshipStatusIntent:   h.handleBattleshipStatus,
		alexa.BattleShipsIntent:        h.handleBattleships,
		alexa.AnimalStatusIntent:       h.handleAnimalStatus,
		alexa.AnimalHintIntent:         h.handleAnimalHint,
		alexa.AnimalGuessIntent:        h.handleAnimalGuess,
		alexa.RandomNumberIntent:       h.handleRandomNumber,
		alexa.LastResponseIntent:       h.handleLastResponse,
		alexa.HelpIntent:               h.handleHelp,
		alexa.CancelIntent:             h.handleCancel,
		alexa.NoIntent:                 h.handleStop,
		alexa.StopIntent:               h.handleStop,
		alexa.FallbackIntent:           h.handleFallback,
	}
}

func (h *Handler) randomFact(ctx context.Context) (string, error) {
	ctx, span := trace.Start(ctx, "randomFact")
	defer span.End()

	return h.ChatGptService.TextGeneration(ctx, "tell me a random fact", h.Model)
}

func (h *Handler) DispatchIntents(ctx context.Context, req alexa.Request) (alexa.Response, error) {
	h.Logger.With("intent", utils.ToJSON(req)).Info("got intent")
	ctx, span := trace.Start(ctx, "DispatchIntents")
	defer span.End()
	xrayID := otelsetup.GetXRayTraceID(ctx)
	span.SetAttributes(attribute.String("intent", req.Body.Intent.Name))

	handler, ok := h.routes[req.Body.Intent.Name]
	if !ok {
		h.Logger.Error("user has invoked unsupported intent")
		return alexa.NewResponse("unsupported intent", "unsupported intent!", false), nil
	}
	return handler(ctx, req, xrayID)
}

func (h *Handler) handlePurge(ctx context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	err := h.ResponsesQueue.Purge(ctx)
	if err != nil {
		return alexa.Response{}, err
	}
	return alexa.NewResponse("Purged", "successfully purged queue", false), nil
}

func (h *Handler) handleModel(_ context.Context, req alexa.Request, _ string) (alexa.Response, error) {
	model := req.Body.Intent.Slots["chatModel"].Value
	return h.getOrSetModel(model)
}

func (h *Handler) handleImage(ctx context.Context, req alexa.Request, xrayID string) (alexa.Response, error) {
	prompt := req.Body.Intent.Slots["prompt"].Value
	h.Logger.With("prompt", prompt).Info("found phrase to autocomplete")

	err := h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{Prompt: prompt, ImageModel: &h.ImageModel, TraceID: xrayID})
	if err != nil {
		return alexa.Response{}, err
	}

	return h.GetResponse(ctx, h.PollDelay, false)
}

func (h *Handler) handleSystemMessage(_ context.Context, req alexa.Request, _ string) (alexa.Response, error) {
	prompt := req.Body.Intent.Slots["prompt"].Value
	h.Logger.With("prompt", prompt).Info("settings system message")
	h.SystemMessage = prompt
	return alexa.NewResponse("System Message Set", "Ok", false), nil
}

func (h *Handler) handleSystemAutoComplete(ctx context.Context, req alexa.Request, xrayID string) (alexa.Response, error) {
	prompt := req.Body.Intent.Slots["prompt"].Value
	h.Logger.With("prompt", prompt).Info("found phrase to autocomplete")

	err := h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{Prompt: prompt, Model: h.Model, SystemPrompt: h.SystemMessage, TraceID: xrayID})
	if err != nil {
		return alexa.Response{}, err
	}
	return h.GetResponse(ctx, h.PollDelay, false)
}

func (h *Handler) handleTranslate(ctx context.Context, req alexa.Request, xrayID string) (alexa.Response, error) {
	prompt := req.Body.Intent.Slots["prompt"].Value
	spaces := strings.Split(prompt, " ")
	sourceLanguage := spaces[0]
	targetLanguage := spaces[2]
	promptPhrase := strings.SplitN(prompt, targetLanguage, 2)
	promptToTranslate := promptPhrase[1]

	err := h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{
		Prompt:         promptToTranslate,
		TargetLanguage: targetLanguage,
		SourceLanguage: sourceLanguage,
		Model:          chatmodels.CHAT_MODEL_TRANSLATIONS,
		TraceID:        xrayID,
	})
	if err != nil {
		return alexa.Response{}, err
	}

	return h.GetResponse(ctx, h.PollDelay, false)
}

func (h *Handler) handleAutoComplete(ctx context.Context, req alexa.Request, xrayID string) (alexa.Response, error) {
	prompt := req.Body.Intent.Slots["prompt"].Value
	h.Logger.With("prompt", prompt).Info("found phrase to autocomplete")

	err := h.RequestsQueue.PushMessage(ctx, &chatmodels.Request{Prompt: prompt, Model: h.Model, TraceID: xrayID})
	if err != nil {
		return alexa.Response{}, err
	}

	return h.GetResponse(ctx, h.PollDelay, false)
}

func (h *Handler) handleRandomFact(ctx context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	h.Logger.Debug("random fact")

	execTime := time.Now().UTC()
	randomFact, err := h.randomFact(ctx)
	if err != nil {
		return alexa.Response{}, err
	}
	h.lastResponse = &chatmodels.LastResponse{Response: randomFact, TimeDiff: fmt.Sprintf("%.0f", time.Since(execTime).Seconds()), Model: h.Model.String()}
	return alexa.NewResponse("Random Fact", randomFact, false), nil
}

func (h *Handler) handleBattleshipStatus(ctx context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	alive, killed := h.BattleShips.ShipsTotals()
	hits, misses := h.BattleShips.TotalHitsAndMisses()

	statusStr := "the user is playing a game of battleships, tell the status update of their game, ther are %d boats still alive, %d boats have been killed. Their total hits are %d, their total misses are %d."
	statement, _ := h.ChatGptService.TextGeneration(ctx, fmt.Sprintf(statusStr, alive, killed, hits, misses), h.Model)
	return alexa.NewResponse("BattleShips", statement, false), nil
}

func (h *Handler) handleBattleships(ctx context.Context, req alexa.Request, _ string) (alexa.Response, error) {
	execTime := time.Now().UTC()

	x, ok := req.Body.Intent.Slots["x"]
	if !ok || x.Value == "" || x.Value == "?" {
		return alexa.NewResponse("Try again!", "Try again!", false), nil
	}

	y, ok := req.Body.Intent.Slots["y"]
	if !ok || y.Value == "" || y.Value == "?" {
		return alexa.NewResponse("Try again!", "Try again!", false), nil
	}

	x_cord, _ := strconv.Atoi(x.Value)
	y_cord, _ := strconv.Atoi(y.Value)

	var statement string
	switch h.BattleShips.Attack(x_cord, y_cord) {
	case Hit:
		statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they hit a ship", h.Model)
	case Miss:
		statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they missed a ship", h.Model)
	case Sink:
		statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they sunk a ship", h.Model)
	case GameOver:
		statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they won the game", h.Model)
		h.BattleShips = NewBattleShipSetup()
	case Invalid:
		statement, _ = h.ChatGptService.TextGeneration(ctx, "playing battleships, tell the user they made an invalid move", h.Model)
	}
	h.lastResponse = &chatmodels.LastResponse{Response: statement, TimeDiff: fmt.Sprintf("%.0f", time.Since(execTime).Seconds()), Model: h.Model.String()}
	return alexa.NewResponse("BattleShips", statement, false), nil
}

func (h *Handler) handleAnimalStatus(ctx context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	status := h.AnimalGame.GetStatus()

	systemPrompt := "You are a friendly game host for children aged 7 and up. Be encouraging, enthusiastic, and use simple language."
	statusStr := "The player has %d guesses left and %d hints remaining in the animal guessing game. Tell them this information."
	statement, _ := h.ChatGptService.TextGenerationWithSystem(ctx, systemPrompt, fmt.Sprintf(statusStr, status.GuessesLeft, status.HintsLeft), h.Model)
	return alexa.NewResponse("Animal Game", statement, false), nil
}

func (h *Handler) handleAnimalHint(ctx context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	execTime := time.Now().UTC()

	hintResult := h.AnimalGame.RequestHint()
	var statement string
	systemPrompt := "You are a friendly game host for children aged 7 and up. Be encouraging, enthusiastic, and use simple language."

	switch hintResult.Status {
	case GameInactive:
		statement, _ = h.ChatGptService.TextGenerationWithSystem(ctx, systemPrompt, "Tell the player there's no active animal guessing game. They need to start a new game by making a guess.", h.Model)
	case NoHintsLeft:
		statement, _ = h.ChatGptService.TextGenerationWithSystem(ctx, systemPrompt, "Tell the player they have no hints left. They need to keep guessing!", h.Model)
	case HintAvailable:
		hintSystemPrompt := "You are a friendly game host helping children aged 7 and up guess animals. Give educational, age-appropriate hints without revealing the animal's name. Be fun and encouraging."
		hintPrompt := fmt.Sprintf("Give hint number %d about a %s. The child has %d guesses left.",
			hintResult.HintNumber, hintResult.Animal, hintResult.GuessesLeft)
		statement, _ = h.ChatGptService.TextGenerationWithSystem(ctx, hintSystemPrompt, hintPrompt, h.Model)
	}
	h.lastResponse = &chatmodels.LastResponse{Response: statement, TimeDiff: fmt.Sprintf("%.0f", time.Since(execTime).Seconds()), Model: h.Model.String()}
	return alexa.NewResponse("Animal Game", statement, false), nil
}

func (h *Handler) handleAnimalGuess(ctx context.Context, req alexa.Request, _ string) (alexa.Response, error) {
	execTime := time.Now().UTC()

	guessSlot, ok := req.Body.Intent.Slots["animal"]
	if !ok || guessSlot.Value == "" || guessSlot.Value == "?" {
		return alexa.NewResponse("Try again!", "Please say an animal name to make your guess!", false), nil
	}

	guess := guessSlot.Value
	result := h.AnimalGame.MakeGuess(guess)

	var statement string
	systemPrompt := "You are a friendly game host for children aged 7 and up. Be encouraging, enthusiastic, and use simple language."

	switch result.Status {
	case GameWon:
		congratsPrompt := fmt.Sprintf("The player correctly guessed the animal '%s'! Congratulate them enthusiastically. Then say: Here's what a %s sounds like: %s",
			result.Animal, result.Animal, result.Sound)
		statement, _ = h.ChatGptService.TextGenerationWithSystem(ctx, systemPrompt, congratsPrompt, h.Model)
		h.AnimalGame.ResetGame()
	case GuessIncorrect:
		incorrectPrompt := fmt.Sprintf("The player guessed '%s' but it's wrong. They have %d guesses left. Encourage them to try again and suggest they can ask for a hint!",
			guess, result.GuessesLeft)
		statement, _ = h.ChatGptService.TextGenerationWithSystem(ctx, systemPrompt, incorrectPrompt, h.Model)
	case GameLost:
		losePrompt := fmt.Sprintf("The player ran out of guesses! The correct animal was '%s'. Be supportive and encourage them to play again.",
			result.Animal)
		statement, _ = h.ChatGptService.TextGenerationWithSystem(ctx, systemPrompt, losePrompt, h.Model)
		h.AnimalGame.ResetGame()
	case GameInactive:
		statement = "There's no active game! Say an animal name to start a new guessing game!"
		h.AnimalGame.ResetGame()
	}
	h.lastResponse = &chatmodels.LastResponse{Response: statement, TimeDiff: fmt.Sprintf("%.0f", time.Since(execTime).Seconds()), Model: h.Model.String()}
	return alexa.NewResponse("Animal Game", statement, false), nil
}

func (h *Handler) handleRandomNumber(ctx context.Context, req alexa.Request, _ string) (alexa.Response, error) {
	if req.Body.Intent.Slots["number"].Value == "cheat" {
		res := alexa.NewResponse("Random Fact", fmt.Sprintf("You gave up so easily, your number is %d", h.RandomNumberSvc.Number), false)
		h.RandomNumberSvc.ShuffleRandomNumber()
		return res, nil
	}
	return h.RandomNumberGame(ctx, req)
}

func (h *Handler) handleLastResponse(ctx context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	h.Logger.Debug("fetching last response")
	return h.GetResponse(ctx, h.PollDelay, true)
}

func (h *Handler) handleHelp(_ context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	return alexa.NewResponse(
		"Help",
		"simply repeat, question followed by a desired sentence, to change model simply say 'use' followed by 'gpt' or 'gemini'",
		false,
	), nil
}

func (h *Handler) handleCancel(_ context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	h.Logger.Debug("user has invoked cancelled intent")
	return alexa.NewResponse("Next Question", "okay, i'm listening", false), nil
}

func (h *Handler) handleStop(_ context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	h.Logger.Debug("user has invoked no or stop intent")
	return alexa.NewResponse("Bye", "Good bye", true), nil
}

func (h *Handler) handleFallback(_ context.Context, _ alexa.Request, _ string) (alexa.Response, error) {
	h.Logger.Debug("user has invoked fallback intent")
	return alexa.NewResponse("Try again!", "Try again!", false), nil
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
