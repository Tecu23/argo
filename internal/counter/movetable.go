package counter

import (
	"github.com/Tecu23/argov2/pkg/move"
)

// MoveTable stores good response moves to specific opponent moves
type MoveTable struct {
	// [piece][square] -> move that's good against it
	moves [12][64]move.Move
}

// New creates a new counter move table
func New() *MoveTable {
	return &MoveTable{}
}

// Clear resets all counter moves
func (c *MoveTable) Clear() {
	c.moves = [12][64]move.Move{}
}

// Update adds a counter move for a previous opponent move
func (c *MoveTable) Update(prevMove, goodMove move.Move) {
	if prevMove == move.NoMove {
		return
	}

	piece := prevMove.GetPiece()
	to := prevMove.GetTarget()

	// Store this move as a good response to the previous move
	c.moves[piece][to] = goodMove
}

// Get returns the counter move score if the current move is a known counter to the previous move
func (c *MoveTable) Get(prevMove, currMove move.Move) int {
	if prevMove == move.NoMove {
		return 0
	}

	piece := prevMove.GetPiece()
	to := prevMove.GetTarget()

	if c.moves[piece][to] == currMove {
		return 500 // Counter move bonus
	}

	return 0
}

// GetMove returns the stored counter move for a specific previous move
func (c *MoveTable) GetMove(prevMove move.Move) move.Move {
	if prevMove == move.NoMove {
		return move.NoMove
	}

	piece := prevMove.GetPiece()
	to := prevMove.GetTarget()

	return c.moves[piece][to]
}
