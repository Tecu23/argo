package engine

import (
	"context"
	"sort"

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

// type SearchStats struct {
// 	killerHits    int
// 	historyHits   int
// 	ttHits        int
// 	betaCutoffs   int
// 	movesSearched int
// }
//
// // Print stats after search
// func (e *Engine) printStats() {
// 	total := float64(e.stats.movesSearched)
// 	fmt.Printf(
// 		"Beta cutoffs: %d (%.2f%%)\n",
// 		e.stats.betaCutoffs,
// 		float64(e.stats.betaCutoffs)/total*100,
// 	)
// 	fmt.Printf("TT hits: %d (%.2f%%)\n", e.stats.ttHits, float64(e.stats.ttHits)/total*100)
// 	fmt.Printf(
// 		"Killer hits: %d (%.2f%%)\n",
// 		e.stats.killerHits,
// 		float64(e.stats.killerHits)/total*100,
// 	)
// 	fmt.Printf(
// 		"History hits: %d (%.2f%%)\n",
// 		e.stats.historyHits,
// 		float64(e.stats.historyHits)/total*100,
// 	)
// }

func (e *Engine) orderMoves(
	moves []move.Move,
	b *board.Board,
	ttMove move.Move,
	ply int,
) []move.Move {
	// fmt.Printf("Ordering moves at ply %d\n", ply)

	scores := make([]MoveScore, len(moves))
	stm := b.Side

	for i, mv := range moves {
		score := 0
		// reason := "quiet move"

		// TT move gets highest priority
		if mv == ttMove {
			// reason = "TT move"
			score = 20000
		} else if mv.GetCapture() != 0 {
			// MVV-LVA scoring
			victim := b.GetPieceAt(mv.GetTarget())
			aggressor := b.GetPieceAt(mv.GetSource())
			score = 10000 + (evaluation.GetPieceValue(victim) - evaluation.GetPieceValue(aggressor)/10)
			// reason = fmt.Sprintf("capture (victim: %v, aggressor: %v)", victim, aggressor)

		} else {

			killerScore := e.killerTable.GetScore(mv, ply)
			if killerScore > 0 {
				score = killerScore
				// reason = fmt.Sprintf("killer (slot: %d)", killerScore/100)
			} else {
				score = e.historyTable.Get(stm, mv.GetSource(), mv.GetTarget())
				// reason = fmt.Sprintf("history (score: %d)", score)
			}
		}
		scores[i] = MoveScore{mv, score}
		// fmt.Printf("  Move: %s Score: %d Reason: %s\n", mv.String(), score, reason)

	}

	// Sort moves by score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	// Extract sorted moves
	sortedMoves := make([]move.Move, len(moves))

	// fmt.Println("After sorting:")
	// for _, ms := range scores {
	// 	fmt.Printf("  Move: %s Score: %d\n", ms.move.String(), ms.score)
	// }

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
	e.killerTable.Clear()

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
	originalAlpha := alpha

	// Generate moves at root
	moves := b.GenerateMoves()

	var ttMove move.Move
	if entry, ok := e.tt.Probe(b.Hash()); ok {
		ttMove = entry.BestMove
	}

	moves = e.orderMoves(moves, b, ttMove, 0)

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
	// Store in TT
	flag := TTExact
	if alpha <= originalAlpha {
		flag = TTAlpha
	} else if alpha >= beta {
		flag = TTBeta
	}
	e.tt.Store(b.Hash(), alpha, depth, flag, bestMove)

	return alpha, bestMove
}

// alphaBeta performs the main alpha-beta search
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

	// TT Lookup
	hash := b.Hash()
	if entry, ok := e.tt.Probe(hash); ok {
		if entry.Depth >= depth {
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

	if b.IsStalemate() || b.IsInsufficientMaterial() {
		return 0
	}

	// Base case: evaluate leaf nodes
	if depth <= 0 {
		return e.quiescence(ctx, b, alpha, beta, ply, tm)
	}

	// Generate moves
	moves := b.GenerateMoves()
	var ttMove move.Move

	if entry, ok := e.tt.Probe(hash); ok {
		ttMove = entry.BestMove
	}

	moves = e.orderMoves(moves, b, ttMove, ply)

	hasLegalMoves := false
	var bestMove move.Move
	bestScore := -Infinity
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

		// if score >= beta {
		// 	if mv == ttMove {
		// 		e.stats.ttHits++
		// 	} else if e.killerTable.IsKiller(mv, ply) {
		// 		e.stats.killerHits++
		// 	} else if score := e.historyTable.Get(copyB.Side, mv.GetSource(), mv.GetTarget()); score > 0 {
		// 		e.stats.historyHits++
		// 	}
		// 	e.stats.betaCutoffs++
		// }
		//
		// e.stats.movesSearched++

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
						e.killerTable.Update(mv, ply)
						e.historyTable.Update(copyB.Side, mv.GetSource(), mv.GetTarget(), depth)
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

		score := -e.quiescence(ctx, &copyB, -beta, -alpha, ply+1, tm)

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

func adjustScore(score, ply int) int {
	if score >= MateScore-MaxDepth {
		return score - ply
	}
	if score <= -MateScore+MaxDepth {
		return score + ply
	}
	return score
}
