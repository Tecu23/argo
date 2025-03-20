// Package nnue keeps the NNUE (Efficiently Updated Neural Network) responsible for
// evaluation the current position
package nnue

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
)

// Evaluator uses a neural network (NNUE) to evaluate chess positions.
// It maintains a history of accumulator states to allow incremental updates.
type Evaluator struct {
	History                  []Accumulator     // History stack of accumulators states for undo/redo moves
	HistoryIndex             int               // Current index in the history stack
	AccumulatorTable         *AccumulatorTable // Cached accumulators based on king positions
	AccumulatorIsInitialized [2]bool           // Flags to track whether accumulators have been initialized for each color
}

// NewEvaluator creates and initializes a new NNUE evaluator instance.
func NewEvaluator() *Evaluator {
	evaluator := &Evaluator{
		History:          []Accumulator{{}}, // Start with an initial accumulator state
		HistoryIndex:     0,
		AccumulatorTable: &AccumulatorTable{}, // Create a new table for caching accumulators
	}

	evaluator.AccumulatorTable.Reset()
	return evaluator
}

// Reset reinitializes the evaluator for a new board position.
// It resets the accumulator history and reinitializes accumulators for both colors.
func (e *Evaluator) Reset(b *board.Board) {
	e.History = []Accumulator{{}} // Clear history to initial state
	e.HistoryIndex = 0
	e.ResetAccumulator(b, White)
	e.ResetAccumulator(b, Black)
}

// ResetAccumulator reinitializes the accumulator for a specific perspective using the current board
func (e *Evaluator) ResetAccumulator(b *board.Board, color int) {
	e.AccumulatorTable.Use(color, b, e)
	e.AccumulatorIsInitialized[color] = true
}

var phaseValues = [5]float64{
	0.552938, 1.55294, 1.50862, 2.64379, 4.0,
}

// Evaluate computes a positional evaluation score for the current board.
// It scales between middlegame and endgame scores based on the phase of the game.
func (e *Evaluator) Evaluate(b *board.Board) int {
	const (
		evaluationMgScalar = 1.5     // Middlegame scaling factor
		evaluationEgScalar = 1.15    // Endgame scaling factor
		phaseSum           = 39.6684 // Total phase value sum for normalization
	)

	// Start with full phase and substract phase values based on the pieces remaining
	phase := phaseSum

	phase -= float64((b.Bitboards[WP] | b.Bitboards[BP]).Count()) * phaseValues[Pawn]
	phase -= float64((b.Bitboards[WN] | b.Bitboards[BN]).Count()) * phaseValues[Knight]
	phase -= float64((b.Bitboards[WB] | b.Bitboards[BB]).Count()) * phaseValues[Bishop]
	phase -= float64((b.Bitboards[WR] | b.Bitboards[BR]).Count()) * phaseValues[Rook]
	phase -= float64((b.Bitboards[WQ] | b.Bitboards[BQ]).Count()) * phaseValues[Queen]

	phase /= phaseSum // Normalize phase to a value between 0 and 1

	return int(
		(evaluationMgScalar - phase*(evaluationMgScalar-evaluationEgScalar)) * float64(
			e.eval(int(b.Side)),
		),
	)
}

// eval computes the raw neural network evaluation score using the current  accumulator state.
func (e *Evaluator) eval(activePlayer int) int {
	// Get accumulator values for the active and inactive sides
	accActive := e.History[e.HistoryIndex].Summation[activePlayer][:]
	accInactive := e.History[e.HistoryIndex].Summation[1-activePlayer][:]

	// Compute the output score
	var sum int32

	// Apply ReLU (max(0, x)) and compute the dot product for the active side
	for i := 0; i < HiddenSize; i++ {
		// ReLU: max(0, x)
		if accActive[i] > 0 {
			sum += int32(accActive[i]) * int32(HiddenWeights[0][i])
		}
	}

	// Similarly, for the inactive side using the corresponding half of hidden weights
	for i := 0; i < HiddenSize; i++ {
		if accInactive[i] > 0 {
			sum += int32(accInactive[i]) * int32(HiddenWeights[0][i+HiddenSize])
		}
	}

	sum += HiddenBias[0]

	// Scale the sum based on the weight multipliers to obtain the final evaluation score
	result := int(
		float64(sum) / float64(InputWeightMultiplier) / float64(HiddenWeightMultiplier),
	)

	return result
}

