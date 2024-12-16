package engine

import (
	"context"
	"fmt"
	"math"
	"time"

	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/move"
)

// Constants defining difficulty and branch factors for time calculation
const (
	maxDifficulty   = 2    // The maximum difficulty scaling factor
	minBranchFactor = 0.75 // Minimum branch factor used in time calculation
	maxBranchFactor = 1.5  // Maximum branch factor used in time calculation
)

// timeManager is responsible for determining and enforcing time limits during the search.
// It decides how long the engine can spend on the current move based on the time controls,
// difficulty adjustments, and ongoing search feedback (like changing best moves or scores)
type timeManager struct {
	start        time.Time          // The moment the search started
	limits       LimitsType         // UCI-style time control limits (depth, nodes, movetime, etc..)
	side         bool               // Side to move: true = White, false = Black
	difficulty   float64            // A scaling factor  influencing time usage (adjusted based on search results)
	lastScore    int                // The best score from the previous iteration
	lastBestMove move.Move          // the best move found in the previous iteration
	done         <-chan struct{}    // A channel that signals when the allowed time or conditions are met (cancellation)
	cancel       context.CancelFunc // a function to cancel the ongoing context, stopping the search
}

// newTimeManager creates and initializes a timeManager instance. It sets up a context with appropriate
// deadlines or cancellations based on the provided LimisType and the current game situation.
func newTimeManager(
	ctx context.Context,
	start time.Time,
	limits LimitsType,
	b *board.Board,
) *timeManager {
	tm := &timeManager{
		start:      start,
		limits:     limits,
		side:       b.Side == color.WHITE,
		difficulty: 1, // Start with a neutral difficulty factor
	}

	var cancel context.CancelFunc

	// If a MoveTime or classical clock times (WhiteTime/BlackTime) are set,
	// determine the maximum allowed time for this move.
	// Otherwise, we just allow the infinite search untill stopped manually or by conditions.
	if limits.MoveTime > 0 || limits.WhiteTime > 0 || limits.BlackTime > 0 {
		var maximum time.Duration
		if limits.MoveTime > 0 {
			// If MoveTime is specified, use it directly as the max time.
			maximum = time.Duration(limits.MoveTime) * time.Millisecond
		} else {
			// Otherwise, calculate a time limit based on difficulty and maximum branch factor.
			maximum = tm.calculateTimeLimit(maxDifficulty, maxBranchFactor)
		}

		// Create a context that will expire once we reach the computed maximum time.
		ctx, cancel = context.WithDeadline(ctx, start.Add(maximum))
	} else {
		// No explicit time constratints, use a cancelable context that can be ended on conditions.
		ctx, cancel = context.WithCancel(ctx)
	}

	// Store the done channel and the cancel; function to use later
	tm.done = ctx.Done()
	tm.cancel = cancel
	return tm
}

// IsDone check if the time manager's context is already signaled as done (i.e., time is up or canceled)
func (tm *timeManager) IsDone() bool {
	select {
	case <-tm.done:
		return true
	default:
		return false
	}
}

// OnNodesChanged is called when the search has processed a certain number of nodes.
// If the limits specify a node limit and we surpass it, we stop the search
func (tm *timeManager) OnNodesChanged(nodes int) {
	if tm.limits.Nodes > 0 && nodes >= tm.limits.Nodes {
		tm.cancel() // Stop the search id we've hit the node limit
	}
}

// On IterationComplete is called after an iterative deepening search iteration completes.
// It receives the current mainLine (best line of moves found) and decides whether to continue searching
// or stop based on factors like depth reached, score changes, and allocated time.
func (tm *timeManager) OnIterationComplete(line mainLine) {
	// If running in "infinite" mode (like analysis mode), never cancel due to time/depth.
	if tm.limits.Infinite {
		return
	}

	// If a depth limit is set and we reached it, stop searching
	if tm.limits.Depth != 0 && line.depth >= tm.limits.Depth {
		tm.cancel()
		return
	}

	// If a depth limit is set and we reached it, stop searching
	if line.score >= winIn(line.depth-5) || line.score <= lossIn(line.depth-5) {
		tm.cancel()
		return
	}

	// If we have a winning or losing position near mate (score close to ±MateScore) at some depth,
	// we can stop searching early since the outcome is already known.
	if tm.limits.WhiteTime > 0 || tm.limits.BlackTime > 0 {
		// Once we have some reasonable depth, start adjusting difficulty based on changes in score or best move.
		if line.depth >= 5 {
			scoreDrop := tm.lastScore - line.score
			if scoreDrop > 50 {
				// If the score dropped significantly, increase difficulty (spend more time)
				tm.difficulty = math.Min(tm.difficulty*1.3, maxDifficulty)
			} else if line.moves[0] != tm.lastBestMove {
				// If the best move changed from the last iteration, slightly increase difficulty
				tm.difficulty = math.Min(tm.difficulty*1.2, 1.5)
			} else {
				// If the best move stayed the same, slightly reduce difficulty to save time
				tm.difficulty = math.Max(0.8, tm.difficulty*0.8)
			}
		}

		// Update last known best score and best move
		tm.lastScore = line.score
		tm.lastBestMove = line.moves[0]

		// Calculate an "optimum" time limit based on the current difficulty and a minimal branch factor
		optimum := tm.calculateTimeLimit(tm.difficulty, minBranchFactor)

		// If we have already spent more than 'optimum' time, stop the search
		if time.Since(tm.start) >= optimum {
			tm.cancel()
			return
		}
	}
}

