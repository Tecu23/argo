// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
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

func nonPawnMaterial(b *board.Board, clr color.Color) int {
	score := 0

	if clr == color.WHITE {
		for pc := WN; pc < WK; pc++ {
			bb := b.Bitboards[pc]
			for bb.Count() > 0 {
				_ = bb.FirstOne()

				switch pc {
				case BN, WN:
					score += KnightBonus[0]
				case BB, WB:
					score += BishopBonus[0]
				case BR, WR:
					score += RookBonus[0]
				case BQ, WQ:
					score += QueenBonus[0]
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
					score += KnightBonus[0]
				case BB, WB:
					score += BishopBonus[0]
				case BR, WR:
					score += RookBonus[0]
				case BQ, WQ:
					score += QueenBonus[0]
				}
			}
		}
	}

	return score
}
