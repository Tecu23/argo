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

// search performs the actual search logic
func (e *Engine) search(ctx context.Context, b *board.Board, tm *timeManager) SearchInfo {
	e.nodes = 0
	e.tt.NewSearch()

	var bestMove move.Move
	var bestScore int

	// Engine will only loop up to depth limit, regardless of time
	maxDepth := MaxDepth
	if tm.limits.Depth > 0 {
		maxDepth = tm.limits.Depth
	}

	// Iterative deepeing
	for depth := 1; depth <= maxDepth; depth++ {
		score, mv := e.searchRoot(ctx, b, depth)

		if mv != move.NoMove {
			// Store best move and score
			bestMove = mv
			bestScore = score

			// Update search info
			e.mainLine.moves = []move.Move{bestMove}
			e.mainLine.score = bestScore
			e.mainLine.depth = depth
			e.mainLine.nodes = e.nodes
		}

		// Report progress
		if e.progress != nil {
			e.progress(e.createSearchInfo())
		}

		// Check if we should stop
		if tm.IsDone() || ctx.Err() != nil {
			break
		}

		// Update time manager
		tm.OnNodesChanged(int(e.nodes))
	}

	searchInfo := e.createSearchInfo()

	// Ensure we have a move to return
	if len(searchInfo.MainLine) == 0 && bestMove != move.NoMove {
		searchInfo.MainLine = []move.Move{bestMove}
	}

	return e.createSearchInfo()
}

// searchRoot performs alpha-beta search at the root level
func (e *Engine) searchRoot(ctx context.Context, b *board.Board, depth int) (int, move.Move) {
	alpha := -Infinity
	beta := Infinity
	var bestMove move.Move

	// Generate moves at root
	moves := b.GenerateMoves()

	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !b.MakeMove(mv, board.AllMoves) {
			continue
		}

		// Search this position
		score := -e.alphaBeta(ctx, b, depth-1, -beta, -alpha)

		b.TakeBack(copyB)

		// Check for search abort
		if ctx.Err() != nil {
			return 0, move.NoMove
		}

		// Update best score if we found a better score
		if score > alpha {
			alpha = score
			bestMove = mv
		}

	}

	return alpha, bestMove
}

// alphaBeta performs the main alpha-beta search
func (e *Engine) alphaBeta(ctx context.Context, b *board.Board, depth, alpha, beta int) int {
	// Increment node counter
	e.nodes++

	// TT Lookup
	hash := b.Hash()
	if entry, ok := e.tt.Probe(hash); ok {
		if entry.Depth >= depth {
			switch entry.Flag {
			case TTExact:
				return entry.Score
			case TTAlpha:
				if entry.Score <= alpha {
					return alpha
				}
			case TTBeta:
				if entry.Score >= beta {
					return beta
				}
			}
		}
	}

	// exit early if time exceeeded
	if (e.nodes & 1023) == 0 {
		if ctx.Err() != nil {
			// This ensures the move won't be selected
			// as it will always be worse than any real evaluation
			if depth&1 == 0 {
				return -Infinity // for maximizing player
			}
			return Infinity // for minimizing player
		}
	}

	// Check for terminal positions
	if b.IsCheckmate() {
		return -MateScore + int(e.nodes) // Prefer shorter mates
	}

	if b.IsStalemate() || b.IsInsufficientMaterial() {
		return 0
	}

	// Base case: evaluate leaf nodes
	if depth <= 0 {
		return e.quiescence(ctx, b, alpha, beta)
	}

	// Generate moves
	moves := b.GenerateMoves()
	hasLegalMoves := false
	var bestMove move.Move
	bestScore := -Infinity

	// Search all moves {
	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !b.MakeMove(mv, board.AllMoves) {
			continue
		}

		hasLegalMoves = true
		score := -e.alphaBeta(ctx, b, depth-1, -beta, -alpha)
		b.TakeBack(copyB)

		// Check for search abort
		if ctx.Err() != nil {
			// This ensures the move won't be selected
			// as it will always be worse than any real evaluation
			if depth&1 == 0 {
				return -Infinity // for maximizing player
			}
			return Infinity // for minimizing player
		}

		if score > bestScore {
			bestScore = score
			bestMove = mv
		}

		alpha = max(alpha, score)
		if alpha >= beta {
			e.tt.Store(hash, beta, depth, TTBeta, mv)
			return beta
		}
	}

	// Check for chechmate/stalemate
	if !hasLegalMoves {
		if b.InCheck() {
			return -MateScore + int(e.nodes)
		}
		return 0
	}

	// Store position in TT
	var flag TTFlag
	if bestScore <= alpha {
		flag = TTAlpha
	} else {
		flag = TTExact
	}
	e.tt.Store(hash, bestScore, depth, flag, bestMove)

	return bestScore
}

// quiescence performs capture-only search to reach quiet position
func (e *Engine) quiescence(ctx context.Context, b *board.Board, alpha, beta int) int {
	e.nodes++

	// TT Lookup for quiescence
	hash := b.Hash()
	if entry, ok := e.tt.Probe(hash); ok {
		if entry.Flag == TTExact {
			return entry.Score
		}
	}

	if (e.nodes & 4095) == 0 {
		if ctx.Err() != nil {
			// Use same logic as alpha-beta for consistency
			if e.nodes&1 == 0 {
				return -Infinity
			}
			return Infinity
		}
	}

	// Stand-pat score
	score := e.evaluator.Evaluate(b)
	if score >= beta {
		return beta
	}

	if score > alpha {
		alpha = score
	}

	// Generate captures
	moves := b.GenerateCaptures()

	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !b.MakeMove(mv, board.OnlyCaptures) {
			continue
		}

		score := -e.quiescence(ctx, b, -beta, -alpha)

		b.TakeBack(copyB)

		if ctx.Err() != nil {
			if e.nodes&1 == 0 {
				return -Infinity
			}
			return Infinity
		}

		if score >= beta {
			return beta
		}

		if score > alpha {
			alpha = score
		}
	}

	e.tt.Store(hash, alpha, 0, TTExact, move.Move(0))
	return alpha
}
