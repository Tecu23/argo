package evaluation

import (
	"fmt"

	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

var (
	PawnBonus   = [2]int{124, 206}
	KnightBonus = [2]int{781, 854}
	BishopBonus = [2]int{825, 915}
	RookBonus   = [2]int{1276, 1380}
	QueenBonus  = [2]int{2538, 2682}
)

func rule50(b *board.Board) int {
	return int(b.Rule50)
}

func phase(b *board.Board) int {
	midgameLimit := 15258
	endgameLimit := 3915

	//
	npm := nonPawnMaterial(b, color.WHITE) + nonPawnMaterial(b, color.BLACK)

	fmt.Println(npm)

	npm = max(endgameLimit, min(npm, midgameLimit))

	return (((npm - endgameLimit) * 128) / (midgameLimit - endgameLimit)) << 0
}

func nonPawnMaterial(b *board.Board, clr color.Color) int {
	score := 0

	if clr == color.WHITE {
		for pc := WN; pc < WK; pc++ {
			bb := b.Bitboards[pc]
			for bb.Count() > 0 {
				_ = bb.FirstOne()

				switch pc {
				case BN, WN:
					score += KnightBonus[1]
				case BB, WB:
					score += BishopBonus[1]
				case BR, WR:
					score += RookBonus[1]
				case BQ, WQ:
					score += QueenBonus[1]
				}
			}
		}
	} else {
		for pc := BN; pc < BK; pc++ {
			bb := b.Bitboards[pc]
			for bb.Count() > 0 {
				_ = bb.FirstOne()

				switch pc {
				case BN, WN:
					score += KnightBonus[1]
				case BB, WB:
					score += BishopBonus[1]
				case BR, WR:
					score += RookBonus[1]
				case BQ, WQ:
					score += QueenBonus[1]
				}
			}
		}
	}

	return score
}

func scaleFactor(b *board.Board, eg int) int {
	sf := 64

	attackingSide := color.WHITE
	opponentSide := color.BLACK

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
			sf = 22 + 4*candidatePassed(pos_w) // Get passed pawns for white pos
		} else if ob {
			sf == 22+3*pieceCount(pos_W)
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
