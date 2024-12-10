// Package move contains the move and move list representation and all move helper functions.
// Move is represented as 64 bit unsigned integer where some bits represent some part of the move
// The first 6 bits keep the source square, the next 6 bits keep the target square and so on...
package move

import "fmt"

// Movelist is a dynamic list of moves. It can be appended to as moves are generated.
type Movelist []Move

// AddMove adds a move to the movelist.
func (m *Movelist) AddMove(move Move) {
	*m = append(*m, move)
}

// PrintMovelist prints all moves in the movelist with details for debugging.
func (m Movelist) PrintMovelist() {
	fmt.Printf("Move   Piece   Capture   Double   Enpassant   Castling\n")
	for _, move := range m {
		move.PrintMove()
	}
	fmt.Printf("\n\n Move Count: %d \n\n", len(m))
}
