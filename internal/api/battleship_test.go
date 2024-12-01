package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBattleshipsCreatesEmptyBoard(t *testing.T) {
	size := 5
	game := NewBattleships(size)

	assert.Equal(t, size, game.size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			assert.Equal(t, Empty, game.board[i][j])
		}
	}
}

func TestPlaceShipHorizontalWithinBounds(t *testing.T) {
	game := NewBattleships(5)
	err := game.PlaceShip(0, 0, 3, true)

	assert.NoError(t, err)
	assert.Equal(t, IntShip, game.board[0][0])
	assert.Equal(t, IntShip, game.board[0][1])
	assert.Equal(t, IntShip, game.board[0][2])
}

func TestPlaceShipVerticalWithinBounds(t *testing.T) {
	game := NewBattleships(5)
	err := game.PlaceShip(0, 0, 3, false)

	assert.NoError(t, err)
	assert.Equal(t, IntShip, game.board[0][0])
	assert.Equal(t, IntShip, game.board[1][0])
	assert.Equal(t, IntShip, game.board[2][0])
}

func TestPlaceShipOutOfBounds(t *testing.T) {
	game := NewBattleships(5)
	err := game.PlaceShip(0, 3, 3, true)

	assert.Error(t, err)
}

func TestPlaceShipOverlapping(t *testing.T) {
	game := NewBattleships(5)
	_ = game.PlaceShip(0, 0, 3, true)
	err := game.PlaceShip(0, 1, 3, true)

	assert.Error(t, err)
}

func TestAttackMiss(t *testing.T) {
	game := NewBattleships(5)
	result := game.Attack(0, 0)

	assert.Equal(t, Miss, result)
	assert.Equal(t, Miss, game.board[0][0])
}

func TestAttackHit(t *testing.T) {
	game := NewBattleships(5)
	_ = game.PlaceShip(0, 0, 5, true)
	result := game.Attack(0, 0)

	assert.Equal(t, Hit, result)
	assert.Equal(t, Hit, game.board[0][0])
}

func TestAttackSinksShip(t *testing.T) {
	game := NewBattleships(5)
	_ = game.PlaceShip(0, 0, 1, true)
	_ = game.PlaceShip(4, 0, 2, true)
	result := game.Attack(0, 0)

	assert.Equal(t, Sink, result)
	assert.Equal(t, Hit, game.board[0][0])
}

func TestAttackOutOfBounds(t *testing.T) {
	game := NewBattleships(5)
	result := game.Attack(5, 5)

	assert.Equal(t, Invalid, result)
}

func TestAttackAlreadyAttacked(t *testing.T) {
	game := NewBattleships(5)
	_ = game.Attack(0, 0)
	result := game.Attack(0, 0)

	assert.Equal(t, Invalid, result)
}

func TestAttackEndsGame(t *testing.T) {
	game := NewBattleships(5)
	_ = game.PlaceShip(0, 0, 1, true)
	result := game.Attack(0, 0)

	assert.Equal(t, GameOver, result)
}
