package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
)

type Evaluator struct {
	material      *materialEvaluator
	mobility      *mobilityEvaluator
	position      *positionalEvaluator
	pawnStructure *pawnStructureEvaluator
}

func NewEvaluator() *Evaluator {
	return &Evaluator{
		material:      newMaterialEvaluator(),
		position:      newPositionalEvaluator(),
		mobility:      newMobilityEvaluator(),
		pawnStructure: newPawnStructureEvaluator(),
	}
}

func (e *Evaluator) Evaluate(board *board.Board) int {
	// Early game termination checks would go here
	// if board.isCheckMate() etc.

	score := 0
	score += e.material.Evaluate(board)
	score += e.position.Evaluate(board)
	score += e.mobility.Evaluate(board)
	score += e.pawnStructure.Evaluate(board)

	if board.Side == color.BLACK {
		score = -score
	}

	return score
}
