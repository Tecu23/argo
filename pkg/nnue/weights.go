// Package nnue keeps the NNUE (Efficiently Updated Neural Network) responsible for
// evaluation the current position
package nnue

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