// GetPieceValue returns a static bonus value for a given piece in the middlegame.
// (The commented code hints at an endgame evaluation variant.)
func GetPieceValue(piece int) int {
	// if isEG {
	// 	switch piece {
	// 	case WP, BP:
	// 		return pawnBonusEG
	// 	case WN, BN:
	// 		return knightBonusEG
	// 	case WB, BB:
	// 		return bishopBonusEG
	// 	case WR, BR:
	// 		return rookBonusEG
	// 	case WQ, BQ:
	// 		return queenBonusEG
	// 	case WK, BK:
	// 		return 200_000
	// 	default:
	// 		return 0
	// 	}
	// }
	switch piece {
	case WP, BP:
		return pawnBonusMG
	case WN, BN:
		return knightBonusMG
	case WB, BB:
		return bishopBonusMG
	case WR, BR:
		return rookBonusMG
	case WQ, BQ:
		return queenBonusMG
	case WK, BK:
		return 200_000
	default:
		return 0
	}
}

// AddNewAccumulation adds a new accumulator state to the history stack,
// so that subsequent move updates are applied on a new state.
func (e *Evaluator) AddNewAccumulation() {
	e.HistoryIndex++

	// If the history slice is not long enough, expand it
	if e.HistoryIndex >= len(e.History) {
		e.History = append(e.History, Accumulator{})
	}

	// Mark both accumulators as not yet initialized for the new state
	e.AccumulatorIsInitialized[White] = false
	e.AccumulatorIsInitialized[Black] = false
}

// PopAccumulation undoes the last move by moving back in the history stack.
func (e *Evaluator) PopAccumulation() {
	e.HistoryIndex--
	e.AccumulatorIsInitialized[White] = true
	e.AccumulatorIsInitialized[Black] = true
}

// ClearHistory resets the accumulator history completely.
func (e *Evaluator) ClearHistory() {
	e.History = []Accumulator{{}}
	e.HistoryIndex = 0
}

// SetPieceOnSquare updates the accumulator for both perspectives when a piece is added or removed.
// It calls the per-color update function using both king positions.
func (e *Evaluator) SetPieceOnSquare(
	add bool,
	pieceType, pieceColor, square, wKingSq, bKingSq int,
) {
	square = square
	wKingSq = wKingSq
	bKingSq = bKingSq

	e.SetPieceOnSquareAccumulator(add, White, pieceType, pieceColor, square, wKingSq)
	e.SetPieceOnSquareAccumulator(add, Black, pieceType, pieceColor, square, bKingSq)
}

// SetPieceOnSquareAccumulator updates a single accumulator for a piece addition or removal.
// If the accumulator was not initialized, it uses the previous state's summation as a baseline.
func (e *Evaluator) SetPieceOnSquareAccumulator(
	add bool,
	side, pieceType, pieceColor, square, kingSq int,
) {
	idx := Index(pieceType, pieceColor, square, side, kingSq)

	if !e.AccumulatorIsInitialized[side] {
		// Use previous state as a baseline for the update
		AddWeightsToAccumulator(
			add,
			idx,
			e.History[e.HistoryIndex-1].Summation[side][:],
			e.History[e.HistoryIndex].Summation[side][:],
		)
		e.AccumulatorIsInitialized[side] = true
	} else {
		// Directly update the current accumulator
		AddWeightsToAccumulator(add, idx, e.History[e.HistoryIndex].Summation[side][:], e.History[e.HistoryIndex].Summation[side][:])
	}
}

