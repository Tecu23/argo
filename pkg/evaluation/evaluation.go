package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
)

type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Evaluate(board *board.Board) int {
	return 0
}
