// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// Evaluator represents the main evaluator of the position
type Evaluator struct{}

// NewEvaluator creates a new evaluator
func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(board *board.Board) int {
	wMG, wEG := e.EvaluateOneSide(board, false)
	bMG, bEG := e.EvaluateOneSide(board.Mirror(), false)

	mg := wMG - bMG
	eg := wEG - bEG

	p := e.Phase(board)
	r50 := e.Rule50(board)

	eg = eg * e.ScaleFactor(board, eg) / 64

	v := (((mg*p + ((eg * (128 - p)) << 0)) / 128) << 0)

	v = ((v / 16) << 0) * 16
	v += e.Tempo(board)
	v = (v * (100 - r50) / 100) << 0

	return v
}

// Evaluate returns the current position evaluation
func (e *Evaluator) EvaluateOneSide(board *board.Board, noWinnable bool) (int, int) {
	mg, eg := e.MaterialEvaluation(board)

	posMg, posEg := e.PositionalEvaluation(board)
	mg, eg = mg+posMg, eg+posEg

	imb := imbalance(board)
	mg, eg = mg+imb, eg+imb

	pawnMg, pawnEg := e.PawnsEvaluation(board)
	mg, eg = mg+pawnMg, eg+pawnEg

	piecesMg, piecesEg := e.PiecesEvaluation(board)
	mg, eg = mg+piecesMg, eg+piecesEg

	mobMg, mobEg := e.MobilityEvaluation(board)
	mg, eg = mg+mobMg, eg+mobEg

	thMg, thEg := e.ThreatsEvaluation(board)
	mg, eg = mg+thMg, eg+thEg

	passMg, passEg := e.PassedPawnEvaluation(board)
	mg, eg = mg+passMg, eg+passEg

	space := Space(board)
	mg, eg = mg+space, eg+space

	kingMg, kingEg := e.KingEvaluation(board)
	mg, eg = mg+kingMg, eg+kingEg

	if !noWinnable {
		wMG, wEG := e.WinnableEvaluation(board, mg, eg)
		mg, eg = mg+wMG, eg+wEG
	}

	return mg, eg
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

func (e *Evaluator) Tempo(b *board.Board) int {
	factor := 1
	if b.Side == color.BLACK {
		factor = -1
	}
	return 28 * factor
}

func (e *Evaluator) Rule50(b *board.Board) int {
	return int(b.Rule50)
}

func (e *Evaluator) Phase(b *board.Board) int {
	midgameLimit := 15258
	endgameLimit := 3915

	//
	npm := nonPawnMaterial(b, color.WHITE) + nonPawnMaterial(b, color.BLACK)
	npm = max(endgameLimit, min(npm, midgameLimit))

	return (((npm - endgameLimit) * 128) / (midgameLimit - endgameLimit)) << 0
}

func (e *Evaluator) ScaleFactor(b *board.Board, eg int) int {
	sf := 64

	attackingSide := color.WHITE
	opponentSide := color.BLACK

	if eg < 0 {
		attackingSide = color.BLACK
		opponentSide = color.WHITE
	}

	pawnCountWhite := b.GetPieceCountForSide(Pawn, attackingSide)
	knightCountWhite := b.GetPieceCountForSide(Knight, attackingSide)
	bishopCountWhite := b.GetPieceCountForSide(Bishop, attackingSide)
	queenCountWhite := b.GetPieceCountForSide(Queen, attackingSide)

	pawnCountBlack := b.GetPieceCountForSide(Pawn, opponentSide)
	knightCountBlack := b.GetPieceCountForSide(Knight, opponentSide)
	bishopCountBlack := b.GetPieceCountForSide(Bishop, opponentSide)
	queenCountBlack := b.GetPieceCountForSide(Queen, opponentSide)

	npmWhite := nonPawnMaterial(b, attackingSide)
	npmBlack := nonPawnMaterial(b, opponentSide)

	bishopValueMg := 825
	bishopValueEg := 915
	rookValueMg := 1276

	if pawnCountWhite == 0 && npmWhite-npmBlack <= bishopValueMg {
		if npmWhite < rookValueMg {
			sf = 0
		} else if npmBlack <= bishopValueEg {
			sf = 4
		} else {
			sf = 14
		}
	}

	if sf == 64 {
		ob := b.OppositeBishops()

		if ob && npmWhite == bishopValueMg && npmBlack == bishopValueMg {
			sf = 22 + 4*b.CandidatePassed(attackingSide) // Get passed pawns for white pos
		} else if ob {
			sf = 22 + 3*b.PieceCount(attackingSide)
		} else {
			if npmWhite == rookValueMg && npmBlack == rookValueMg && pawnCountWhite-pawnCountBlack <= 1 {
				pawnKingBlack := 0
				pawnCountWhiteFlank := []int{0, 0}

				for sq := A8; sq <= H1; sq++ {
					x := sq / 8
					y := sq % 8
					if b.GetPieceAt(sq) == WP {
						if x < 4 {
							pawnCountWhiteFlank[1] = 1
						} else {
							pawnCountWhiteFlank[0] = 1
						}
					}

					if b.GetPieceAt(sq) == BK {
						for ix := -1; ix <= 1; ix++ {
							for iy := -1; iy <= 1; iy++ {
								if b.GetPieceAt((x+ix)*8+(y+iy)) == BP {
									pawnKingBlack = 1
								}
							}
						}
					}
				}

				if pawnCountWhiteFlank[0] != pawnCountWhiteFlank[1] && pawnKingBlack != 0 {
					return 36
				}
			}

			if queenCountWhite+queenCountBlack == 1 {
				if queenCountWhite == 1 {
					sf = 37 + 3*(bishopCountBlack+knightCountBlack)
				} else {
					sf = 37 + 3*(bishopCountWhite+knightCountWhite)
				}
			} else {
				sf = min(sf, 36+7*pawnCountWhite)
			}
		}
	}
	return sf
}
