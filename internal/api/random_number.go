package api

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jackmcguire1/alexa-chatgpt/internal/pkg/alexa"
)

type RandomNumberGame struct {
	Number   int
	MaxLimit int
}

func NewRandomNumberGame(maxLimit int) *RandomNumberGame {
	r := &RandomNumberGame{
		Number:   rand.Intn(maxLimit),
		MaxLimit: maxLimit,
	}
	r.ShuffleRandomNumber()
	return r
}

const (
	winningprompt    = "please respond as if you were the person controlling the game, a person guessed a number between 0 and 100. they got the number correct, tell them a congratulatory statement. the winning number was %d. please also mention the target number will be changed for the next round"
	higherThanprompt = "please respond as if you were the person controlling the game, a person guessed a number between 0 and 100. the number was higher than target value, tell them advice to guess lower"
	lessThanprompt   = "please respond as if you were the person controlling the game, a person guessed a number between 0 and 100. the number was lower than the target value, tell them advice to guess higher"
)

func (svc *RandomNumberGame) ShuffleRandomNumber() {
	rand.Seed(time.Now().Unix())
	for {
		newNumber := rand.Intn(svc.MaxLimit)
		if newNumber != svc.Number {
			svc.Number = newNumber
			break
		}
	}
}

func (h *Handler) RandomNumberGame(ctx context.Context, req alexa.Request) (res alexa.Response, err error) {
	h.Logger.With("intent", req.Body.Intent.Name).Debug("got intent")
	guess := req.Body.Intent.Slots["number"].Value
	guessInt, _ := strconv.Atoi(guess)
	number := h.RandomNumberSvc.Number

	h.Logger.With("guess", guess).With("current number number", number).Info("got guess")

	if guessInt > number {
		statement, _ := h.ChatGptService.TextGeneration(ctx, higherThanprompt, h.Model)
		res = alexa.NewResponse("Random Number Game", statement, false)
	}
	if guessInt < number {
		statement, _ := h.ChatGptService.TextGeneration(ctx, lessThanprompt, h.Model)
		res = alexa.NewResponse("Random Number Game", statement, false)
	}
	if guessInt == number {
		winningStatement := fmt.Sprintf(winningprompt, h.RandomNumberSvc.Number)
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
