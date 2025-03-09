package engine

import (
	"context"
	"sort"

	"github.com/Tecu23/argov2/internal/transposition"
	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/eval"
	"github.com/Tecu23/argov2/pkg/move"
)

// Constants used for controlling search depth, defining large numeric values, and scoring mating sequences
const (
	MaxDepth         = 64     // Maximum search depth
	Infinity         = 50_000 // A large value representing 'infinity' for alpha-beta cutoffs
	MateScore        = 49_000 // A bse score for checkmate scores
	MateDepth        = 48_000 // Used for adjusting mate scores closer to actual ply counts
	MaxKillers       = 2      // Maximum number of killer moves stored per depth
	NullMoveR        = 3      // Reduction for null move search
	NullMoveMinDepth = 3      // Minimum depth to try null move pruning
)

// MoveScore associates a move with a heuristic score to order moves
type MoveScore struct {
	move  move.Move
	score int
}

// orderMoves takes a list of moves and sorts them according to various heuristics:
//   - Transposition table move priority`
//   - Captures ordered by material gain (victim - ahhressor heuristic)
//   - Killer moves (non-captures previously leading to cutoffs)
//   - History heuristic (moves that were good in previous searches)
func (e *Engine) orderMoves(
	moves []move.Move,
	b *board.Board,
	ttMove move.Move,
	ply int,
) ([]move.Move, []int) {
	scores := make([]MoveScore, len(moves))
	stm := b.Side // side to move

	for i, mv := range moves {
		score := 0

		// If the move is the transposition table (TT) move, give it a large score
		if mv == ttMove {
			score = 20000
		} else if mv.GetCapture() != 0 {
			// Capture scoring: prioritize captures with a high "victim minus aggressor" value
			victim := b.GetPieceAt(mv.GetTarget())
			aggressor := b.GetPieceAt(mv.GetSource())
			score = 10000 + (eval.GetPieceValue(victim) - eval.GetPieceValue(aggressor)/10)
		} else {
			// For non-captures, check if it's a killer move, otherwise use history scores
			killerScore := e.killerTable.GetScore(mv, ply)
			if killerScore > 0 {
				score = killerScore
			} else {
				score = e.historyTable.Get(stm, mv.GetSource(), mv.GetTarget())
			}
		}
		scores[i] = MoveScore{mv, score}
	}

	// Sort mvoes by their scores in descending order (best first)
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	sortedMoves := make([]move.Move, len(moves))
	sortedScores := make([]int, len(moves))

	for i, ms := range scores {
		sortedMoves[i] = ms.move
		sortedScores[i] = ms.score
	}

	return sortedMoves, sortedScores
}

// nullMovePruning performs a null move search to determine if we can prune the current subtree.
// Returns true if the position is so good that we can likely prune this node.
func (e *Engine) nullMovePruning(
	ctx context.Context,
	b *board.Board,
	depth, beta, ply int,
	tm *timeManager,
) bool {
	// Skip null move pruning in this cases
	// 1. If in check (tactically volatile)
	// 2. If depth is too shallow
	// 3. If the endgame (material count is low - too dangerous)
	if b.InCheck() || depth < NullMoveMinDepth || e.evaluator.IsEndgame(b) {
		return false
	}

	// Make a null move (skip turn)
	copyB := b.CopyBoard()
	copyB.MakeNullMove()

	// Search with reduce depth
	// Use -beta+1 and -beta as bounds since we only need to know if score >= beta
	score := -e.alphaBeta(ctx, &copyB, depth-1-NullMoveR, -beta, -beta+1, ply+1, tm)

	return score >= beta
}

// search initiates an iterative deepening search, starting from depth = 1 up to maxDepth.
// It updates best move and score as depth increases and supports time management and aborting.
func (e *Engine) search(ctx context.Context, b *board.Board, tm *timeManager) SearchInfo {
	e.nodes = 0
	e.tt.NewSearch()       // Reset the transposition table for a new search
	e.historyTable.Clear() // Clear history heuristic table
	e.killerTable.Clear()  // Clear killer moves

	var bestMove move.Move
	var bestScore int

	maxDepth := MaxDepth
	if tm.limits.Depth > 0 {
		maxDepth = tm.limits.Depth
	}

	// Iterative deepening: incrementally increase search depth until max depth or time is up
	for depth := 1; depth <= maxDepth; depth++ {
		if tm.IsDone() {
			// Time's up or context canceled, stop searching
			break
		}

		// Perform a root-level alpha-beta search at the current depth
		score, mv := e.searchRoot(ctx, b, depth, tm)
		if mv != move.NoMove {
			bestMove = mv
			bestScore = score

			// Update principal variation information
			e.mainLine.moves = []move.Move{bestMove}
			e.mainLine.score = bestScore
			e.mainLine.depth = depth
			e.mainLine.nodes = e.nodes
		}

		// If there's a progress callback, update it with current search info
		if e.progress != nil {
			e.progress(e.createSearchInfo())
		}

		if tm.IsDone() {
			break
		}

		// Update time manager with node count changes for time allocation
		tm.OnNodesChanged(int(e.nodes))
	}

	// Create the final search info result
	searchInfo := e.createSearchInfo()
	if len(searchInfo.MainLine) == 0 && bestMove != move.NoMove {
		searchInfo.MainLine = []move.Move{bestMove}
	}

	return searchInfo
}

