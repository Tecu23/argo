package engine

import (
	"context"
	"math/rand"
	"time"

	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/move"
)

// search performs the actual search logic
func (e *Engine) search(ctx context.Context, b board.Board, tm *timeManager) SearchInfo {
	// Generate all legal moves
	moves := b.GenerateMoves()
	legalMoves := []move.Move{}

	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !b.MakeMove(mv, board.AllMoves) {
			continue
		}
		b.TakeBack(copyB)
		legalMoves = append(legalMoves, mv)
	}

	if len(legalMoves) == 0 {
		return SearchInfo{}
	}

	// Select random move
	randomIndex := rand.Intn(len(legalMoves))
	bestMove := legalMoves[randomIndex]

	// Store in main line
	e.mainLine.moves = []move.Move{bestMove}
	e.mainLine.score = 0
	e.mainLine.depth = 1
	e.mainLine.nodes = 1

	// Simulate "thinking"
	select {
	case <-time.After(500 * time.Millisecond):

		// Send Progress info
		if e.progress != nil {
			e.progress(e.createSearchInfo())
		}
	case <-ctx.Done():
		return e.createSearchInfo()
	}

	return e.createSearchInfo()
}

func (e *Engine) minimax(board board.Board, depth int) int {
	return 0
}

func (e *Engine) alphaBeta(board board.Board, depth, alpha, beta int) int {
	return 0
}
