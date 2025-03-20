package nnue

const (
	// Network dimensions
	InputSize   = 6 * 64 * 2 * 16 // Piece types * squares * colors * king buckets
	HiddenSize  = 512
	HiddenDSize = HiddenSize * 2
	OutputSize  = 1

	// Weight multipliers
	InputWeightMultiplier  = 32
	HiddenWeightMultiplier = 128
)

// Color constants
const (
	White = 0
	Black = 1
)

func FileIndex(square int) int {
	return square & 7
}
