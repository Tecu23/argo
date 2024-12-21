package evaluation

import (
	"github.com/Tecu23/argov2/pkg/board"
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

	// score := 0
	// score += e.material.Evaluate(board)
	// score += e.position.Evaluate(board)
	// score += e.mobility.Evaluate(board)
	// score += e.pawnStructure.Evaluate(board)
	//
	// if board.Side == color.BLACK {
	// 	score = -score
	// }
	//
	// return score

	a := phase(board)

	return a
}

func GetPieceValue(piece int) int {
	switch piece {
	case WP, BP:
		return PawnValue
	case WN, BN:
		return KnightValue
	case WB, BB:
		return BishopValue
	case WR, BR:
		return RookValue
	case WQ, BQ:
		return QueenValue
	case WK, BK:
		return KingValue
	default:
		return 0
	}
}

// IsEndgame returns true if the position is in endgame phase
func (e *Evaluator) IsEndgame(b *board.Board) bool {
	// Simple material-based endgame detection
	// Consider it endgame if either side has no queens or
	// if total material (excluding kings and pawns) is low
	const endgameMaterialThreshold = 1300 // roughly equivalent to a rook + bishop

	materialWithoutKingsPawns := 0

	// Count material for both sides, excluding kings and pawns
	for sq := 0; sq < 64; sq++ {
		piece := b.GetPieceAt(sq)
		if piece == Empty {
			continue
		}

		pieceType := piece % 6 // Get piece type without color
		if pieceType != King && pieceType != Pawn {
			materialWithoutKingsPawns += GetPieceValue(piece)
		}
	}

	return materialWithoutKingsPawns <= endgameMaterialThreshold
}