// searchRoot performs an alpha-beta search from the root position.
// It returns the best score and best move found at this depth
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

	moves := b.GenerateMoves() // Generate all pseudo-legal moves at root

	// Attempt to retrieve a transposition table (TT) entry to identify a good move to try first
	var ttMove move.Move
	if entry, ok := e.tt.Probe(b.Hash()); ok {
		ttMove = entry.BestMove
	}

	// Order moves using heuristics (TT move, captures, killers, history)
	moves, _ = e.orderMoves(moves, b, ttMove, 0)

	// Evaluate each move, searching one ply deeper with alpha-beta
	for _, mv := range moves {
		copyB := b.CopyBoard()
		// If the move is not legal, skip it
		if !copyB.MakeMove(mv, board.AllMoves) {
			continue
		}

		// Negamax formation: switch perspective by multiplying by -1
		score := -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, 1, tm)

		// Check if we ran out of time or context canceled
		if tm.IsDone() {
			return 0, move.NoMove
		}

		// Update alpha and best move if we found a better move
		if score > alpha {
			alpha = score
			bestMove = mv
			// Alpha-beta cutoff check
			if alpha >= beta {
				break
			}
		}
	}

	// Store the result in TT
	flag := transposition.TTExact
	if alpha <= originalAlpha {
		flag = transposition.TTAlpha
	} else if alpha >= beta {
		flag = transposition.TTBeta
	}
	e.tt.Store(b.Hash(), alpha, depth, flag, bestMove)

	return alpha, bestMove
}

