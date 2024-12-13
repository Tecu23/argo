package evaluation

import "github.com/Tecu23/argov2/pkg/board"

type mobilityEvaluator struct{}

func newMobilityEvaluator() *mobilityEvaluator {
	return &mobilityEvaluator{}
}

func (m *mobilityEvaluator) Evaluate(board *board.Board) int {
	// Count potential mobility using attack maps
	whiteMobility := m.calculateMobility(board, true)
	blackMobility := m.calculateMobility(board, false)

	return (whiteMobility - blackMobility) * 5
}

func (m *mobilityEvaluator) calculateMobility(board *board.Board, isWhite bool) int {
	mobility := 0

	// TODO: Use existing attack tables to generate mobility

	return mobility
}
