package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnimalGameInitialization(t *testing.T) {
	game := NewAnimalGame()

	assert.NotNil(t, game)
	assert.Equal(t, MaxGuesses, game.guessesLeft)
	assert.Equal(t, MaxHints, game.hintsLeft)
	assert.True(t, game.gameActive)
	assert.NotEmpty(t, game.animal)
	assert.Empty(t, game.hintsGiven)
}

func TestNewAnimalGameSelectsRandomAnimal(t *testing.T) {
	// Create multiple games and verify different animals are selected
	animals := make(map[string]bool)

	for range 20 {
		game := NewAnimalGame()
		animals[game.animal] = true
	}

	// With 20 animals in the database and 20 iterations,
	// we should get at least a few different animals
	assert.Greater(t, len(animals), 1, "Should select different animals randomly")
}

func TestAnimalDatabaseIsNotEmpty(t *testing.T) {
	assert.Greater(t, len(animalDatabase), 0, "Animal database should not be empty")

	// Verify each animal has a name and sound
	for _, animal := range animalDatabase {
		assert.NotEmpty(t, animal.Name, "Animal name should not be empty")
		assert.NotEmpty(t, animal.Sound, "Animal sound should not be empty")
	}
}

func TestMakeGuessCorrect(t *testing.T) {
	game := NewAnimalGame()
	correctAnimal := game.animal

	result := game.MakeGuess(correctAnimal)

	assert.Equal(t, GameWon, result.Status)
	assert.Equal(t, correctAnimal, result.Animal)
	assert.NotEmpty(t, result.Sound)
	assert.False(t, game.gameActive)
}

func TestMakeGuessCaseInsensitive(t *testing.T) {
	game := NewAnimalGame()
	correctAnimal := game.animal

	// Test with uppercase
	result := game.MakeGuess(strings.ToUpper(correctAnimal))

	assert.Equal(t, GameWon, result.Status)
	assert.Equal(t, correctAnimal, result.Animal)
	assert.False(t, game.gameActive)
}

func TestMakeGuessWithWhitespace(t *testing.T) {
	game := NewAnimalGame()
	correctAnimal := game.animal

	// Test with leading/trailing whitespace
	result := game.MakeGuess("  " + correctAnimal + "  ")

	assert.Equal(t, GameWon, result.Status)
	assert.Equal(t, correctAnimal, result.Animal)
	assert.False(t, game.gameActive)
}

func TestMakeGuessIncorrect(t *testing.T) {
	game := NewAnimalGame()
	initialGuesses := game.guessesLeft

	result := game.MakeGuess("nonexistent_animal_xyz")

	assert.Equal(t, GuessIncorrect, result.Status)
	assert.Equal(t, initialGuesses-1, result.GuessesLeft)
	assert.Equal(t, initialGuesses-1, game.guessesLeft)
	assert.True(t, game.gameActive)
}

func TestMakeGuessDecrementsGuesses(t *testing.T) {
	game := NewAnimalGame()

	for i := MaxGuesses; i > 1; i-- {
		result := game.MakeGuess("wrong_animal")
		assert.Equal(t, GuessIncorrect, result.Status)
		assert.Equal(t, i-1, game.guessesLeft)
		assert.True(t, game.gameActive)
	}
}

func TestMakeGuessGameLost(t *testing.T) {
	game := NewAnimalGame()
	correctAnimal := game.animal

	// Use up all guesses with wrong answers
	for i := range MaxGuesses {
		result := game.MakeGuess("wrong_animal")

		if i < MaxGuesses-1 {
			assert.Equal(t, GuessIncorrect, result.Status)
			assert.True(t, game.gameActive)
		} else {
			assert.Equal(t, GameLost, result.Status)
			assert.Equal(t, correctAnimal, result.Animal)
			assert.False(t, game.gameActive)
			assert.Equal(t, 0, result.GuessesLeft)
		}
	}
}

func TestMakeGuessWhenGameInactive(t *testing.T) {
	game := NewAnimalGame()
	game.gameActive = false

	result := game.MakeGuess("any_animal")

	assert.Equal(t, GameInactive, result.Status)
}

func TestRequestHintSuccess(t *testing.T) {
	game := NewAnimalGame()
	initialHints := game.hintsLeft

	result := game.RequestHint()

	assert.Equal(t, HintAvailable, result.Status)
	assert.Equal(t, 1, result.HintNumber)
	assert.Equal(t, initialHints-1, result.HintsLeft)
	assert.Equal(t, game.animal, result.Animal)
	assert.Equal(t, MaxGuesses, result.GuessesLeft)
	assert.Contains(t, game.hintsGiven, 1)
}

func TestRequestHintDecrementsHints(t *testing.T) {
	game := NewAnimalGame()

	for i := 1; i <= MaxHints; i++ {
		result := game.RequestHint()

		assert.Equal(t, HintAvailable, result.Status)
		assert.Equal(t, i, result.HintNumber)
		assert.Equal(t, MaxHints-i, result.HintsLeft)
		assert.Equal(t, MaxHints-i, game.hintsLeft)
	}
}