// alphaBeta implements a negamax variant of the alpha-beta search algorithm.
// It uses transposition tables, move ordering, and other heuristics.
// 'depth' is how deep to search, 'alpha' and 'beta' are the search windows,
// 'ply' is the current depth from the root, and 'tm' manages time
func (e *Engine) alphaBeta(
	ctx context.Context,
	b *board.Board,
	depth, alpha, beta, ply int,
	tm *timeManager,
) int {
	// Periodically check for interruptions (time)
	if (e.nodes & 1023) == 0 {
		if tm.IsDone() {
			// If stopped, return a large negative or positive value depending on the ply parity
			if depth&1 == 0 {
				return -Infinity // for maximizing player
			}
			return Infinity // for minimizing player
		}
	}
	e.nodes++ // count the node

	originalAlpha := alpha
	hash := b.Hash()

	// e.logger.LogSearchNode(b, depth, ply, alpha, beta, 0, move.NoMove)

	// Probe the transposition table to see if we already know something about this position
	if entry, ok := e.tt.Probe(hash); ok {
		// e.logger.LogTT(hash, entry, depth)

		if entry.Depth >= depth {
			score := adjustScore(entry.Score, ply)
			// Depending on the entry's flag, we can return early or narrow our bound
			switch entry.Flag {
			case transposition.TTExact:
				return entry.Score
			case transposition.TTAlpha:
				if score <= alpha {
					return alpha
				}
			case transposition.TTBeta:
				if score >= beta {
					return beta
				}
			}
		}
	}

	// Try null move pruning
	if depth >= NullMoveMinDepth && !b.InCheck() {
		if e.nullMovePruning(ctx, b, depth, beta, ply, tm) {
			// e.logger.LogPruning("null move", depth, beta, beta)
			return beta
		}
	}

	// If we are in check, increase search depth to resolve the check
	if b.InCheck() {
		depth++
	}

	// Check terminal conditions
	if b.IsCheckmate() {
		// Return a large negative score plus the node count, preferring shorter mates
		return -MateScore + int(e.nodes) // Prefer shorter mates
	}

	if b.IsStalemate() || b.IsInsufficientMaterial() {
		// Draw or conditions return a score of 0
		return 0
	}

	// If we have reached zero depth, switch to quiescence search
	if depth <= 0 {
		return e.quiescence(ctx, b, alpha, beta, ply, tm)
	}
	moves := b.GenerateMoves()

	var ttMove move.Move
	if entry, ok := e.tt.Probe(hash); ok {
		ttMove = entry.BestMove
	}

	moves, _ = e.orderMoves(moves, b, ttMove, ply)
	// e.logger.LogMoveOrdering(ply, moves, scores)

	hasLegalMoves := false
	var bestMove move.Move
	bestScore := -Infinity
	moveCount := 0
	inCheck := b.InCheck()

	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.AllMoves) {
			continue // Illegal move
		}

		hasLegalMoves = true
		moveCount++

		isCapture := mv.GetCapture() != 0
		givesCheck := copyB.InCheck()

		// Late Move Reductions (LMR):
		// If move is not forcing (no check, no capture, not urgent), try a reduced depth search first
		reduction := 0
		if depth >= 3 && moveCount > 4 && !inCheck && !isCapture && !givesCheck {
			// reduction = e.reductionTable.Get(depth, moveCount)
			reduction = 1 // A simple LMR heuristics: reduce depth for "quiet" moves
		}

		var score int
		if reduction > 0 {
			// First do a reduced-depth search
			score = -e.alphaBeta(ctx, &copyB, depth-1-reduction, -(alpha + 1), -alpha, ply+1, tm)

			// If it improves alpha, do a full depth re-search
			if score > alpha {
				score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, ply+1, tm)
			}
		} else {
			// Regular full-depth search
			score = -e.alphaBeta(ctx, &copyB, depth-1, -beta, -alpha, ply+1, tm)
		}

		// Update best score and alpha
		if score > bestScore {
			bestScore = score
			bestMove = mv
			if score > alpha {
				// Update history heuristics for quiet moves
				if !isCapture && ply < MaxDepth {
					e.historyTable.Update(copyB.Side, mv.GetSource(), mv.GetTarget(), 1)
				}
				alpha = score
				if alpha >= beta {
					// Store killer move and update history on a beta cutoff
					if !isCapture && ply < MaxDepth {
						e.killerTable.Update(mv, ply)
						e.historyTable.Update(copyB.Side, mv.GetSource(), mv.GetTarget(), depth)
					}
					// Store this node as a fail-high node in TT
					e.tt.Store(hash, beta, depth, transposition.TTBeta, mv)
					return beta
				}
			}
		}
	}

	// If not legal moves, it's either checkmate or stalemate.
	if !hasLegalMoves {
		if b.InCheck() {
			return -MateScore + int(e.nodes)
		}
		return 0
	}

	// Store TT Entry
	flag := transposition.TTExact
	if bestScore <= originalAlpha {
		flag = transposition.TTAlpha
	}
	e.tt.Store(hash, bestScore, depth, flag, bestMove)

	return bestScore
}

// quiescence search is a specialized search used at leaf nodes to evaluate "quiet" positions.
// Instead of evaluating at leaves directly, it only considers capture moves and checks,
// trying to resolve unstable positions so that no immediate tactical gains are missed.
func (e *Engine) quiescence(
	ctx context.Context,
	b *board.Board,
	alpha, beta, ply int,
	tm *timeManager,
) int {
	e.nodes++

	// Peridically check for interrupts
	if (e.nodes & 4095) == 0 {
		if ctx.Err() != nil || tm.IsDone() {
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

	// Evaluate the static position
	score := e.evaluator.Evaluate(b)

	// Stand pat check: if current position is already good enough to cause beta cutoff, just return
	if score >= beta {
		return beta
	}

	// Raise alpha if this stand-pat evaluation improves it
	alpha = max(alpha, score)

	// Generate only captures for quiescence search
	moves := b.GenerateCaptures()
	moves, _ = e.orderMoves(moves, b, move.NoMove, ply) // Order captures

	// Explore all captures to see if any improve the evaluation
	for _, mv := range moves {
		copyB := b.CopyBoard()
		if !copyB.MakeMove(mv, board.OnlyCaptures) {
			continue // illegal move
		}

		// Negamax search on the resulting position after the capture
		score := -e.quiescence(ctx, &copyB, -beta, -alpha, ply+1, tm)

		if tm.IsDone() {
			// Interruption handling within quiescence
			if ply&1 == 0 {
				return -Infinity
			}
			return Infinity
		}

		if score >= beta {
			return beta // beta cutoff
		}

		// Update alpha if we found a better capture line
		alpha = max(alpha, score)
	}
	return alpha
}

// adjustScore adjusts mate scores to reflect the distance to mate.
// This ensures that closer mates have higher priority (or more negative for mates against us).
func adjustScore(score, ply int) int {
	if score >= MateScore-MaxDepth {
		// Reducing score by ply makes a mate closer more valuable (less negative offset)
		return score - ply
	}
	if score <= -MateScore+MaxDepth {
		// Increasing score by ply for negative mate scores makes closer mates worse for the side to move
		return score + ply
	}
	return score
}
