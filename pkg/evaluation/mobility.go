// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// MobilityEvaluation evaluates the mobility of the pieces
func (e *Evaluator) MobilityEvaluation(b *board.Board) (mg, eg int) {
	mg, eg = 0, 0
	sq := 0

	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		sq = knightBB.FirstOne()
		mobMg, mobEg := mobilityBonus(b, sq)
		mg, eg = mg+mobMg, eg+mobEg
	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq = bishopBB.FirstOne()
		mobMg, mobEg := mobilityBonus(b, sq)
		mg, eg = mg+mobMg, eg+mobEg
	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		sq = rookBB.FirstOne()
		mobMg, mobEg := mobilityBonus(b, sq)
		mg, eg = mg+mobMg, eg+mobEg
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		sq = queenBB.FirstOne()
		mobMg, mobEg := mobilityBonus(b, sq)
		mg, eg = mg+mobMg, eg+mobEg
	}

	return mg, eg
}

// MobilityBonus attaches bonuses for middlegame and endgame by piece type and Mobility
func mobilityBonus(b *board.Board, sq int) (mg, eg int) {
	if b.Bitboards[WN].Test(sq) {
		return mobilityBonusValues[0][0][mobility(b, sq)], mobilityBonusValues[1][0][mobility(b, sq)]
	}

	if b.Bitboards[WB].Test(sq) {
		return mobilityBonusValues[0][1][mobility(b, sq)], mobilityBonusValues[1][1][mobility(b, sq)]
	}

	if b.Bitboards[WR].Test(sq) {
		return mobilityBonusValues[0][2][mobility(b, sq)], mobilityBonusValues[1][2][mobility(b, sq)]
	}

	if b.Bitboards[WQ].Test(sq) {
		return mobilityBonusValues[0][3][mobility(b, sq)], mobilityBonusValues[1][3][mobility(b, sq)]
	}

	return 0, 0
}

// mobility is the number of attacked squares in the Mobility area. For queens squares
// defended by opponent knight, bishop or rook are ignored. For minor pieces squares
// occupied by our queen are ignored
func mobility(b *board.Board, sq int) int {
	if !b.Bitboards[WN].Test(sq) && !b.Bitboards[WB].Test(sq) && !b.Bitboards[WR].Test(sq) &&
		!b.Bitboards[WQ].Test(sq) {
		return 0
	}

	score := 0

	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if !mobilityArea(b, y*8+x) {
				continue
			}

			if b.Bitboards[WN].Test(sq) && KnightAttack(b, y*8+x, sq) > 0 &&
				!b.Bitboards[WQ].Test(y*8+x) {
				score++
			}
			if b.Bitboards[WB].Test(sq) && BishopXrayAttack(b, y*8+x, sq) > 0 &&
				!b.Bitboards[WQ].Test(y*8+x) {
				score++
			}
			if b.Bitboards[WR].Test(sq) && RookXrayAttack(b, y*8+x, sq) > 0 {
				score++
			}
			if b.Bitboards[WQ].Test(sq) && QueenAttack(b, y*8+x, sq) > 0 {
				score++
			}
		}
	}

	return score
}

// mobilityArea  do not include in mobility area squares protected by enemy pawns,
// or occupied by our blocked pawns or king. Pawns blocked or on ranks 2 and 3
// will be excluded from the mobility area. Also excludes blockers for king from
// mobility area - blockers for king can't really move until king moves (in most cases)
// so logic behind it is the same as behind excluding king square from mobility area.
func mobilityArea(b *board.Board, sq int) bool {
	if b.Bitboards[WK].Test(sq) {
		return false
	}

	if b.Bitboards[WQ].Test(sq) {
		return false
	}

	rank := sq / 8
	file := sq % 8

	if b.Bitboards[BP].Test((rank-1)*8+file-1) && file > 0 && rank > 0 {
		return false
	}

	if b.Bitboards[BP].Test((rank-1)*8+file+1) && file < 7 && rank > 0 {
		return false
	}

	if b.Bitboards[WP].Test(sq) &&
		((8-rank) < 4 || b.Occupancies[color.BOTH].Test((rank-1)*8+file)) {
		return false
	}

	mirror := b.Mirror()

	if blockersForKing(mirror, (7-rank)*8+file) > 0 {
		return false
	}

	return true
}

// blockersForKing returns if a particular piece on a particular square is a blocker
// for the king for a pin
func blockersForKing(b *board.Board, sq int) int {
	mirror := b.Mirror()
	rank := sq / 8

	if PinnedDirection(mirror, (7-rank)*8+(sq%8)) > 0 {
		return 1
	}

	return 0
}
