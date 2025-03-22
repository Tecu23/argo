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

// Global arrays to store the network weights and biases.
// These are populated by the LoadWeights function and used during evaluation.

var (
	// InputWeights maps each input feature (based on piece, square, etc.) to the hidden layer neurons.
	InputWeights [InputSize][HiddenSize]int16
	// HiddenWeights connects the concatenated hidden layer output to the final evaluation output.
	HiddenWeights [OutputSize][HiddenDSize]int16
	// InputBias provides the baseline activation for the hidden layer neurons.
	InputBias [HiddenSize]int16
	// HiddenBias is added to the final weighted sum from the hidden layer.
	HiddenBias [OutputSize]int32
)

// Static piece value constants used during evaluation.
// They provide bonus values for pieces in the middlegame and endgame.
// The bonus for the king is a high constant to reflect its critical importance.
const (
	// for middlegame
	pawnBonusMG   = 124
	knightBonusMG = 781
	bishopBonusMG = 825
	rookBonusMG   = 1276
	queenBonusMG  = 2538

	// for endgame
	pawnBonusEG   = 206
	knightBonusEG = 854
	bishopBonusEG = 915
	rookBonusEG   = 1380
	queenBonusEG  = 2682
)

// FileIndex returns the file (column) index for a given square.
// It masks the square index with 7 (binary 111) to get a value in the range 0-7
func FileIndex(square int) int {
	return square & 7
}
