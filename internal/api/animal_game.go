package api

import (
	"math/rand"
	"strings"
	"time"
)

// AnimalGame represents the state of the animal guessing game
type AnimalGame struct {
	animal        string
	guessesLeft   int
	hintsLeft     int
	hintsGiven    []int
	gameActive    bool
	maxGuesses    int
	maxHints      int
	animalList    []Animal
	animalSounds  map[string]string
}

// Animal represents an animal with its properties
type Animal struct {
	Name  string
	Sound string
}

// Game constants
const (
	MaxGuesses = 10
	MaxHints   = 5
)

// animalDatabase is the hardcoded list of animals
var animalDatabase = []Animal{
	// Farm Animals
	{Name: "dog", Sound: "woof woof"},
	{Name: "cat", Sound: "meow meow"},
	{Name: "cow", Sound: "moo moo"},
	{Name: "horse", Sound: "neigh neigh"},
	{Name: "pig", Sound: "oink oink"},
	{Name: "sheep", Sound: "baa baa"},
	{Name: "duck", Sound: "quack quack"},
	{Name: "chicken", Sound: "cluck cluck"},
	{Name: "rooster", Sound: "cock a doodle doo"},
	{Name: "goat", Sound: "maa maa"},
	{Name: "donkey", Sound: "hee haw"},
	{Name: "turkey", Sound: "gobble gobble"},

	// Wild Animals
	{Name: "lion", Sound: "roar roar"},
	{Name: "tiger", Sound: "roar roar"},
	{Name: "elephant", Sound: "trumpet trumpet"},
	{Name: "monkey", Sound: "ooh ooh ah ah"},
	{Name: "bear", Sound: "growl growl"},
	{Name: "wolf", Sound: "howl howl"},
	{Name: "fox", Sound: "yip yip"},
	{Name: "zebra", Sound: "whinny bark"},
	{Name: "giraffe", Sound: "hum hum"},
	{Name: "hippopotamus", Sound: "grunt grunt"},
	{Name: "rhinoceros", Sound: "snort snort"},
	{Name: "gorilla", Sound: "hoo hoo grunt"},
	{Name: "chimpanzee", Sound: "hoo hoo pant"},
	{Name: "kangaroo", Sound: "chuckle grunt"},
	{Name: "koala", Sound: "grunt snore"},
	{Name: "panda", Sound: "chirp bleat"},
	{Name: "sloth", Sound: "whistle scream"},
	{Name: "alligator", Sound: "hiss bellow"},
	{Name: "crocodile", Sound: "hiss growl"},

	// Birds
	{Name: "bird", Sound: "tweet tweet"},
	{Name: "owl", Sound: "hoot hoot"},
	{Name: "crow", Sound: "caw caw"},
	{Name: "parrot", Sound: "squawk squawk"},
	{Name: "eagle", Sound: "screech screech"},
	{Name: "penguin", Sound: "honk honk trumpet"},
	{Name: "flamingo", Sound: "honk grunt"},
	{Name: "peacock", Sound: "meow meow"},
	{Name: "woodpecker", Sound: "tap tap tap"},
	{Name: "seagull", Sound: "screech screech"},
	{Name: "pigeon", Sound: "coo coo"},

	// Small Animals & Insects
	{Name: "frog", Sound: "ribbit ribbit"},
	{Name: "snake", Sound: "hiss hiss"},
	{Name: "bee", Sound: "buzz buzz"},
	{Name: "cricket", Sound: "chirp chirp"},
	{Name: "mosquito", Sound: "buzz buzz"},
	{Name: "mouse", Sound: "squeak squeak"},
	{Name: "rat", Sound: "squeak squeak"},
	{Name: "hamster", Sound: "squeak squeak"},
	{Name: "rabbit", Sound: "thump thump"},
	{Name: "squirrel", Sound: "chatter chatter"},

	// Marine Animals
	{Name: "dolphin", Sound: "click click whistle"},
	{Name: "whale", Sound: "whoosh sing"},
	{Name: "seal", Sound: "arf arf"},
	{Name: "sea lion", Sound: "bark bark"},
	{Name: "walrus", Sound: "bellow roar"},
	{Name: "shark", Sound: "splash splash"},
	{Name: "octopus", Sound: "squirt squirt"},

	// Other Animals
	{Name: "hyena", Sound: "laugh cackle"},
	{Name: "cheetah", Sound: "chirp purr"},
	{Name: "leopard", Sound: "roar growl"},
	{Name: "raccoon", Sound: "chitter chitter"},
	{Name: "skunk", Sound: "hiss stomp"},
	{Name: "bat", Sound: "screech chirp"},
	{Name: "deer", Sound: "bleat grunt"},
	{Name: "moose", Sound: "bellow grunt"},
	{Name: "camel", Sound: "grunt groan"},
	{Name: "llama", Sound: "hum hum"},
	{Name: "armadillo", Sound: "grunt squeak"},
	{Name: "hedgehog", Sound: "snuffle grunt"},
	{Name: "porcupine", Sound: "chatter whine"},
}

