package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/evaluation/tables"
)

type positionalEvaluator struct {
	pawnTable   [64]int
	knightTable [64]int
	bishopTable [64]int
}

func newPositionalEvaluator() *positionalEvaluator {
	return &positionalEvaluator{
		pawnTable:   tables.PawnTable,
		knightTable: tables.KnightTable,
		bishopTable: tables.BishopTable,
	}
}

func (p *positionalEvaluator) Evaluate(board *board.Board) int {
	score := 0

	// Evaluate white pieces
	whitePawns := board.Bitboards[WP]
	for whitePawns != 0 {
		sq := whitePawns.FirstOne()
		score += p.pawnTable[sq]
	}

	whiteKnights := board.Bitboards[WN]
	for whiteKnights != 0 {
		sq := whiteKnights.FirstOne()
		score += p.pawnTable[sq]
	}

	whiteBishops := board.Bitboards[WB]
	for whiteBishops != 0 {
		sq := whiteBishops.FirstOne()
		score += p.pawnTable[sq]
	}

	// Evaluate white pieces
	blackPawns := board.Bitboards[BP]
	for blackPawns != 0 {
		sq := blackPawns.FirstOne()
		score -= p.pawnTable[sq]
	}

	blackKnights := board.Bitboards[BN]
	for blackKnights != 0 {
		sq := blackKnights.FirstOne()
		score -= p.pawnTable[sq]
	}

	blackBishops := board.Bitboards[BB]
	for blackBishops != 0 {
		sq := blackBishops.FirstOne()
		score -= p.pawnTable[sq]
	}

	return score
}
