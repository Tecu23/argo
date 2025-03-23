// Copyright (C) 2025 Tecu23
// Port of Koivisto evaluation, licensed under GNU GPL v3

// File: accumulator.go

// Package nnue keeps the NNUE (Efficiently Updated Neural Network) responsible for
// evaluation the current position
package nnue

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// Accumulator represents the first layer of the neural network that is efficiently updatable.
// It stores a summation vector per color (0: White, 1: Black) with a fixed hidden size.
type Accumulator struct {
	Summation [2][HiddenSize]int16 // [color][neuron] stores the sum for each hidden neuron
}

// AccumulatorTableEntry caches the accumulator state for a specific board positon,
// along with the bitboards for piece occupancy per color and piece type.
type AccumulatorTableEntry struct {
	PieceOcc    [2][6]bitboard.Bitboard // [color][piece type] bitboards for caching updates
	Accumulator Accumulator             // The associated accumulator for this entry
}

// AccumulatorTable caches accumulators for different king positions.
// It is indexed by color and a precomputed king bucket index (0-31).
type AccumulatorTable struct {
	Entries [2][32]AccumulatorTableEntry // [color][kingIndex] mapping of cached accumulator entries
}

// Reset initializes the accumulator table with the network's input bias values.
// This is called to start evaluation with a baseline accumulator state.
func (a *AccumulatorTable) Reset() {
	for c := range 2 {
		for s := range 32 {
			for i := range HiddenSize {
				a.Entries[c][s].Accumulator.Summation[c][i] = InputBias[i]
			}
		}
	}
}

// Use updates the accumulator for a given perspective (view) based on the current board position.
// It calculates which king bucket (entry) to use based on the king's square and then adjustes the accumulator.
func (a *AccumulatorTable) Use(view int, b *board.Board, evaluator *Evaluator) {
	var kingSq int
	// Determine the king's square based on the perspective (White or Black)
	if view == White {
		kingBB := b.Bitboards[WK]
		kingSq = kingBB.FirstOne()
	} else {
		kingBB := b.Bitboards[BK]
		kingSq = kingBB.FirstOne()
	}

	// Convert from engine square representation to NNUE expected format
	kingSq = ConvertSquare(kingSq)

	// Determine if the king is on the king-side (file > 3) to choose the correct bucket
	kingSide := FileIndex(kingSq) > 3
	ksIndex := KingSquareIndex(kingSq, view)

	// Determine entry index based on king side
	entryIdx := 0
	if kingSide {
		entryIdx = 16 + ksIndex
	} else {
		entryIdx = ksIndex
	}

	// Get a reference to the cached accumulator entry for this king position and color
	entry := &a.Entries[view][entryIdx]

	// Loop over both colors and each piece type to update the accumulator based on changes in piece occupancy
	for c := range 2 {
		for pt := range 6 {
			boardBB := b.GetPieceBB(c, pt)   // Current board bitboard for this color and piece type
			entryBB := entry.PieceOcc[c][pt] // Cached bitboard from the previous state

			// Identify squares where pieces have been added (present on board but not in cache)
			toSet := boardBB & ^entryBB

			// Identify squares where pieces have been removed (present in cache but not on board)
			toUnset := entryBB & ^boardBB

			// Process pieces that need to be added
			for toSet != 0 {
				sq := toSet.FirstOne()
				sq = ConvertSquare(sq) // Convert Square to input into the NNUE
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
			// Update the cached piece occupancy to match the current board
			entry.PieceOcc[c][pt] = boardBB
		}
	}

	// Copy the updated accumulator to the evaluator's current history state
	copy(
		evaluator.History[evaluator.HistoryIndex].Summation[view][:],
		entry.Accumulator.Summation[view][:],
	)
}

// AddWeightsToAccumulator adds (or subtracts) network input weights to/from the accumulator.
// The 'add' flag determines if weights are added (true) or substracted (false)
func AddWeightsToAccumulator(add bool, idx int, src, target []int16) {
	addWeightsToAccumulatorASM(add, src, target, InputWeights[idx][:])
}

func addWeightsToAccumulatorASM(add bool, src, target, weights []int16)

// SetUnsetPieceBothColors applies the piece move update for both White and Black perspective
func SetUnsetPieceBothColors(input, output *Accumulator, set, unset FeatureIndex) {
	SetUnsetPiece(input, output, White, set, unset)
	SetUnsetPiece(input, output, Black, set, unset)
}

// SetUnsetPiece updates the accumulator when a piece moves from one square to another.
// It substracts the weights from the source square and adds the weights for the target square.
func SetUnsetPiece(input, output *Accumulator, side int, set, unset FeatureIndex) {
	idx1 := set.Get(side)
	idx2 := unset.Get(side)

	setUnsetPieceASM(
		input.Summation[side][:],
		output.Summation[side][:],
		InputWeights[idx1][:],
		InputWeights[idx2][:],
	)
}

func setUnsetPieceASM(input, output []int16, weightsSet, weightsUnset []int16)

// SetUnsetUnsetPiece updates the accumulator for moves involving a piece move with an additional removal.
// For example, when capturing, it adds the moving piece's weight, substracts the weight from the origin,
// and substracts the captured piece's weights
func SetUnsetUnsetPiece(input, output *Accumulator, side int, set, unset1, unset2 FeatureIndex) {
	idx1 := set.Get(side)
	idx2 := unset1.Get(side)
	idx3 := unset2.Get(side)

	setUnsetUnsetPieceASM(
		input.Summation[side][:],
		output.Summation[side][:],
		InputWeights[idx1][:],
		InputWeights[idx2][:],
		InputWeights[idx3][:],
	)
}

func setUnsetUnsetPieceASM(input, output []int16, set, unset1, unset2 []int16)

// SetUnsetUnsetPieceBothColors applies the above update for both colors.
func SetUnsetUnsetPieceBothColors(input, output *Accumulator, set, unset1, unset2 FeatureIndex) {
	SetUnsetUnsetPiece(input, output, White, set, unset1, unset2)
	SetUnsetUnsetPiece(input, output, Black, set, unset1, unset2)
}

// SetSetUnsetUnsetPiece handles cases where two pieces are set and two are removed in a single update.
// It is used for moves like castling where multiple pieces change positions.
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

// SetSetUnsetUnsetPieceBothColors applies the above castling or double-update for both colors
func SetSetUnsetUnsetPieceBothColors(
	input, output *Accumulator,
	set1, set2, unset1, unset2 FeatureIndex,
) {
	SetSetUnsetUnsetPiece(input, output, White, set1, set2, unset1, unset2)
	SetSetUnsetUnsetPiece(input, output, Black, set1, set2, unset1, unset2)
}