// NewAnimalGame creates a new animal guessing game
func NewAnimalGame() *AnimalGame {
	game := &AnimalGame{
		guessesLeft: MaxGuesses,
		hintsLeft:   MaxHints,
		maxGuesses:  MaxGuesses,
		maxHints:    MaxHints,
		hintsGiven:  []int{},
		gameActive:  true,
		animalList:  animalDatabase,
	}
	game.selectRandomAnimal()
	return game
}

// selectRandomAnimal picks a random animal from the database
func (g *AnimalGame) selectRandomAnimal() {
	rand.Seed(time.Now().UnixNano())
	selectedAnimal := g.animalList[rand.Intn(len(g.animalList))]
	g.animal = selectedAnimal.Name
}

// MakeGuess processes a player's guess and returns the result
func (g *AnimalGame) MakeGuess(guess string) GuessResult {
	if !g.gameActive {
		return GuessResult{Status: GameInactive}
	}

	if g.guessesLeft <= 0 {
		g.gameActive = false
		return GuessResult{
			Status:      GameLost,
			Animal:      g.animal,
			GuessesLeft: 0,
		}
	}

	g.guessesLeft--

	// Case-insensitive comparison
	if strings.EqualFold(strings.TrimSpace(guess), g.animal) {
		g.gameActive = false
		animalSound := g.getAnimalSound()
		return GuessResult{
			Status:      GameWon,
			Animal:      g.animal,
			GuessesLeft: g.guessesLeft,
			Sound:       animalSound,
		}
	}

	if g.guessesLeft == 0 {
		g.gameActive = false
		return GuessResult{
			Status:      GameLost,
			Animal:      g.animal,
			GuessesLeft: 0,
		}
	}

	return GuessResult{
		Status:      GuessIncorrect,
		GuessesLeft: g.guessesLeft,
	}
}

// RequestHint returns a hint number if available
func (g *AnimalGame) RequestHint() HintResult {
	if !g.gameActive {
		return HintResult{Status: GameInactive}
	}

	if g.hintsLeft <= 0 {
		return HintResult{
			Status:    NoHintsLeft,
			HintsLeft: 0,
		}
	}

	g.hintsLeft--
	hintNumber := g.maxHints - g.hintsLeft
	g.hintsGiven = append(g.hintsGiven, hintNumber)

	return HintResult{
		Status:      HintAvailable,
		HintNumber:  hintNumber,
		HintsLeft:   g.hintsLeft,
		Animal:      g.animal,
		GuessesLeft: g.guessesLeft,
	}
}

// GetStatus returns the current game status
func (g *AnimalGame) GetStatus() GameStatus {
	return GameStatus{
		GuessesLeft: g.guessesLeft,
		HintsLeft:   g.hintsLeft,
		GameActive:  g.gameActive,
	}
}

// getAnimalSound returns the sound for the current animal
func (g *AnimalGame) getAnimalSound() string {
	for _, animal := range g.animalList {
		if animal.Name == g.animal {
			return animal.Sound
		}
	}
	return "unknown animal sound"
}

// ResetGame resets the game with a new animal
func (g *AnimalGame) ResetGame() {
	g.guessesLeft = MaxGuesses
	g.hintsLeft = MaxHints
	g.hintsGiven = []int{}
	g.gameActive = true
	g.selectRandomAnimal()
}

// Result status constants
const (
	GuessCorrect = iota
	GuessIncorrect
	GameWon
	GameLost
	GameInactive
	HintAvailable
	NoHintsLeft
)

// GuessResult represents the result of a guess
type GuessResult struct {
	Status      int
	Animal      string
	GuessesLeft int
	Sound       string
}

// HintResult represents the result of a hint request
type HintResult struct {
	Status      int
	HintNumber  int
	HintsLeft   int
	Animal      string
	GuessesLeft int
}

// GameStatus represents the current state of the game
type GameStatus struct {
	GuessesLeft int
	HintsLeft   int
	GameActive  bool
}
