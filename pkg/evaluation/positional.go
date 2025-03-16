// Package evaluation is responsible with evaluating the current board position
package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// PositionalEvaluation returns the positional evaluation of the board pieces
// using piece square table bonuses. It returns bothe middlegame and endgame evaluations
func (e *Evaluator) PositionalEvaluation(b *board.Board) (mg, eg int) {
	pawnBB := b.Bitboards[WP]
	for pawnBB != 0 {
		sq := pawnBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtPawnBonus[0][7-rank][file]
		eg += psqtPawnBonus[1][7-rank][file]
	}

	knightBB := b.Bitboards[WN]
	for knightBB != 0 {
		sq := knightBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][0][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][0][7-rank][min(file, 7-file)]
	}

	bishopBB := b.Bitboards[WB]
	for bishopBB != 0 {
		sq := bishopBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][1][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][1][7-rank][min(file, 7-file)]
	}

	rookBB := b.Bitboards[WR]
	for rookBB != 0 {
		sq := rookBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][2][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][2][7-rank][min(file, 7-file)]
	}

	queenBB := b.Bitboards[WQ]
	for queenBB != 0 {
		sq := queenBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][3][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][3][7-rank][min(file, 7-file)]
	}

	kingBB := b.Bitboards[WK]
	for kingBB != 0 {
		sq := kingBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][4][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][4][7-rank][min(file, 7-file)]
	}

	return mg, eg
}
