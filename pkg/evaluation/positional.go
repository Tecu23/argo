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
	var pawnMg, pawnEg int
	for pawnBB != 0 {
		sq := pawnBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtPawnBonus[0][7-rank][file]
		eg += psqtPawnBonus[1][7-rank][file]

		pawnMg, pawnEg = e.evaluatePawn(b, sq)
		mg, eg = mg+pawnMg, eg+pawnEg
	}

	knightBB := b.Bitboards[WN]
	var knightMg, knightEg int
	for knightBB != 0 {
		sq := knightBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][0][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][0][7-rank][min(file, 7-file)]

		knightMg, knightEg = e.evaluateKnight(b, sq)
		mg, eg = mg+knightMg, eg+knightEg
	}

	bishopBB := b.Bitboards[WB]
	var bishopMg, bishopEg int
	for bishopBB != 0 {
		sq := bishopBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][1][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][1][7-rank][min(file, 7-file)]

		bishopMg, bishopEg = e.evaluateBishop(b, sq)
		mg, eg = mg+bishopMg, eg+bishopEg
	}

	rookBB := b.Bitboards[WR]
	var rookMg, rookEg int
	for rookBB != 0 {
		sq := rookBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][2][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][2][7-rank][min(file, 7-file)]

		rookMg, rookEg = e.evaluateRook(b, sq)
		mg, eg = mg+rookMg, eg+rookEg
	}

	queenBB := b.Bitboards[WQ]
	var queenMg, queenEg int
	for queenBB != 0 {
		sq := queenBB.FirstOne()
		rank := sq / 8
		file := sq % 8

		mg += psqtBonus[0][3][7-rank][min(file, 7-file)]
		eg += psqtBonus[1][3][7-rank][min(file, 7-file)]

		queenMg, queenEg = e.evaluateQueen(b, sq)
		mg, eg = mg+queenMg, eg+queenEg
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