// ProcessMove updates the accumulator based on a move on the board.
// It determines the move type (normal, capture, castling, en passant, promotion)
// and applies the appropriate update to the accumulator history.
func (e *Evaluator) ProcessMove(b *board.Board, m move.Move) {
	from := m.GetSource()
	to := m.GetTarget()
	piece := m.GetPiece()

	from = ConvertSquare(from)
	to = ConvertSquare(to)

	// moveType := m.MoveType
	isCapture := m.GetCapture() != 0
	isCastling := m.GetCastling() != 0
	isQueenCast := m.IsQueenCastle()
	enPass := m.GetEnpassant()

	capturedPiece := -1
	if isCapture {
		capturedPiece = b.GetPieceAt(to)
	}

	promPc := m.GetPromoted()
	c := White
	if b.Side == color.BLACK {
		c = Black
	}

	wKingBB := b.Bitboards[WK]
	wKingSq := wKingBB.FirstOne()
	bKingBB := b.Bitboards[BK]
	bKingSq := bKingBB.FirstOne()

	// Initialize for the new move
	e.AddNewAccumulation()

	if piece == King {
		// Handle king moves separately - may require full reset
		requiresReset := KingSquareIndex(to, c) != KingSquareIndex(from, c) ||
			FileIndex(from)+FileIndex(to) == 7

		if !requiresReset {
			if isCapture {
				SetUnsetUnsetPieceBothColors(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					FeatureIndex{piece, c, to, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
					FeatureIndex{capturedPiece, 1 - c, to, wKingSq, bKingSq},
				)
			} else if isCastling {
				val := 3
				if isQueenCast {
					val = -4
				}
				rookFrom := from + val
				val = -1
				if isQueenCast {
					val = 1
				}
				rookTo := to + val

				SetSetUnsetUnsetPieceBothColors(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					FeatureIndex{piece, c, to, wKingSq, bKingSq},
					FeatureIndex{Rook, c, rookTo, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
					FeatureIndex{Rook, c, rookFrom, wKingSq, bKingSq},
				)
			} else {
				SetUnsetPieceBothColors(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					FeatureIndex{piece, c, to, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
				)
			}
		} else {
			// Handle the opponent's view, then reset for the king's side
			if isCapture {
				SetUnsetUnsetPiece(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					1-c,
					FeatureIndex{piece, c, to, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
					FeatureIndex{capturedPiece, 1 - c, to, wKingSq, bKingSq},
				)
			} else if isCastling {
				val := 3
				if isQueenCast {
					val = -4
				}
				rookFrom := from + val
				val = -1
				if isQueenCast {
					val = 1
				}
				rookTo := to + val

				SetSetUnsetUnsetPiece(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					1-c,
					FeatureIndex{piece, c, to, wKingSq, bKingSq},
					FeatureIndex{Rook, c, rookTo, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
					FeatureIndex{Rook, c, rookFrom, wKingSq, bKingSq},
				)
			} else {
				SetUnsetPiece(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					1-c,
					FeatureIndex{piece, c, to, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
				)
			}
			e.ResetAccumulator(b, c)
		}
	} else {
		// Handle non-king moves
		movingPiece := piece
		if promPc != 0 {
			movingPiece = promPc
		}

		if isCapture {
			if enPass != 0 {
				epSquare := to - 8
				if c == Black {
					epSquare = to + 8
				}

				SetUnsetUnsetPieceBothColors(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					FeatureIndex{movingPiece, c, to, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
					FeatureIndex{Pawn, 1 - c, epSquare, wKingSq, bKingSq},
				)
			} else {
				SetUnsetUnsetPieceBothColors(
					&e.History[e.HistoryIndex-1],
					&e.History[e.HistoryIndex],
					FeatureIndex{movingPiece, c, to, wKingSq, bKingSq},
					FeatureIndex{piece, c, from, wKingSq, bKingSq},
					FeatureIndex{capturedPiece, 1 - c, to, wKingSq, bKingSq},
				)
			}
		} else {
			SetUnsetPieceBothColors(
				&e.History[e.HistoryIndex-1],
				&e.History[e.HistoryIndex],
				FeatureIndex{movingPiece, c, to, wKingSq, bKingSq},
				FeatureIndex{piece, c, from, wKingSq, bKingSq},
			)
		}
	}
}
