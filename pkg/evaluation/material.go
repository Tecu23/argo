package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

type materialEvaluator struct{}

func newMaterialEvaluator() *materialEvaluator {
	return &materialEvaluator{}
}

func (m *materialEvaluator) Evaluate(board *board.Board) int {
	var score int

	// White pieces evaluation
	score += board.Bitboards[WP].Count() * PawnValue
	score += board.Bitboards[WN].Count() * KnightValue
	score += board.Bitboards[WB].Count() * BishopValue
	score += board.Bitboards[WR].Count() * RookValue
	score += board.Bitboards[WQ].Count() * QueenValue

	// Black pieces evaluation
	score -= board.Bitboards[BP].Count() * PawnValue
	score -= board.Bitboards[BN].Count() * KnightValue
	score -= board.Bitboards[BB].Count() * BishopValue
	score -= board.Bitboards[BR].Count() * RookValue
	score -= board.Bitboards[BQ].Count() * QueenValue

	return score
}
