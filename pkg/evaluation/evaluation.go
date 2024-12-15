package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
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

func GetPieceValue(piece int) int {
	switch piece {
	case WP, BP:
		return 100 // Pawn value
	case WN, BN:
		return 320 // Knight value
	case WB, BB:
		return 330 // Bishop value
	case WR, BR:
		return 500 // Rook value
	case WQ, BQ:
		return 900 // Queen value
	case WK, BK:
		return 20000 // King value
	default:
		return 0
	}
}
