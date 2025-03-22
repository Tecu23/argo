// Package engine keeps the running the search and engine logic
package engine

import (
	"context"

	"github.com/Tecu23/argov2/internal/reduction"
	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/move"
)

// search performs the actual search logic
func (e *Engine) search(ctx context.Context, b *board.Board, tm *timeManager) SearchInfo {
	e.nodes = 0
	e.tt.NewSearch()
	e.historyTable.Clear()
	e.killerMoves = [MaxDepth][MaxKillers]move.Move{}

	e.evaluator.Reset(b)

	var bestMove move.Move
	var bestScore int

	// Engine will only loop up to depth limit, regardless of time
	maxDepth := MaxDepth
	if tm.limits.Depth > 0 {
		maxDepth = tm.limits.Depth
	}

	// Iterative deepeing
	for depth := 1; depth <= maxDepth; depth++ {
		if tm.IsDone() || ctx.Err() != nil {
			break
		}

		// Report current search depth
		if e.progress != nil {
			info := e.createSearchInfo()
			info.Depth = depth
			e.progress(info)
		}

		score, mv := e.searchRoot(ctx, b, depth, tm)
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

		// If we found a forced mate, no need to search deeper
		if bestScore > MateScore-MaxDepth || bestScore < -MateScore+MaxDepth {
			break
		}
	}

	searchInfo := e.createSearchInfo()
	// Ensure we have a move to return
	if len(searchInfo.MainLine) == 0 && bestMove != move.NoMove {
		searchInfo.MainLine = []move.Move{bestMove}
	}

	return searchInfo
}

// searchRoot performs alpha-beta search at the root level
func (e *Engine) searchRoot(
	ctx context.Context,
	b *board.Board,
	depth int,
	tm *timeManager,
) (int, move.Move) {
	alpha := -Infinity
	beta := Infinity
	var bestMove move.Move
	originalAlpha := alpha

	// Generate moves at root
	moves := b.GenerateMoves()

	// Check for single legal move - if only 1 move is available, return it immediately
	if len(moves) == 1 {
		cpy := b.CopyBoard()
		if cpy.MakeMove(moves[0], board.AllMoves) {
			e.evaluator.ProcessMove(&cpy, moves[0])
			score := e.evaluator.Evaluate(&cpy)
			e.evaluator.PopAccumulation()
			return score, moves[0]
		}
	}

	var ttMove move.Move
	if entry, ok := e.tt.Probe(b.Hash()); ok {
		ttMove = entry.BestMove
	}

	moves = e.orderMoves(moves, b, ttMove, 0)

	bestScore := -Infinity
	moveCount := 0

	for i, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.AllMoves) {
			continue
		}

		e.evaluator.ProcessMove(&copyB, mv)
		moveCount++

		var score int

		// For the first move or promising moves, do a full-window search
		if i == 0 {
			score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, 1, tm)
		} else {
			// Use zero-window search for other moves
			score = -e.alphaBeta(ctx, &copyB, depth-1, -alpha-1, -alpha, 1, tm)

			// If the score exceeds alpha but is below beta, re-search with full window
			if score > alpha && score < beta {
				score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, 1, tm)
			}
		}

		e.evaluator.PopAccumulation()

		// Check for search abort
		if ctx.Err() != nil || tm.IsDone() {
			return 0, move.NoMove
		}

		// Update best score if we found a better score
		if score > bestScore {
			bestScore = score
			bestMove = mv

			if score > alpha {
				alpha = score

				// For the main variation, we store the move in the PV
				if e.mainLine.moves == nil {
					e.mainLine.moves = make([]move.Move, 1)
					e.mainLine.moves[0] = mv
				} else if len(e.mainLine.moves) > 0 {
					e.mainLine.moves[0] = mv
				}

				if alpha >= beta {
					break
				}
			}
		}
	}

	// If no legal moves were found
	if moveCount == 0 {
		if b.InCheck() {
			return -MateScore, move.NoMove
		}

		return 0, move.NoMove
	}

	// Store in TT
	flag := TTExact
	if bestScore <= originalAlpha {
		flag = TTAlpha
	} else if bestScore >= beta {
		flag = TTBeta
	}
	e.tt.Store(b.Hash(), alpha, depth, flag, bestMove)

	return bestScore, bestMove
}

