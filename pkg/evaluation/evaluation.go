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
	mg, eg := e.EvaluateOneSide(board, false)

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
	mg, eg := 0, 0
	mirror := board.Mirror()

	wMatMg, wMatEg := e.MaterialEvaluation(board)
	bMatMg, bMatEg := e.MaterialEvaluation(mirror)

	mg = mg + wMatMg - bMatMg
	eg = eg + wMatEg - bMatEg

	wPosMg, wPosEg := e.PositionalEvaluation(board)
	bPosMg, bPosEg := e.PositionalEvaluation(mirror)

	mg = mg + wPosMg - bPosMg
	eg = eg + wPosEg - bPosEg

	imb := e.ImbalanceEvaluation(board, mirror)
	mg, eg = mg+imb, eg+imb

	wPawnMg, wPawnEg := e.PawnsEvaluation(board)
	bPawnMg, bPawnEg := e.PawnsEvaluation(mirror)

	mg = mg + wPawnMg - bPawnMg
	eg = eg + wPawnEg - bPawnEg

	wPiecesMg, wPiecesEg := e.PiecesEvaluation(board)
	bPiecesMg, bPiecesEg := e.PiecesEvaluation(mirror)

	mg = mg + wPiecesMg - bPiecesMg
	eg = eg + wPiecesEg - bPiecesEg

	wMobMg, wMobEg := e.MobilityEvaluation(board)
	bMobMg, bMobEg := e.MobilityEvaluation(mirror)

	mg = mg + wMobMg - bMobMg
	eg = eg + wMobEg - bMobEg

	wThMg, wThEg := e.ThreatsEvaluation(board)
	bThMg, bThEg := e.ThreatsEvaluation(mirror)

	mg = mg + wThMg - bThMg
	eg = eg + wThEg - bThEg

	wPassMg, wPassEg := e.PassedPawnEvaluation(board)
	bPassMg, bPassEg := e.PassedPawnEvaluation(mirror)

	mg = mg + wPassMg - bPassMg
	eg = eg + wPassEg - bPassEg

	space := Space(board) - Space(mirror)
	mg = mg + space

	wKingMg, wKingEg := e.KingEvaluation(board)
	bKingMg, bKingEg := e.KingEvaluation(mirror)

	mg = mg + wKingMg - bKingMg
	eg = eg + wKingEg - bKingEg

	if !noWinnable {
		wWinMg, wWinEg := e.WinnableEvaluation(board, mg, eg)

		mg = mg + wWinMg
		eg = eg + wWinEg
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
	if eg == -1 {
		_, egT := e.EvaluateOneSide(b, false)
		return egT
	}

	mirror := b.Mirror()
	var pos_w *board.Board
	var pos_b *board.Board
	if eg > 0 {
		pos_w = b
		pos_b = mirror
	} else {
		pos_w = mirror
		pos_b = b
	}

	sf := 64

	pc_w := pos_w.Bitboards[WP].Count()
	pc_b := pos_b.Bitboards[WP].Count()

	qc_w := pos_w.Bitboards[WQ].Count()
	qc_b := pos_b.Bitboards[WQ].Count()

	bc_w := pos_w.Bitboards[WB].Count()
	bc_b := pos_b.Bitboards[WB].Count()

	nc_w := pos_w.Bitboards[WN].Count()
	nc_b := pos_b.Bitboards[WN].Count()

	npm_w := nonPawnMaterial(pos_w, color.WHITE)
	npm_b := nonPawnMaterial(pos_b, color.WHITE)

	bishopValueMg := 825
	bishopValueEg := 915
	rookValueMg := 1276

	if pc_w == 0 && npm_w-npm_b <= bishopValueMg {
		if npm_w < rookValueMg {
			sf = 0
		} else if npm_b <= bishopValueEg {
			sf = 4
		} else {
			sf = 14
		}
	}

	if sf == 64 {
		ob := b.OppositeBishops()

		if ob && npm_w == bishopValueMg && npm_b == bishopValueMg {
			sf = 22 + 4*pos_w.CandidatePassed() // Get passed pawns for white pos
		} else if ob {
			sf = 22 + 3*PieceCount(pos_w)
		} else {
			if npm_w == rookValueMg && npm_b == rookValueMg && pc_w-pc_b <= 1 {
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

			if qc_w+qc_b == 1 {
				if qc_w == 1 {
					sf = 37 + 3*(bc_b+nc_b)
				} else {
					sf = 37 + 3*(bc_w+nc_w)
				}
			} else {
				sf = min(sf, 36+7*pc_w)
			}
		}
	}

	return sf
}

func PieceCount(b *board.Board) int {
	return b.Occupancies[color.WHITE].Count()
}
