// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// MaterialEvaluation returns the sum of all the pieces multiplied by each piece bonus
func (e *Evaluator) MaterialEvaluation(b *board.Board) (mg int, eg int) {
	mg = 0
	eg = 0

	pawnCount := b.Bitboards[WP].Count()
	knightCount := b.Bitboards[WN].Count()
	bishopCount := b.Bitboards[WB].Count()
	rookCount := b.Bitboards[WR].Count()
	queenCount := b.Bitboards[WQ].Count()

	mg = pawnCount*pawnBonusMG + knightCount*knightBonusMG + bishopCount*bishopBonusMG + rookCount*rookBonusMG + queenCount*queenBonusMG
	eg = pawnCount*pawnBonusEG + knightCount*knightBonusEG + bishopCount*bishopBonusEG + rookCount*rookBonusEG + queenCount*queenBonusEG

	return mg, eg
}
