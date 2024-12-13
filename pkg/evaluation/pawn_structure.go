package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
)

type pawnStructureEvaluator struct{}

func newPawnStructureEvaluator() *pawnStructureEvaluator {
	return &pawnStructureEvaluator{}
}

func (p *pawnStructureEvaluator) Evaluate(board *board.Board) int {
	score := 0

	// TODO: Should evaluate doubled pawns, isolated pawns and pawn structure

	return score
}
