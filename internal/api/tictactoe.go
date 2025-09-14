package api

import (
	"errors"
	"fmt"
	"strings"
)

const (
	TicTacToeEmpty = iota
	TicTacToeX
	TicTacToeO
)

type TicTacToe struct {
	board  [][]int
	size   int
	player int // Current player (TicTacToeX or TicTacToeO)
	winner int // Winner of the game
	moves  int // Number of moves made
}

func NewTicTacToe() *TicTacToe {
	size := 3
	board := make([][]int, size)
	for i := range board {
		board[i] = make([]int, size)
	}
	return &TicTacToe{
		board:  board,
		size:   size,
		player: TicTacToeX, // X goes first (human player)
		winner: TicTacToeEmpty,
		moves:  0,
	}
}

func (t *TicTacToe) MakeMove(row, col int) error {
	if row < 0 || row >= t.size || col < 0 || col >= t.size {
		return errors.New("move out of bounds")
	}
	if t.board[row][col] != TicTacToeEmpty {
		return errors.New("position already occupied")
	}
	if t.winner != TicTacToeEmpty {
		return errors.New("game already over")
	}

	t.board[row][col] = t.player
	t.moves++

	// Check for winner
	if t.checkWinner() {
		t.winner = t.player
	}

	// Switch player
	if t.player == TicTacToeX {
		t.player = TicTacToeO
	} else {
		t.player = TicTacToeX
	}

	return nil
}

func (t *TicTacToe) checkWinner() bool {
	// Check rows
	for i := 0; i < t.size; i++ {
		if t.board[i][0] != TicTacToeEmpty &&
			t.board[i][0] == t.board[i][1] &&
			t.board[i][1] == t.board[i][2] {
			return true
		}
	}

	// Check columns
	for j := 0; j < t.size; j++ {
		if t.board[0][j] != TicTacToeEmpty &&
			t.board[0][j] == t.board[1][j] &&
			t.board[1][j] == t.board[2][j] {
			return true
		}
	}

	// Check diagonals
	if t.board[0][0] != TicTacToeEmpty &&
		t.board[0][0] == t.board[1][1] &&
		t.board[1][1] == t.board[2][2] {
		return true
	}

	if t.board[0][2] != TicTacToeEmpty &&
		t.board[0][2] == t.board[1][1] &&
		t.board[1][1] == t.board[2][0] {
		return true
	}

	return false
}

func (t *TicTacToe) IsDraw() bool {
	return t.moves == t.size*t.size && t.winner == TicTacToeEmpty
}

func (t *TicTacToe) IsGameOver() bool {
	return t.winner != TicTacToeEmpty || t.IsDraw()
}

func (t *TicTacToe) GetWinner() int {
	return t.winner
}

func (t *TicTacToe) GetCurrentPlayer() int {
	return t.player
}

// MakeAIMove makes the best move for the AI player (O)
func (t *TicTacToe) MakeAIMove() (int, int, error) {
	if t.IsGameOver() {
		return -1, -1, errors.New("game already over")
	}

	// Simple AI: try to win, block opponent, or take center/corners
	bestRow, bestCol := t.findBestMove()

	if bestRow == -1 || bestCol == -1 {
		// Find first empty spot (fallback)
		for i := 0; i < t.size; i++ {
			for j := 0; j < t.size; j++ {
				if t.board[i][j] == TicTacToeEmpty {
					bestRow, bestCol = i, j
					break
				}
			}
			if bestRow != -1 {
				break
			}
		}
	}

	if bestRow == -1 || bestCol == -1 {
		return -1, -1, errors.New("no valid moves available")
	}

	err := t.MakeMove(bestRow, bestCol)
	return bestRow, bestCol, err
}

func (t *TicTacToe) findBestMove() (int, int) {
	// Try to win
	for i := 0; i < t.size; i++ {
		for j := 0; j < t.size; j++ {
			if t.board[i][j] == TicTacToeEmpty {
				t.board[i][j] = TicTacToeO
				if t.checkWinner() {
					t.board[i][j] = TicTacToeEmpty
					return i, j
				}
				t.board[i][j] = TicTacToeEmpty
			}
		}
	}

	// Try to block opponent
	for i := 0; i < t.size; i++ {
		for j := 0; j < t.size; j++ {
			if t.board[i][j] == TicTacToeEmpty {
				t.board[i][j] = TicTacToeX
				if t.checkWinner() {
					t.board[i][j] = TicTacToeEmpty
					return i, j
				}
				t.board[i][j] = TicTacToeEmpty
			}
		}
	}

	// Take center if available
	if t.board[1][1] == TicTacToeEmpty {
		return 1, 1
	}

	// Take corners
	corners := [][]int{{0, 0}, {0, 2}, {2, 0}, {2, 2}}
	for _, corner := range corners {
		if t.board[corner[0]][corner[1]] == TicTacToeEmpty {
			return corner[0], corner[1]
		}
	}

	return -1, -1
}

// GetBoardDescription returns a text description for image generation
func (t *TicTacToe) GetBoardDescription() string {
	var sb strings.Builder
	sb.WriteString("Generate a tic-tac-toe game board image with a 3x3 grid. ")
	sb.WriteString("The board should look clean and modern with clear grid lines. ")
	sb.WriteString("Place the following marks on the board:\n")

	for i := 0; i < t.size; i++ {
		for j := 0; j < t.size; j++ {
			if t.board[i][j] != TicTacToeEmpty {
				position := fmt.Sprintf("Row %d, Column %d", i+1, j+1)
				mark := "X"
				if t.board[i][j] == TicTacToeO {
					mark = "O"
				}
				sb.WriteString(fmt.Sprintf("- %s: %s\n", position, mark))
			}
		}
	}

	if t.winner == TicTacToeX {
		sb.WriteString("\nHighlight the winning line for X (player wins!).")
	} else if t.winner == TicTacToeO {
		sb.WriteString("\nHighlight the winning line for O (AI wins!).")
	} else if t.IsDraw() {
		sb.WriteString("\nShow 'DRAW' text over the board.")
	}

	return sb.String()
}

// GetTextBoard returns a simple text representation of the board
func (t *TicTacToe) GetTextBoard() string {
	var sb strings.Builder
	sb.WriteString("```\n")
	sb.WriteString("  1 2 3\n")
	for i := 0; i < t.size; i++ {
		sb.WriteString(fmt.Sprintf("%d ", i+1))
		for j := 0; j < t.size; j++ {
			switch t.board[i][j] {
			case TicTacToeEmpty:
				sb.WriteString(". ")
			case TicTacToeX:
				sb.WriteString("X ")
			case TicTacToeO:
				sb.WriteString("O ")
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString("```")
	return sb.String()
}
