package reduction

import "math"

type Table struct {
	reductions [64][64]int // [depth][moveNumber]
}

const (
	BaseDepthReduction      = 0.75
	BaseMoveNumReduction    = 0.75
	MaxReduction            = 3
	MinDepthForReduction    = 3
	MinMovesBeforeReduction = 4
)

func New() *Table {
	t := &Table{}
	t.initialize()

	return t
}

func (t *Table) initialize() {
	for depth := 0; depth < 64; depth++ {
		for moveNumber := 0; moveNumber < 64; moveNumber++ {
			t.reductions[depth][moveNumber] = t.calculateReduction(depth, moveNumber)
		}
	}
}

func (t *Table) calculateReduction(depth, moveNumber int) int {
	if depth < MinDepthForReduction || moveNumber < MinMovesBeforeReduction {
		return 0
	}

	reduction := float64(0)

	depthComponent := BaseDepthReduction * math.Log(float64(depth))
	moveComponent := BaseMoveNumReduction * math.Log(float64(moveNumber))

	reduction = depthComponent * moveComponent

	r := int(math.Floor(reduction))

	if r < 0 {
		r = 0
	}
	if r > MaxReduction {
		r = MaxReduction
	}

	if r >= depth/2 {
		r = depth/2 - 1
	}

	return r
}

func (t *Table) Get(depth, moveNumber int) int {
	if depth >= 64 || moveNumber >= 64 {
		return 0
	}

	return t.reductions[depth][moveNumber]
}
