package api

import (
	"errors"
	"fmt"
)

const (
	Empty = iota
	IntShip
	Hit
	Sink
	Miss
	Invalid
	GameOver
)

type Ship struct {
	X, Y, Length int
	Horizontal   bool
	Killed       bool
}

type Battleships struct {
	board [][]int
	size  int
	ships map[int]*Ship
}

func NewBattleships(size int) *Battleships {
	board := make([][]int, size)
	for i := range board {
		board[i] = make([]int, size)
	}
	return &Battleships{board: board, size: size, ships: make(map[int]*Ship)}
}

func (b *Battleships) PlaceShip(x, y, length int, horizontal bool) error {
	if horizontal {
		if y+length > b.size {
			return errors.New("ship placement out of bounds")
		}
		for i := range length {
			if b.board[x][y+i] != Empty {
				return errors.New("ship placement overlaps another ship")
			}
		}
		for i := range length {
			b.board[x][y+i] = IntShip
		}
	} else {
		if x+length > b.size {
			return errors.New("ship placement out of bounds")
		}
		for i := range length {
			if b.board[x+i][y] != Empty {
				return errors.New("ship placement overlaps another ship")
			}
		}
		for i := range length {
			b.board[x+i][y] = IntShip
		}
	}
	b.ships[len(b.ships)] = &Ship{X: x, Y: y, Length: length, Horizontal: horizontal}
	return nil
}

func (b *Battleships) Attack(x, y int) int {
	if x >= b.size || y >= b.size {
		return Invalid
	}
	switch b.board[x][y] {
	case Empty:
		b.board[x][y] = Miss
		return Miss
	case IntShip:
		b.board[x][y] = Hit
		killed := b.checkShipKilled(x, y)
		if killed {
			goto checkShips
		}
		return Hit
	case Hit, Miss, Sink:
		return Invalid
	default:
		return Invalid
	}
checkShips:
	if b.AllShipsKilled() {
		return GameOver
	}

	return Sink
}

func (b *Battleships) AllShipsKilled() bool {
	for _, ship := range b.ships {
		if !ship.Killed {
			return false
		}
	}
	return true
}

func (b *Battleships) checkShipKilled(x, y int) bool {
	for _, ship := range b.ships {
		if ship.Horizontal {
			if ship.X == x && y >= ship.Y && y < ship.Y+ship.Length {
				killed := b.isShipKilled(ship)
				if killed {
					ship.Killed = true
					return true
				}
				return false
			}
		} else {
			if ship.Y == y && x >= ship.X && x < ship.X+ship.Length {
				killed := b.isShipKilled(ship)
				if killed {
					ship.Killed = true
					return true
				}
			}
		}
	}
	return false
}

func (b *Battleships) isShipKilled(ship *Ship) bool {
	if ship.Horizontal {
		for i := 0; i < ship.Length; i++ {
			if b.board[ship.X][ship.Y+i] != Hit {
				return false
			}
		}
	} else {
		for i := 0; i < ship.Length; i++ {
			if b.board[ship.X+i][ship.Y] != Hit {
				return false
			}
		}
	}
	return true
}

func (b *Battleships) ShipsTotals() (int, int) {
	var alive, killed = 0, 0
	for _, ship := range b.ships {
		if ship.Killed {
			killed++
		} else {
			alive++
		}
	}
	return alive, killed
}

func (b *Battleships) TotalHitsAndMisses() (int, int) {
	hits, misses := 0, 0
	for i := 0; i < len(b.board); i++ {
		for j := 0; j < len(b.board[i]); j++ {
			if b.board[i][j] == Miss {
				misses++
			}
			if b.board[i][j] == Hit {
				hits++
			}
		}
	}
	return hits, misses
}

func (b *Battleships) PrintBoard() {
	for _, row := range b.board {
		for _, cell := range row {
			switch cell {
			case Empty:
				fmt.Print(". ")
			case IntShip:
				fmt.Print("S ")
			case Hit:
				fmt.Print("X ")
			case Miss:
				fmt.Print("O ")
			}
		}
		fmt.Println()
	}
}

/*
BOARD SETUP
x 0 1 2 3 4 5 6 7 8 9
0 S S S S S . . . . .
1 . . . . . . . . . .
2 . . S . . . . . . .
3 . . S . . . . . . .
4 . . S . . . . . . .
5 . . S . . S S S . .
6 . . . . . . . . . .
7 . S . . . . . . . .
8 . S . . . . . . . .
9 . . . . . . . . . S
*/
func NewBattleShipSetup() *Battleships {
	board := NewBattleships(10)

	// Place ships at predefined positions
	_ = board.PlaceShip(0, 0, 5, true)  // Horizontal ship of length 5 at (0, 0)
	_ = board.PlaceShip(2, 2, 4, false) // Vertical ship of length 4 at (2, 2)
	_ = board.PlaceShip(5, 5, 3, true)  // Horizontal ship of length 3 at (5, 5)
	_ = board.PlaceShip(7, 1, 2, false) // Vertical ship of length 2 at (7, 1)
	_ = board.PlaceShip(9, 9, 1, true)  // Horizontal ship of length 1 at (9, 9)
	return board
}
