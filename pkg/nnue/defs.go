// Package nnue keeps the NNUE (Efficiently Updated Neural Network) responsible for
// evaluation the current position
package nnue

// Network dimensions for NNUE evaluation
const (
	// InputSize: total number of input features calculated as:
	// piece types (6) * squares (64) * colors (2) * king buckets (16)
	InputSize = 6 * 64 * 2 * 16

	// Hidden layer size (number of hidden neurons)
	HiddenSize = 512

	// Double hidden layer size for output connections (concatenation of two halfs)
	HiddenDSize = HiddenSize * 2

	// Output size: one evaluation score
	OutputSize = 1

	// Multipliers used to scale the network weights during evaluations
	InputWeightMultiplier  = 32
	HiddenWeightMultiplier = 128
)

// Color constants representing White and Black
const (
	White = 0
	Black = 1
)

// FileIndex returns the file (column) index for a given square.
// It masks the square index with 7 (binary 111) to get a value in the range 0-7
func FileIndex(square int) int {
	return square & 7
}