// Close stops the time manager and thus cancels the search
func (tm *timeManager) Close() {
	tm.cancel()
}

// calculateTimeLimit determines how much time to allocate for a move given the current difficulty,
// branch factor, and the known time control parameters (WhiteTime, BlackTime, increments, etc..)
func (tm *timeManager) calculateTimeLimit(difficulty, branchFactor float64) time.Duration {
	const (
		DefaultMovesToGo = 40                     // Assume 40 moves remain if not specified
		MoveOverhead     = 300 * time.Millisecond // Deduct a small overhead from the total time
		MinTimeLimit     = 1 * time.Millisecond   // Minimnum allocated time (avoid zero or negatime)
		EmergencyTime    = 10 * time.Second
		CriticalTime     = 5 * time.Second
	)

	// Determine main time and increment based on which side is moving
	var main, inc time.Duration
	if tm.side {
		main = time.Duration(tm.limits.WhiteTime) * time.Millisecond
		inc = time.Duration(tm.limits.WhiteIncrement) * time.Millisecond
	} else {
		main = time.Duration(tm.limits.BlackTime) * time.Millisecond
		inc = time.Duration(tm.limits.BlackIncrement) * time.Millisecond
	}

	// Deduct the overhead
	main -= MoveOverhead
	if main < MinTimeLimit {
		main = MinTimeLimit
	}

	// Determine how many moves remain until the nest time control (or assume default if none given)
	moves := tm.limits.MovesToGo
	if moves == 0 || moves > DefaultMovesToGo {
		moves = DefaultMovesToGo
	}

	// Total think tume = current main time + increments for the remaining moves
	total := float64(main) + float64(moves-1)*float64(inc)

	// Calculate time allocation using a formula that considers difficult and branch factor.
	// The time given is spread out over the remaining moves, adjusted by difficulty and the branch factor.
	timeLimit := time.Duration(
		difficulty * branchFactor * total / (difficulty*maxBranchFactor + float64(moves-1)),
	)

	// Add time scaling based on remaining time
	timeScaleFactor := 1.0
	if main < 30*time.Second {
		timeScaleFactor = 0.3 // Use only 30% if calculated time when low
	} else if main < 60*time.Second {
		timeScaleFactor = 0.5 // Use 50% when under a minute
	}

	timeLimit = time.Duration(float64(timeLimit) * timeScaleFactor)

	// Add progressive moves estimation
	if moves > 30 {
		// Early game - be more conservative
		timeLimit = timeLimit * 2 / 3
	}

	// Ensure we don't exceed the main time or go below the minimum time limit
	if timeLimit > main {
		timeLimit = main
	}
	if timeLimit < MinTimeLimit {
		timeLimit = MinTimeLimit
	}
	return timeLimit
}

func (tm *timeManager) String() string {
	sideStr := "Black"
	if tm.side {
		sideStr = "White"
	}

	doneSet := (tm.done != nil)
	cancelSet := (tm.cancel != nil)

	return fmt.Sprintf(`timeManager:
  start:        %v
  limits:
    Ponder:         %t
    Infinite:       %t
    WhiteTime:      %d
    BlackTime:      %d
    WhiteIncrement: %d
    BlackIncrement: %d
    MoveTime:       %d
    MovesToGo:      %d
    Depth:          %d
    Nodes:          %d
    Mate:           %d
  side:         %s
  difficulty:   %.2f
  lastScore:    %d
  lastBestMove: %v
  done:         %t
  cancel:       %t
`,
		tm.start,
		tm.limits.Ponder,
		tm.limits.Infinite,
		tm.limits.WhiteTime,
		tm.limits.BlackTime,
		tm.limits.WhiteIncrement,
		tm.limits.BlackIncrement,
		tm.limits.MoveTime,
		tm.limits.MovesToGo,
		tm.limits.Depth,
		tm.limits.Nodes,
		tm.limits.Mate,
		sideStr,
		tm.difficulty,
		tm.lastScore,
		tm.lastBestMove,
		doneSet,
		cancelSet,
	)
}
