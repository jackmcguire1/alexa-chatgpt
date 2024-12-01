package api

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
)

type RandomNumberGame struct {
	Number   int
	MaxLimit int
}

func NewRandomNumberGame(maxLimit int) RandomNumberGame {
	r := RandomNumberGame{
		Number:   rand.Intn(maxLimit),
		MaxLimit: maxLimit,
	}
	r.ShuffleRandomNumber()
	return r
}

const prompt = "You are a person helping someone guess a number between 0 and %d. offer help if they're closer to the target, with an abrupt higher or lower, or maybe a short witty phrase to help them figure out how close they're, if they get it correct, do a congratulatory statement. their guess was %s"
const (
	CORRECT = "CORRECT"
)

func (svc RandomNumberGame) ShuffleRandomNumber() {
	svc.Number = rand.Intn(svc.MaxLimit)
}

func (h *Handler) RandomNumberGame(ctx context.Context, req alexa.Request) (res alexa.Response, err error) {
	h.Logger.With("intent", req.Body.Intent.Name).Debug("got intent")
	guess := req.Body.Intent.Slots["number"].Value
	guessInt, _ := strconv.Atoi(guess)
	number := h.RandomNumberSvc.Number

	h.Logger.With("guess", guess).With("current number number", number).Info("got guess")

	if guessInt < number || guessInt > number {
		statement, _ := h.ChatGptService.TextGeneration(ctx, fmt.Sprintf(prompt, h.RandomNumberSvc.MaxLimit, guess), h.Model)
		res = alexa.NewResponse("Random Number Game", statement, false)
	}
	if guessInt == number {
		winningStatement := fmt.Sprintf(prompt, h.RandomNumberSvc.Number, (CORRECT + " " + guess + ". please also mention the number will be changed for the users to guess again."))
		statement, _ := h.ChatGptService.TextGeneration(
			ctx,
			winningStatement,
			h.Model,
		)
		res = alexa.NewResponse("Random Number Game", statement, false)
		h.RandomNumberSvc.ShuffleRandomNumber()
	}

	if res.Body.OutputSpeech != nil && res.Body.OutputSpeech.Text == "" {
		res = alexa.NewResponse("Random Number Game", "I'm sorry, something went wrong", false)
	}

	return res, nil
}
