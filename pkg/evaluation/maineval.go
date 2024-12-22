package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

func MainEvaluation(b *board.Board) int {
	mg := MiddleGameEvaluation(b, false)
	eg := EndGameEvaluation(b, false)
	p := Phase(b)
	r50 := Rule50(b)

	eg = eg * ScaleFactor(b, eg) / 64

	v := (((mg*p + ((eg * (128 - p)) << 0)) / 128) << 0)

	v = ((v / 16) << 0) * 16
	v += Tempo(b)
	v = (v * (100 - r50) / 100) << 0

	return v
}

func MiddleGameEvaluation(b *board.Board, noWinnable bool) int {
	score := 0
	mirror := b.Mirror()

	score += PieceValueMg(b) - PieceValueMg(mirror)
	score += PsqtMg(b) - PsqtMg(mirror)
	score += ImbalanceTotal(b, mirror)
	score += PawnsMg(b) - PawnsMg(mirror)
	score += PiecesMg(b) - PiecesMg(mirror)
	score += MobilityMg(b) - MobilityMg(mirror)
	score += ThreatsMg(b) - ThreatsMg(mirror)
	score += PassedMg(b) - PassedMg(mirror)
	score += Space(b) - Space(mirror)
	score += KingMg(b) - KingMg(mirror)

	if !noWinnable {
		score += WinnableTotalMg(b, score)
	}

	return score
}

func EndGameEvaluation(b *board.Board, noWinnable bool) int {
	score := 0
	mirror := b.Mirror()

	score += PieceValueEg(b) - PieceValueEg(mirror)
	score += PsqtEg(b) - PieceValueEg(mirror)
	score += ImbalanceTotal(b, mirror)
	score += PawnsEg(b) - PawnsEg(mirror)
	score += PiecesEg(b) - PiecesEg(mirror)
	score += MobilityEg(b) - MobilityEg(mirror)
	score += ThreatsEg(b) - ThreatsEg(mirror)
	score += PassedEg(b) - PassedEg(mirror)
	score += KingEg(b) - KingEg(mirror)

	if !noWinnable {
		score += WinnableTotalEg(b, score)
	}

	return score
}

func Tempo(b *board.Board) int {
	factor := 1
	if b.Side == color.BLACK {
		factor = -1
	}
	return 28 * factor
}

func Rule50(b *board.Board) int {
	return int(b.Rule50)
}

func Phase(b *board.Board) int {
	midgameLimit := 15258
	endgameLimit := 3915

	//
	npm := nonPawnMaterial(b, color.WHITE) + nonPawnMaterial(b, color.BLACK)
	npm = max(endgameLimit, min(npm, midgameLimit))

	return (((npm - endgameLimit) * 128) / (midgameLimit - endgameLimit)) << 0
}

// TODO: Finish and test this
func ScaleFactor(b *board.Board, eg int) int {
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
