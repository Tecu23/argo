package killer

import "github.com/Tecu23/argov2/pkg/move"

const (
	MaxKillers = 2
	MaxPly     = 64
)

type Table struct {
	// [ply][slot]
	moves [MaxPly][MaxKillers]move.Move
}

// New creates a new killer move table
func New() *Table {
	return &Table{}
}

// Clear resets all killer moves
func (t *Table) Clear() {
	t.moves = [MaxPly][MaxKillers]move.Move{}
}

// Update adds a new killer move at the given ply
func (t *Table) Update(mv move.Move, ply int) {
	// Don't store captures as killer moves
	if mv.GetCapture() != 0 {
		return
	}

	// Don't store a move that's already a killer at this ply
	if t.IsKiller(mv, ply) {
		return
	}

	// Shift existing killers and insert new one at first position
	for i := MaxKillers - 1; i > 0; i-- {
		t.moves[ply][i] = t.moves[ply][i-1]
	}
	t.moves[ply][0] = mv
}

// IsKiller checks if a move is a killer move at the given ply
func (t *Table) IsKiller(mv move.Move, ply int) bool {
	for i := 0; i < MaxKillers; i++ {
		if t.moves[ply][i] == mv {
			return true
		}
	}
	return false
}

// GetScore returns the killer move score (0 if not a killer)
func (t *Table) GetScore(mv move.Move, ply int) int {
	for i := 0; i < MaxKillers; i++ {
		if t.moves[ply][i] == mv {
			return 9000 - i*100 // First killer gets higher score
		}
	}
	return 0
}

// Get returns all killer moves for a given ply
func (t *Table) Get(ply int) []move.Move {
	killers := make([]move.Move, 0, MaxKillers)
	for i := 0; i < MaxKillers; i++ {
		if t.moves[ply][i] != move.NoMove {
			killers = append(killers, t.moves[ply][i])
		}
	}
	return killers
}