func TestRequestHintWhenNoHintsLeft(t *testing.T) {
	game := NewAnimalGame()

	// Use up all hints
	for range MaxHints {
		game.RequestHint()
	}

	result := game.RequestHint()

	assert.Equal(t, NoHintsLeft, result.Status)
	assert.Equal(t, 0, result.HintsLeft)
	assert.Equal(t, 0, game.hintsLeft)
}

func TestRequestHintWhenGameInactive(t *testing.T) {
	game := NewAnimalGame()
	game.gameActive = false

	result := game.RequestHint()

	assert.Equal(t, GameInactive, result.Status)
}

func TestRequestHintTracksHintsGiven(t *testing.T) {
	game := NewAnimalGame()

	game.RequestHint()
	game.RequestHint()
	game.RequestHint()

	assert.Equal(t, 3, len(game.hintsGiven))
	assert.Contains(t, game.hintsGiven, 1)
	assert.Contains(t, game.hintsGiven, 2)
	assert.Contains(t, game.hintsGiven, 3)
}

func TestGetStatus(t *testing.T) {
	game := NewAnimalGame()

	status := game.GetStatus()

	assert.Equal(t, MaxGuesses, status.GuessesLeft)
	assert.Equal(t, MaxHints, status.HintsLeft)
	assert.True(t, status.GameActive)
}

func TestGetStatusAfterGuessesAndHints(t *testing.T) {
	game := NewAnimalGame()

	game.MakeGuess("wrong1")
	game.MakeGuess("wrong2")
	game.RequestHint()

	status := game.GetStatus()

	assert.Equal(t, MaxGuesses-2, status.GuessesLeft)
	assert.Equal(t, MaxHints-1, status.HintsLeft)
	assert.True(t, status.GameActive)
}

func TestResetGame(t *testing.T) {
	game := NewAnimalGame()
	originalAnimal := game.animal

	// Make some guesses and get hints
	game.MakeGuess("wrong1")
	game.MakeGuess("wrong2")
	game.RequestHint()
	game.RequestHint()

	game.ResetGame()

	assert.Equal(t, MaxGuesses, game.guessesLeft)
	assert.Equal(t, MaxHints, game.hintsLeft)
	assert.Empty(t, game.hintsGiven)
	assert.True(t, game.gameActive)
	assert.NotEmpty(t, game.animal)
	// New animal might be the same or different (random)
	_ = originalAnimal
}

func TestGetAnimalSound(t *testing.T) {
	game := NewAnimalGame()

	sound := game.getAnimalSound()

	assert.NotEmpty(t, sound)
	assert.NotEqual(t, "unknown animal sound", sound)
}

func TestAnimalSoundsAreCorrect(t *testing.T) {
	testCases := []struct {
		animal string
		sound  string
	}{
		{"dog", "woof woof"},
		{"cat", "meow meow"},
		{"cow", "moo moo"},
		{"lion", "roar roar"},
		{"elephant", "trumpet trumpet"},
	}

	for _, tc := range testCases {
		t.Run(tc.animal, func(t *testing.T) {
			game := NewAnimalGame()
			game.animal = tc.animal

			sound := game.getAnimalSound()
			assert.Equal(t, tc.sound, sound)
		})
	}
}

func TestGameFlowCompleteWin(t *testing.T) {
	game := NewAnimalGame()
	correctAnimal := game.animal

	// Request some hints
	hint1 := game.RequestHint()
	assert.Equal(t, HintAvailable, hint1.Status)

	hint2 := game.RequestHint()
	assert.Equal(t, HintAvailable, hint2.Status)

	// Make a few wrong guesses
	game.MakeGuess("wrong1")
	game.MakeGuess("wrong2")

	// Check status
	status := game.GetStatus()
	assert.Equal(t, MaxGuesses-2, status.GuessesLeft)
	assert.Equal(t, MaxHints-2, status.HintsLeft)

	// Make correct guess
	result := game.MakeGuess(correctAnimal)
	assert.Equal(t, GameWon, result.Status)
	assert.NotEmpty(t, result.Sound)
	assert.False(t, game.gameActive)
}

func TestGameFlowCompleteLoss(t *testing.T) {
	game := NewAnimalGame()
	correctAnimal := game.animal

	// Use all hints
	for range MaxHints {
		result := game.RequestHint()
		assert.Equal(t, HintAvailable, result.Status)
	}

	// Verify no hints left
	noHintResult := game.RequestHint()
	assert.Equal(t, NoHintsLeft, noHintResult.Status)

	// Use all guesses incorrectly
	for i := range MaxGuesses {
		result := game.MakeGuess("wrong_animal")
		if i < MaxGuesses-1 {
			assert.Equal(t, GuessIncorrect, result.Status)
		} else {
			assert.Equal(t, GameLost, result.Status)
			assert.Equal(t, correctAnimal, result.Animal)
		}
	}

	assert.False(t, game.gameActive)
}

func TestGameConstants(t *testing.T) {
	assert.Equal(t, 10, MaxGuesses)
	assert.Equal(t, 5, MaxHints)
}

func TestAnimalDatabaseCount(t *testing.T) {
	assert.GreaterOrEqual(t, len(animalDatabase), 15, "Should have at least 15 animals in database")
}
