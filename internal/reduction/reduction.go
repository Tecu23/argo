package reduction

import "math"

type Table struct {
	reductions [64][64]int // [depth][moveNumber]
}

// Base constants for reduction table
const (
	BaseDepthReduction   = 0.85
	BaseMoveNumReduction = 0.80
	BaseReductionDivisor = 2.2
)

// Limits
const (
	MaxReduction            = 4
	MinDepthForReduction    = 3
	MinMovesBeforeReduction = 3

	// PV related adjustments
	PVNodeReductionPenalty = 1 // Reduce less in PV nodes

	// History heuristic thresholds
	HistoryScoreThreshold = 8000 // Threashold for good history score (will reduce reduction)
)

func New() *Table {
	t := &Table{}
	t.initialize()
	return t
}

func (t *Table) initialize() {
	for depth := range 64 {
		for moveNumber := range 64 {
			t.reductions[depth][moveNumber] = t.calculateReduction(depth, moveNumber)
		}
	}
}

func (t *Table) calculateReduction(depth, moveNumber int) int {
	if depth < MinDepthForReduction || moveNumber < MinMovesBeforeReduction {
		return 0
	}

	reduction := float64(0)

	// Use logarithmic scale for both depth and move number
	// This provides a smoother curve and better balance
	depthComponent := BaseDepthReduction * math.Log(float64(depth))
	moveComponent := BaseMoveNumReduction * math.Log(float64(moveNumber))

	// Formula: R = (ln(depth) * ln(moveNumber)) / divider
	reduction = (depthComponent * moveComponent) / BaseReductionDivisor

	// Convert to integer with proper rounding
	r := int(math.Floor(reduction + 0.5)) // Round to nearest integer

	// Safety constraints
	r = max(r, 0)
	r = min(r, MaxReduction)

	// Ensure we don't reduce too much relative to the remaining depth
	maxAllowedReduction := depth / 2
	maxAllowedReduction = max(maxAllowedReduction, 1)

	if r >= maxAllowedReduction {
		r = maxAllowedReduction - 1
	}

	// Ensure we always have at least one ply left
	if depth-r <= 1 {
		r = depth - 2
		r = max(r, 0)
	}

	return r
}

func (t *Table) Get(depth, moveNumber int) int {
	if depth >= 64 || moveNumber >= 64 || depth < MinDepthForReduction ||
		moveNumber < MinMovesBeforeReduction {
		return 0
	}

	return t.reductions[depth][moveNumber]
}

// GetWithAdjustments provides a more dynamic reduction that takes
// into account factors like PV nodes and history scores
func (t *Table) GetWithAdjustments(depth, moveNumber int, isPV bool, historyScore int) int {
	// Start with base reduction
	r := t.Get(depth, moveNumber)

	// Adjust for PV nodes - reduce less for principal variation
	if isPV {
		r -= PVNodeReductionPenalty
	}

	// Adjust for history score - reduce less for moves with good history
	if historyScore > HistoryScoreThreshold {
		r--
	}

	// Apply bounds
	if r < 0 {
		r = 0
	}
	if r > MaxReduction {
		r = MaxReduction
	}

	return r
}

// Methods for working with different kinds of moves

// ShouldReduce determines if a move should be reduced based on its properties
func (t *Table) ShouldReduce(
	depth, moveNumber int,
	isCheck, isCapture, isPromotion, givesCheck bool,
) bool {
	return depth >= MinDepthForReduction &&
		moveNumber >= MinMovesBeforeReduction &&
		!isCheck &&
		!isCapture &&
		!isPromotion &&
		!givesCheck
}