// alphaBeta performs the main alpha-beta search with Principal Variation Search
func (e *Engine) alphaBeta(
	ctx context.Context,
	b *board.Board,
	depth, alpha, beta, ply int,
	tm *timeManager,
) int {
	if (e.nodes & 1023) == 0 {
		if ctx.Err() != nil || tm.IsDone() {
			// This ensures the move won't be selected
			// as it will always be worse than any real evaluation
			if depth&1 == 0 {
				return -Infinity // for maximizing player
			}
			return Infinity // for minimizing player
		}
	}

	// Increment node counter
	e.nodes++

	originalAlpha := alpha
	isPV := beta > alpha+1 // Check if this is a PV node

	// TT Lookup
	hash := b.Hash()
	var ttMove move.Move
	if entry, ok := e.tt.Probe(hash); ok {
		ttMove = entry.BestMove

		// We can use TT cutoffs in non-PV nodes when depth is sufficient
		if !isPV && entry.Depth >= depth {
			score := adjustScore(entry.Score, ply)
			switch entry.Flag {
			case TTExact:
				return entry.Score
			case TTAlpha:
				if score <= alpha {
					return alpha
				}
			case TTBeta:
				if score >= beta {
					return beta
				}
			}
		}
	}

	// Check extension
	if b.InCheck() {
		depth++
	}

	// Check for terminal positions
	if b.IsCheckmate() {
		return -MateScore + int(e.nodes) // Prefer shorter mates
	}

	// Base case: evaluate leaf nodes
	if depth <= 0 {
		return e.quiescence(ctx, b, alpha, beta, ply, tm)
	}

	// Generate moves
	moves := b.GenerateMoves()
	moves = e.orderMoves(moves, b, ttMove, ply)

	hasLegalMoves := false
	var bestMove move.Move
	bestScore := -Infinity
	moveCount := 0
	inCheck := b.InCheck()

	// Search all moves {
	for i, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.AllMoves) {
			continue
		}

		e.evaluator.ProcessMove(&copyB, mv)

		hasLegalMoves = true
		moveCount++

		var score int
		isCapture := mv.IsCapture()
		givesCheck := copyB.InCheck()

		reduct := 0
		if depth >= reduction.MinDepthForReduction &&
			moveCount > reduction.MinMovesBeforeReduction &&
			!inCheck && !isCapture &&
			mv.GetPromotedPiece() == 0 && !givesCheck {

			// Get history score for this move
			historyScore := e.historyTable.Get(b.Side, mv.GetSourceSquare(), mv.GetTargetSquare())

			// Calculate reduction with adjustments
			reduct = e.reductionTable.GetWithAdjustments(depth, moveCount, isPV, historyScore)

			// Additional dynamic adjustments

			// 1. Reduce less for killer moves
			if mv == e.killerMoves[ply][0] || mv == e.killerMoves[ply][1] {
				reduct = max(0, reduct-1)
			}

			// 2. Ensure we don't reduce into the quiescence search
			if depth-reduct <= 0 {
				reduct = max(0, depth-1)
			}
		}

		// PVS logic
		if i == 0 {
			// Full window search for first move
			score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, ply+1, tm)
		} else {
			// Try with zero window for non-first moves
			if reduct > 0 {
				// Reduced depth zero window search
				score = -e.alphaBeta(ctx, &copyB, depth-1-reduct, -alpha-1, -alpha, ply+1, tm)
			} else {
				// Normal depth zero window search
				score = -e.alphaBeta(ctx, &copyB, depth-1, -alpha-1, -alpha, ply+1, tm)
			}

			if score > alpha && reduct > 0 {
				score = -e.alphaBeta(ctx, &copyB, depth-1, -alpha-1, -alpha, ply+1, tm)
			}

			// If still promising, do a full-window search
			if score > alpha && score < beta {
				score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, ply+1, tm)
			}
		}

		e.evaluator.PopAccumulation()

		if score > bestScore {
			bestScore = score
			bestMove = mv
			if score > alpha {
				if !isCapture && ply < MaxDepth {
					e.historyTable.Update(copyB.Side, mv.GetSourceSquare(), mv.GetTargetSquare(), 1)
				}
				alpha = score
				if alpha >= beta {
					if !isCapture && ply < MaxDepth {
						e.updateKillers(mv, ply)
						e.historyTable.Update(
							copyB.Side,
							mv.GetSourceSquare(),
							mv.GetTargetSquare(),
							depth,
						)
					}
					e.tt.Store(hash, beta, depth, TTBeta, mv)
					return beta
				}
			}
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
	flag := TTExact
	if bestScore <= originalAlpha {
		flag = TTAlpha
	}
	e.tt.Store(hash, bestScore, depth, flag, bestMove)

	return bestScore
}

// quiescence performs capture-only search to reach quiet position
func (e *Engine) quiescence(
	ctx context.Context,
	b *board.Board,
	alpha, beta, ply int,
	tm *timeManager,
) int {
	e.nodes++

	if (e.nodes & 4095) == 0 {
		if ctx.Err() != nil || tm.IsDone() {
			// Use same logic as alpha-beta for consistency
			if e.nodes&1 == 0 {
				return -Infinity
			}
			return Infinity
		}
	}

	// Check for terminal positions
	if b.IsCheckmate() {
		return -MateScore + int(e.nodes) // Prefer shorter mates
	}

	if b.IsStalemate() || b.IsInsufficientMaterial() {
		return 0
	}

	// Stand-pat score
	score := e.evaluator.Evaluate(b)
	if score >= beta {
		return beta
	}

	alpha = max(alpha, score)

	// Generate captures
	moves := b.GenerateCaptures()
	moves = e.orderMoves(moves, b, move.NoMove, ply) // Order captures

	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.OnlyCaptures) {
			continue
		}

		e.evaluator.ProcessMove(&copyB, mv)

		score := -e.quiescence(ctx, &copyB, -beta, -alpha, ply+1, tm)

		e.evaluator.PopAccumulation()

		if ctx.Err() != nil || tm.IsDone() {
			if ply&1 == 0 {
				return -Infinity
			}
			return Infinity
		}

		if score >= beta {
			return beta
		}
		alpha = max(alpha, score)

	}

	return alpha
}
