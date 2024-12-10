package move

import "fmt"

// Movelist will keep track of all the moves being played throught the game
type Movelist []Move

// AddMove should add a move to the movelist
func (m *Movelist) AddMove(move Move) {
	*m = append(*m, move)
}

// PrintMovelist should print all moves that happen during the game
func (m Movelist) PrintMovelist() {
	fmt.Printf("Move   Piece   Capture   Double   Enpassant   Castling\n")
	for _, move := range m {
		move.PrintMove()
	}
	fmt.Printf("\n\n Move Count: %d \n\n", len(m))
}
