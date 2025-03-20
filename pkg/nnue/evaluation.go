package nnue

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
)

// Evaluator performs neural network evaluation of chess positions
type Evaluator struct {
	History                  []Accumulator
	HistoryIndex             int
	AccumulatorTable         *AccumulatorTable
	AccumulatorIsInitialized [2]bool
}

// NewEvaluator creates a new evaluator instance
func NewEvaluator() *Evaluator {
	evaluator := &Evaluator{
		History:          []Accumulator{{}},
		HistoryIndex:     0,
		AccumulatorTable: &AccumulatorTable{},
	}

	evaluator.AccumulatorTable.Reset()
	return evaluator
}

// Reset reinitializes the evaluator for a new board position
func (e *Evaluator) Reset(b *board.Board) {
	e.History = []Accumulator{{}}
	e.HistoryIndex = 0
	e.ResetAccumulator(b, White)
	e.ResetAccumulator(b, Black)
}

// ResetAccumulator reinitializes the accumulator for a specific perspective
func (e *Evaluator) ResetAccumulator(b *board.Board, color int) {
	e.AccumulatorTable.Use(color, b, e)
	e.AccumulatorIsInitialized[color] = true
}

var phaseValues = [5]float64{
	0.552938, 1.55294, 1.50862, 2.64379, 4.0,
}

func (e *Evaluator) Evaluate(b *board.Board) int {
	const (
		evaluationMgScalar = 1.5
		evaluationEgScalar = 1.15
		phaseSum           = 39.6684
	)

	phase := phaseSum

	phase -= float64((b.Bitboards[WP] | b.Bitboards[BP]).Count()) * phaseValues[Pawn]
	phase -= float64((b.Bitboards[WN] | b.Bitboards[BN]).Count()) * phaseValues[Knight]
	phase -= float64((b.Bitboards[WB] | b.Bitboards[BB]).Count()) * phaseValues[Bishop]
	phase -= float64((b.Bitboards[WR] | b.Bitboards[BR]).Count()) * phaseValues[Rook]
	phase -= float64((b.Bitboards[WQ] | b.Bitboards[BQ]).Count()) * phaseValues[Queen]

	phase /= phaseSum

	return int(
		(evaluationMgScalar - phase*(evaluationMgScalar-evaluationEgScalar)) * float64(
			e.eval(int(b.Side)),
		),
	)
}

// Evaluate computes a score for the current position
func (e *Evaluator) eval(activePlayer int) int {
	// if b != nil {
	// 	fmt.Println("reset")
	// 	e.Reset(b)
	// }

	// Get accumulator values
	accActive := e.History[e.HistoryIndex].Summation[activePlayer][:]
	accInactive := e.History[e.HistoryIndex].Summation[1-activePlayer][:]

	// Compute the output score
	var sum int32

	// Debug ReLU and dot product
	activeSum := int32(0)
	inactiveSum := int32(0)
	activeNonZeroCount := 0
	inactiveNonZeroCount := 0

	// Apply ReLU and compute dot product for active player
	for i := 0; i < HiddenSize; i++ {
		// ReLU: max(0, x)
		if accActive[i] > 0 {
			contribution := int32(accActive[i]) * int32(HiddenWeights[0][i])
			activeSum += contribution
			sum += contribution
			activeNonZeroCount++
		}
	}
	// Apply ReLU and compute dot product for inactive player
	for i := 0; i < HiddenSize; i++ {
		if accInactive[i] > 0 {
			contribution := int32(accInactive[i]) * int32(HiddenWeights[0][i+HiddenSize])
			inactiveSum += contribution
			sum += contribution
			inactiveNonZeroCount++
		}
	}
	// Add bias and scale
	sum = activeSum + inactiveSum + HiddenBias[0]

	result := int(
		float64(sum) / float64(InputWeightMultiplier) / float64(HiddenWeightMultiplier),
	)

	return result
}

// GetPieceValue returns each pieces value
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

// AddNewAccumulation adds a new state to the history
func (e *Evaluator) AddNewAccumulation() {
	e.HistoryIndex++

	// Expand history if needed
	if e.HistoryIndex >= len(e.History) {
		e.History = append(e.History, Accumulator{})
	}

	e.AccumulatorIsInitialized[White] = false
	e.AccumulatorIsInitialized[Black] = false
}

// PopAccumulation restores the previous state
func (e *Evaluator) PopAccumulation() {
	e.HistoryIndex--
	e.AccumulatorIsInitialized[White] = true
	e.AccumulatorIsInitialized[Black] = true
}

// ClearHistory resets the history to its initial state
func (e *Evaluator) ClearHistory() {
	e.History = []Accumulator{{}}
	e.HistoryIndex = 0
}

// SetPieceOnSquare updates the accumulator when adding/removing a piece
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

// SetPieceOnSquareAccumulator updates a single accumulator for a piece change
func (e *Evaluator) SetPieceOnSquareAccumulator(
	add bool,
	side, pieceType, pieceColor, square, kingSq int,
) {
	idx := Index(pieceType, pieceColor, square, side, kingSq)

	if !e.AccumulatorIsInitialized[side] {
		AddWeightsToAccumulator(
			add,
			idx,
			e.History[e.HistoryIndex-1].Summation[side][:],
			e.History[e.HistoryIndex].Summation[side][:],
		)
		e.AccumulatorIsInitialized[side] = true
	} else {
		AddWeightsToAccumulator(add, idx, e.History[e.HistoryIndex].Summation[side][:], e.History[e.HistoryIndex].Summation[side][:])
	}
}

// ProcessMove updates the accumulator for a chess move
func (e *Evaluator) ProcessMove(b *board.Board, m move.Move) {
	from := m.GetSource()
	to := m.GetTarget()
	piece := m.GetPiece()

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
