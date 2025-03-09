package history

import "github.com/Tecu23/argov2/pkg/color"

const (
	historyMax = 10000
)

// Table stores success statistics for quiet moves
type Table struct {
	// [color][from][to]
	scores [2][64][64]int
}

// New creates a new history table
func New() *Table {
	return &Table{}
}

// New creates a new history table
func (h *Table) Clear() {
	h.scores = [2][64][64]int{}
}

func (h *Table) Update(color color.Color, from, to int, bonus int) {
	// Apply bonus
	h.scores[color][from][to] += bonus

	// If this score exceeds maximum, scale only this entry
	if h.scores[color][from][to] > historyMax {
		h.scores[color][from][to] = historyMax / 2
	}
}

func (h *Table) Get(color color.Color, from, to int) int {
	return h.scores[color][from][to]
}

func (h *Table) GetButterfly(from, to int) int {
	return h.scores[0][from][to] + h.scores[1][from][to]
}
