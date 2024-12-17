package history

import "github.com/Tecu23/argov2/pkg/color"

const (
	historyMax = 10000
)

type HistoryTable struct {
	// [color][from][to]
	scores [2][64][64]int
}

func New() *HistoryTable {
	return &HistoryTable{}
}

func (h *HistoryTable) Clear() {
	h.scores = [2][64][64]int{}
}

func (h *HistoryTable) Update(color color.Color, from, to int, depth int) {
	score := h.scores[color][from][to] + depth*depth

	if score > historyMax {
		for f := 0; f < 64; f++ {
			for t := 0; t < 64; t++ {
				h.scores[color][f][t] /= 2
			}
		}
		score /= 2
	}

	h.scores[color][from][to] = score
}

func (h *HistoryTable) Get(color color.Color, from, to int) int {
	return h.scores[color][from][to]
}

func (h *HistoryTable) GetButterfly(from, to int) int {
	return h.scores[0][from][to] + h.scores[1][from][to]
}
