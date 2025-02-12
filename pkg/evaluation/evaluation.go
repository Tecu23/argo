package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(board *board.Board) int {
	return MainEvaluation(board)
}

func (e *Evaluator) IsEndgame(b *board.Board) bool {
	return false
}

func GetPieceValue(piece int) int {
	PawnBonus := 124
	KnightBonus := 781
	BishopBonus := 825
	RookBonus := 1276
	QueenBonus := 2538

	value := 0

	if piece == BP || piece == WP {
		value = PawnBonus
	} else if piece == BN || piece == WN {
		value = KnightBonus
	} else if piece == BB || piece == WB {
		value = BishopBonus
	} else if piece == BR || piece == WR {
		value = RookBonus
	} else if piece == BQ || piece == WQ {
		value = QueenBonus
	}

	return value
}
