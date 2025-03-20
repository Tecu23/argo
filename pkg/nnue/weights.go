package nnue

// Network weight arrays
var (
	InputWeights  [InputSize][HiddenSize]int16
	HiddenWeights [OutputSize][HiddenDSize]int16
	InputBias     [HiddenSize]int16
	HiddenBias    [OutputSize]int32
)

// Piece value for evaluation
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
