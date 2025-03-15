package engine

import (
	"context"
	"sort"

	"github.com/Tecu23/argov2/internal/transposition"
	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/evaluation"
	"github.com/Tecu23/argov2/pkg/move"
)

const (
	MaxDepth   = 64
	Infinity   = 50_000
	MateScore  = 49_000
	MateDepth  = 48_000
	MaxKillers = 2
)

type MoveScore struct {
	move  move.Move
	score int
}

func (e *Engine) updateKillers(mv move.Move, ply int) {
	if ply >= MaxDepth {
		return
	}
	// Don't store captures as killer moves
	if mv.GetCapture() != 0 {
		return
	}

	// Don't store a move that's already a killer at this ply
	for i := 0; i < MaxKillers; i++ {
		if e.killerMoves[ply][i] == mv {
			return
		}
	}

	// Shift existing killers and insert new one at first position
	for i := MaxKillers - 1; i > 0; i-- {
		e.killerMoves[ply][i] = e.killerMoves[ply][i-1]
	}
	e.killerMoves[ply][0] = mv
}

func (e *Engine) orderMoves(
	moves []move.Move,
	b *board.Board,
	ttMove move.Move,
	ply int,
) []move.Move {
	scores := make([]MoveScore, len(moves))
	stm := b.Side

	for i, mv := range moves {
		score := 0

		// TT move gets highest priority
		if mv == ttMove {
			score = 2_000_000
		} else if mv.GetCapture() != 0 {
			// MVV-LVA scoring
			victim := b.GetPieceAt(mv.GetTarget())
			aggressor := b.GetPieceAt(mv.GetSource())
			score = 1_000_000 + (evaluation.GetPieceValue(victim) - evaluation.GetPieceValue(aggressor)/10)
		} else {
			for j := 0; j < MaxKillers; j++ {
				if mv == e.killerMoves[ply][j] {
					score = 900_000 - j*1000
					break
				}
			}

			if score == 0 {
				score = e.historyTable.Get(stm, mv.GetSource(), mv.GetTarget())
			}
		}

		scores[i] = MoveScore{mv, score}
	}

	// Sort moves by score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Extract sorted moves
	sortedMoves := make([]move.Move, len(moves))
	for i, ms := range scores {
		sortedMoves[i] = ms.move
	}

	return sortedMoves
}

// search performs the actual search logic
func (e *Engine) search(ctx context.Context, b *board.Board, tm *timeManager) SearchInfo {
	e.nodes = 0
	e.tt.NewSearch()
	e.historyTable.Clear()
	e.killerMoves = [MaxDepth][MaxKillers]move.Move{}

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

	// Generate moves at root
	moves := b.GenerateMoves()

	hashKey := b.Hash()
	_, _, hashMove, _ := e.tt.Probe(hashKey, 0, 0, 0)

	moves = e.orderMoves(moves, b, hashMove, 0)

	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.AllMoves) {
			continue
		}

		// Search this position
		score := -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, 1, tm)

		// Check for search abort
		if ctx.Err() != nil || tm.IsDone() {
			return 0, move.NoMove
		}

		// Update best score if we found a better score
		if score > alpha {
			alpha = score
			bestMove = mv
			if alpha >= beta {
				break
			}
		}

	}
	// Store the result in the transposition table
	e.tt.Store(hashKey, depth, transposition.EXACT, alpha, bestMove)

	return alpha, bestMove
}

