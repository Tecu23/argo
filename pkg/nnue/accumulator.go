package nnue

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// Accumulator represents the efficiently updatable first layer
type Accumulator struct {
	Summation [2][HiddenSize]int16 // [color][neuron]
}

// AccumulatorTableEntry caches a position's accumulator state
type AccumulatorTableEntry struct {
	PieceOcc    [2][6]bitboard.Bitboard // [color][piece type] bitboards
	Accumulator Accumulator
}

// AccumulatorTable caches accumulators for different king positions
type AccumulatorTable struct {
	Entries [2][32]AccumulatorTableEntry // [color][kingIndex]
}

// Reset initializes the accumulator table with bias values
func (a *AccumulatorTable) Reset() {
	for c := 0; c < 2; c++ {
		for s := 0; s < 32; s++ {
			// Copy input bias to initialize accumulator
			for i := 0; i < HiddenSize; i++ {
				a.Entries[c][s].Accumulator.Summation[c][i] = InputBias[i]
			}
		}
	}
}

// Use updates the accumulator based on current board position
func (a *AccumulatorTable) Use(view int, b *board.Board, evaluator *Evaluator) {
	var kingSq int
	if view == White {
		kingBB := b.Bitboards[WK]
		kingSq = kingBB.FirstOne()
	} else {
		kingBB := b.Bitboards[BK]
		kingSq = kingBB.FirstOne()

	}

	kingSq = ConvertSquare(kingSq)

	kingSide := FileIndex(kingSq) > 3
	ksIndex := KingSquareIndex(kingSq, view)

	// Determine entry index
	entryIdx := 0
	if kingSide {
		entryIdx = 16 + ksIndex
	} else {
		entryIdx = ksIndex
	}

	// Get the entry
	entry := &a.Entries[view][entryIdx]

	// Update the accumulator based on piece differences
	for c := 0; c < 2; c++ {
		for pt := 0; pt < 6; pt++ {
			boardBB := b.GetPieceBB(c, pt)
			entryBB := entry.PieceOcc[c][pt]

			// Squares where pieces need to be added
			toSet := boardBB & ^entryBB

			// Squares where pieces need to be removed
			toUnset := entryBB & ^boardBB

			// Add new pieces
			for toSet != 0 {
				sq := toSet.FirstOne()
				sq = ConvertSquare(sq)
				idx := Index(pt, c, sq, view, kingSq)
				AddWeightsToAccumulator(
					true,
					idx,
					entry.Accumulator.Summation[view][:],
					entry.Accumulator.Summation[view][:],
				)
			}
			// Remove old pieces
			for toUnset != 0 {
				sq := toUnset.FirstOne()
				sq = ConvertSquare(sq)
				idx := Index(pt, c, sq, view, kingSq)
				AddWeightsToAccumulator(
					false,
					idx,
					entry.Accumulator.Summation[view][:],
					entry.Accumulator.Summation[view][:],
				)
			}

			// Store the updated piece bitboard
			entry.PieceOcc[c][pt] = boardBB
		}
	}

	// Copy the updated accumulator to the evaluator
	copy(
		evaluator.History[evaluator.HistoryIndex].Summation[view][:],
		entry.Accumulator.Summation[view][:],
	)
}

// AddWeightsToAccumulator adds or subtracts weights to/from an accumulator
func AddWeightsToAccumulator(add bool, idx int, src, target []int16) {
	for i := 0; i < len(src); i++ {
		if add {
			target[i] = src[i] + InputWeights[idx][i]
		} else {
			target[i] = src[i] - InputWeights[idx][i]
		}
	}
}

// Update helpers for different move types
func SetUnsetPiece(input, output *Accumulator, side int, set, unset FeatureIndex) {
	idx1 := set.Get(side)
	idx2 := unset.Get(side)

	for i := 0; i < HiddenSize; i++ {
		output.Summation[side][i] = input.Summation[side][i] +
			InputWeights[idx1][i] -
			InputWeights[idx2][i]
	}
}

func SetUnsetPieceBothColors(input, output *Accumulator, set, unset FeatureIndex) {
	SetUnsetPiece(input, output, White, set, unset)
	SetUnsetPiece(input, output, Black, set, unset)
}

func SetUnsetUnsetPiece(input, output *Accumulator, side int, set, unset1, unset2 FeatureIndex) {
	idx1 := set.Get(side)
	idx2 := unset1.Get(side)
	idx3 := unset2.Get(side)

	for i := 0; i < HiddenSize; i++ {
		output.Summation[side][i] = input.Summation[side][i] +
			InputWeights[idx1][i] -
			InputWeights[idx2][i] -
			InputWeights[idx3][i]
	}
}

func SetUnsetUnsetPieceBothColors(input, output *Accumulator, set, unset1, unset2 FeatureIndex) {
	SetUnsetUnsetPiece(input, output, White, set, unset1, unset2)
	SetUnsetUnsetPiece(input, output, Black, set, unset1, unset2)
}

func SetSetUnsetUnsetPiece(
	input, output *Accumulator,
	side int,
	set1, set2, unset1, unset2 FeatureIndex,
) {
	idx1 := set1.Get(side)
	idx2 := set2.Get(side)
	idx3 := unset1.Get(side)
	idx4 := unset2.Get(side)

	for i := 0; i < HiddenSize; i++ {
		output.Summation[side][i] = input.Summation[side][i] +
			InputWeights[idx1][i] +
			InputWeights[idx2][i] -
			InputWeights[idx3][i] -
			InputWeights[idx4][i]
	}
}

func SetSetUnsetUnsetPieceBothColors(
	input, output *Accumulator,
	set1, set2, unset1, unset2 FeatureIndex,
) {
	SetSetUnsetUnsetPiece(input, output, White, set1, set2, unset1, unset2)
	SetSetUnsetUnsetPiece(input, output, Black, set1, set2, unset1, unset2)
}
