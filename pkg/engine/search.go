package engine

import (
	"context"

	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/move"
)

const (
	MaxDepth  = 64
	Infinity  = 50_000
	MateScore = 49_000
	MateDepth = 48_000
)

func (e *Engine) search(ctx context.Context, b board.Board, tm *timeManager) SearchInfo {
	e.nodes = 0
	var bestMove move.Move

	// CHANGED: Added variables to track best move per iteration
	var bestMoveAtDepth move.Move
	var bestScoreAtDepth int

	for depth := 1; depth <= MaxDepth; depth++ {
		// CHANGED: Initialize best score for this iteration
		bestScoreAtDepth = -Infinity
		bestMoveAtDepth = move.Move(0)

		moves := b.GenerateMoves()

		// CHANGED: Root move search loop
		for _, mv := range moves {
			copyB := b.CopyBoard()
			if !b.MakeMove(mv, board.AllMoves) {
				continue
			}

			// CHANGED: Proper negamax call with correct bounds
			score := -e.alphaBeta(ctx, b, depth-1, -Infinity, -bestScoreAtDepth)

			b.TakeBack(copyB)

			if ctx.Err() != nil {
				if bestMove != move.Move(0) {
					return e.createSearchInfo()
				}
			}

			// CHANGED: Update best move if score is better
			if score > bestScoreAtDepth {
				bestScoreAtDepth = score
				bestMoveAtDepth = mv
			}
		}

		// CHANGED: Update best move after iteration is complete
		bestMove = bestMoveAtDepth

		e.mainLine.moves = []move.Move{bestMove}
		e.mainLine.score = bestScoreAtDepth // CHANGED: Use score from iteration
		e.mainLine.depth = depth
		e.mainLine.nodes = e.nodes

		if e.progress != nil {
			e.progress(e.createSearchInfo())
		}

		tm.updateSearch(bestScoreAtDepth, bestMove)

		if tm.shouldStop() || ctx.Err() != nil {
			break
		}
	}

	return e.createSearchInfo()
}

func (e *Engine) alphaBeta(ctx context.Context, b board.Board, depth, alpha, beta int) int {
	if ctx.Err() != nil {
		return 0
	}

	e.nodes++

	// CHANGED: Better base cases
	if depth == 0 {
		return e.quiescence(ctx, b, alpha, beta)
	}

	// CHANGED: Check for draws first
	if b.IsStalemate() || b.IsInsufficientMaterial() {
		return 0
	}

	// CHANGED: Proper mate detection
	if b.IsCheckmate() {
		return -MateScore + int(e.nodes) // Shorter mates are preferred
	}

	moves := b.GenerateMoves()
	if len(moves) == 0 {
		if b.InCheck() {
			return -MateScore + int(e.nodes)
		}
		return 0 // Stalemate
	}

	// CHANGED: Proper move loop with alpha-beta
	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !b.MakeMove(mv, board.AllMoves) {
			continue
		}

		score := -e.alphaBeta(ctx, b, depth-1, -beta, -alpha)

		b.TakeBack(copyB)

		if ctx.Err() != nil {
			return 0
		}

		if score >= beta {
			return beta // Beta cutoff
		}
		if score > alpha {
			alpha = score // Alpha update
		}
	}

	return alpha
}

func (e *Engine) quiescence(ctx context.Context, b board.Board, alpha, beta int) int {
	if ctx.Err() != nil {
		return 0
	}

	e.nodes++

	// CHANGED: Added stand-pat evaluation
	standPat := e.evaluator.Evaluate(&b)

	if standPat >= beta {
		return beta
	}
	if alpha < standPat {
		alpha = standPat
	}

	// CHANGED: Use GenerateCaptures instead of GenerateMoves
	moves := b.GenerateCaptures()

	for _, mv := range moves {
		copyB := b.CopyBoard()
		// CHANGED: Use OnlyCaptures flag
		if !b.MakeMove(mv, board.OnlyCaptures) {
			continue
		}

		score := -e.quiescence(ctx, b, -beta, -alpha)

		b.TakeBack(copyB)

		if ctx.Err() != nil {
			return 0
		}

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	return alpha
}