// alphaBeta performs the main alpha-beta search
func (e *Engine) alphaBeta(
	ctx context.Context,
	b *board.Board,
	depth, alpha, beta, ply int,
	tm *timeManager,
) int {
	// Check timeout periodically
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

	// Check extension
	if b.InCheck() {
		depth++
	}

	// Check for terminal positions
	if b.IsCheckmate() {
		return -MateScore + int(e.nodes) // Prefer shorter mates
	}

	if b.IsStalemate() || b.IsInsufficientMaterial() {
		return 0
	}

	// Get position key
	hashKey := b.Hash()

	// If this is not a PV node (zero window search), check the transposition table
	isPVNode := beta > alpha+1

	if !isPVNode {
		found, score, _, _ := e.tt.Probe(hashKey, depth, alpha, beta)
		if found {
			return score
		}
	}

	// Base case: evaluate leaf nodes
	if depth <= 0 {
		return e.quiescence(ctx, b, alpha, beta, ply, tm)
	}

	// Generate moves
	moves := b.GenerateMoves()

	// Try hash move first if available
	_, _, hashMove, _ := e.tt.Probe(hashKey, 0, 0, 0)

	moves = e.orderMoves(moves, b, hashMove, ply)

	hasLegalMoves := false

	bestScore := -Infinity
	bestMove := move.NoMove
	originalAlpha := alpha

	moveCount := 0
	inCheck := b.InCheck()

	// Search all moves {
	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.AllMoves) {
			continue
		}

		hasLegalMoves = true
		moveCount++

		var score int
		isCapture := mv.GetCapture() != 0
		givesCheck := copyB.InCheck()

		reduction := 0
		if depth >= 3 && moveCount > 4 && !inCheck && !isCapture && !givesCheck {
			// reduction = e.reductionTable.Get(depth, moveCount)
			reduction = 1
		}

		if reduction > 0 {
			score = -e.alphaBeta(ctx, &copyB, depth-1-reduction, -(alpha + 1), -alpha, ply+1, tm)

			if score > alpha {
				score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, ply+1, tm)
			}
		} else {
			score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, ply+1, tm)
		}

		if score > bestScore {
			bestScore = score
			bestMove = mv
			if score > alpha {
				if !isCapture && ply < MaxDepth {
					e.historyTable.Update(copyB.Side, mv.GetSource(), mv.GetTarget(), 1)
				}
				alpha = score
				if alpha >= beta {
					if !isCapture && ply < MaxDepth {
						e.updateKillers(mv, ply)
						e.historyTable.Update(copyB.Side, mv.GetSource(), mv.GetTarget(), depth)
					}
					e.tt.Store(hashKey, depth, transposition.BETA, beta, mv)
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

	var flag int
	if bestScore <= originalAlpha {
		flag = transposition.ALPHA
	} else if bestScore >= transposition.BETA {
		flag = transposition.BETA
	} else {
		flag = transposition.EXACT
	}

	e.tt.Store(hashKey, depth, flag, bestScore, bestMove)

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

	hashKey := b.Hash()
	found, score, _, _ := e.tt.Probe(hashKey, 0, alpha, beta)
	if found {
		return score
	}

	// Stand-pat score
	score = e.evaluator.Evaluate(b)
	if score >= beta {
		return beta
	}

	alpha = max(alpha, score)

	// Generate captures
	moves := b.GenerateCaptures()
	moves = e.orderMoves(moves, b, move.NoMove, ply) // Order captures

	bestScore := score
	bestMove := move.NoMove
	originalAlpha := alpha

	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.OnlyCaptures) {
			continue
		}

		score := -e.quiescence(ctx, &copyB, -beta, -alpha, ply+1, tm)

		if ctx.Err() != nil || tm.IsDone() {
			if ply&1 == 0 {
				return -Infinity
			}
			return Infinity
		}

		if score > bestScore {
			bestScore = score
			bestMove = mv

			alpha = max(alpha, score)

			if alpha >= beta {
				e.tt.Store(hashKey, 0, transposition.BETA, beta, mv)
				return beta
			}
		}
	}

	var flag int
	if bestScore <= originalAlpha {
		flag = transposition.ALPHA
	} else if bestScore >= beta {
		flag = transposition.BETA
	} else {
		flag = transposition.EXACT
	}

	e.tt.Store(hashKey, 0, flag, bestScore, bestMove)

	return bestScore
}
