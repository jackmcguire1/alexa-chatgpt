package api

import (
	"testing"
)

func TestNewTicTacToe(t *testing.T) {
	game := NewTicTacToe()
	if game == nil {
		t.Fatal("Expected non-nil game")
	}
	if game.size != 3 {
		t.Errorf("Expected size 3, got %d", game.size)
	}
	if game.player != TicTacToeX {
		t.Errorf("Expected X to go first")
	}
	if game.winner != TicTacToeEmpty {
		t.Errorf("Expected no winner at start")
	}
}

func TestTicTacToeMakeMove(t *testing.T) {
	game := NewTicTacToe()

	// Test valid move
	err := game.MakeMove(0, 0)
	if err != nil {
		t.Errorf("Expected valid move, got error: %v", err)
	}
	if game.board[0][0] != TicTacToeX {
		t.Errorf("Expected X at (0,0)")
	}
	if game.player != TicTacToeO {
		t.Errorf("Expected player to switch to O")
	}

	// Test invalid move (same position)
	err = game.MakeMove(0, 0)
	if err == nil {
		t.Errorf("Expected error for occupied position")
	}

	// Test out of bounds
	err = game.MakeMove(3, 3)
	if err == nil {
		t.Errorf("Expected error for out of bounds move")
	}
}

func TestTicTacToeWinConditions(t *testing.T) {
	// Test horizontal win
	game := NewTicTacToe()
	game.MakeMove(0, 0) // X
	game.MakeMove(1, 0) // O
	game.MakeMove(0, 1) // X
	game.MakeMove(1, 1) // O
	game.MakeMove(0, 2) // X wins

	if game.winner != TicTacToeX {
		t.Errorf("Expected X to win horizontally")
	}

	// Test vertical win
	game = NewTicTacToe()
	game.MakeMove(0, 0) // X
	game.MakeMove(0, 1) // O
	game.MakeMove(1, 0) // X
	game.MakeMove(1, 1) // O
	game.MakeMove(2, 0) // X wins

	if game.winner != TicTacToeX {
		t.Errorf("Expected X to win vertically")
	}

	// Test diagonal win
	game = NewTicTacToe()
	game.MakeMove(0, 0) // X
	game.MakeMove(0, 1) // O
	game.MakeMove(1, 1) // X
	game.MakeMove(0, 2) // O
	game.MakeMove(2, 2) // X wins

	if game.winner != TicTacToeX {
		t.Errorf("Expected X to win diagonally")
	}
}

func TestTicTacToeDraw(t *testing.T) {
	game := NewTicTacToe()
	// Play a draw game
	game.MakeMove(0, 0) // X
	game.MakeMove(0, 1) // O
	game.MakeMove(0, 2) // X
	game.MakeMove(1, 1) // O
	game.MakeMove(1, 0) // X
	game.MakeMove(1, 2) // O
	game.MakeMove(2, 1) // X
	game.MakeMove(2, 0) // O
	game.MakeMove(2, 2) // X

	if !game.IsDraw() {
		t.Errorf("Expected game to be a draw")
	}
	if game.winner != TicTacToeEmpty {
		t.Errorf("Expected no winner in a draw")
	}
}

func TestTicTacToeAIMove(t *testing.T) {
	game := NewTicTacToe()

	// Make a move as X
	game.MakeMove(0, 0)

	// AI should make a valid move
	row, col, err := game.MakeAIMove()
	if err != nil {
		t.Errorf("AI move failed: %v", err)
	}
	if row < 0 || row >= 3 || col < 0 || col >= 3 {
		t.Errorf("AI made invalid move: (%d, %d)", row, col)
	}
	if game.board[row][col] != TicTacToeO {
		t.Errorf("AI didn't place O on the board")
	}
}

func TestTicTacToeAIBlocking(t *testing.T) {
	game := NewTicTacToe()

	// Set up a scenario where X is about to win
	game.MakeMove(0, 0) // X
	game.MakeMove(1, 1) // O
	game.MakeMove(0, 1) // X (X has two in a row)

	// AI should block at (0, 2)
	row, col, _ := game.MakeAIMove()
	if row != 0 || col != 2 {
		t.Errorf("AI didn't block winning move, played at (%d, %d) instead of (0, 2)", row, col)
	}
}

func TestTicTacToeGetBoardDescription(t *testing.T) {
	game := NewTicTacToe()
	game.MakeMove(0, 0) // X
	game.MakeMove(1, 1) // O

	desc := game.GetBoardDescription()
	if desc == "" {
		t.Errorf("Expected non-empty board description")
	}
	// Check that description mentions the marks
	if !contains(desc, "Row 1, Column 1: X") {
		t.Errorf("Expected description to mention X at Row 1, Column 1")
	}
	if !contains(desc, "Row 2, Column 2: O") {
		t.Errorf("Expected description to mention O at Row 2, Column 2")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}
