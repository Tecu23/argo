package killer

import "github.com/Tecu23/argov2/pkg/move"

// Table max sizes
const (
	MaxKillers = 2
	MaxPly     = 64
)

// Table represent the implementation of the killer table
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

	// Also update killer moves at sibling plies for better move ordering
	t.UpdateSiblingKillers(mv, ply)
}

// UpdateSiblingKillers adds killer move to sibling moves at similar plies
func (t *Table) UpdateSiblingKillers(mv move.Move, ply int) {
	// Update killer moves for nearby plies (siblings)
	if ply > 0 && ply < MaxPly-1 {
		// Add with lower priority to adjacent plies
		// First check if it's already a killer at those plies
		if !t.IsKiller(mv, ply-1) {
			// Shift existing killers, keeping the first one
			t.moves[ply-1][MaxKillers-1] = mv
		}

		if !t.IsKiller(mv, ply+1) {
			// Shift existing killers, keeping the first one
			t.moves[ply+1][MaxKillers-1] = mv
		}
	}
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

	// Check for killers in sibling plies with lower priority
	if ply > 0 {
		for i := 0; i < MaxKillers; i++ {
			if t.moves[ply-1][i] == mv {
				return 8000 - i*100 // Sibling killer gets lower score
			}
		}
	}

	if ply < MaxPly-1 {
		for i := 0; i < MaxKillers; i++ {
			if t.moves[ply+1][i] == mv {
				return 8000 - i*100 // Sibling killer gets lower score
			}
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
