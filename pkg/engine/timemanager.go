package engine

import (
	"context"
	"time"

	. "github.com/Tecu23/argov2/internal/types"
	"github.com/Tecu23/argov2/pkg/move"
)

// timeManager handle all time-related decisions during search
type timeManager struct {
	start             time.Time          // When search started
	allocated         time.Duration      // Time allocated for this move
	maxTime           time.Duration      // Maximum time allowed
	limits            LimitsType         // UCI time control limits
	side              bool               // Side to move (false=white, true=black)
	lastScore         int                // Previous iteration score
	scoreChanged      bool               // If score changed significantly from the last iteration
	lastBestMove      move.Move          // Previous iteration best move
	searchStable      bool               // If best move hasn't changed in recent iterations
	movesUntilControl int                // Estimated moves until next time control
	done              <-chan struct{}    // Channel to signal completion
	cancel            context.CancelFunc // Function to cancel search
}

func newTimeManager(
	limits LimitsType,
	side bool,
	done <-chan struct{},
	cancel context.CancelFunc,
) *timeManager {
	tm := &timeManager{
		start:  time.Now(),
		limits: limits,
		side:   side,
		done:   done,
		cancel: cancel,
	}

	tm.calculateTimeAllocation()
	return tm
}

func (tm *timeManager) calculateTimeAllocation() {
	// Handle fixed move time
	if tm.limits.MoveTime > 0 {
		tm.allocated = time.Duration(tm.limits.MoveTime) * time.Millisecond
		tm.maxTime = tm.allocated
		return
	}

	// Handle infinite/ponder
	if tm.limits.Infinite || tm.limits.Ponder {
		tm.allocated = 24 * time.Hour // Effectively infinite
		tm.maxTime = tm.allocated
		return
	}

	// Get relevant time and increment
	var baseTime, increment int
	if !tm.side { // White
		baseTime = tm.limits.WhiteTime
		increment = tm.limits.WhiteIncrement
	} else {
		baseTime = tm.limits.WhiteTime
		increment = tm.limits.WhiteIncrement
	}

	// Calculate moves until time control
	if tm.limits.MovesToGo > 0 {
		tm.movesUntilControl = tm.limits.MovesToGo
	} else {
		tm.movesUntilControl = 30 // Default estimate
	}

	// Calculate base allocation
	remainingTime := time.Duration(baseTime) * time.Millisecond
	incTime := time.Duration(increment) * time.Millisecond

	// Basic time management: allocate remaining time / moves plus some increment
	tm.allocated = (remainingTime / time.Duration(tm.movesUntilControl)) + (incTime / 2)

	// Set maximum time to prevent going over
	tm.maxTime = minTime(remainingTime/4, tm.allocated*2)
}

func (tm *timeManager) shouldStop() bool {
	// Checkm for external stop signal
	select {
	case <-tm.done:
		return true
	default:
	}

	// Don't stop if infinite or pondering
	if tm.limits.Infinite || tm.limits.Ponder {
		return false
	}

	elapsed := time.Since(tm.start)

	// Always stop if exceeded maximum time
	if elapsed >= tm.maxTime {
		return true
	}

	// Consider stopping at allocated time if search is stable
	if elapsed >= tm.allocated && tm.searchStable {
		return true
	}

	// Consider stopping at allocated time if score has not changed
	if elapsed >= tm.allocated && !tm.scoreChanged {
		return true
	}

	return false
}

func (tm *timeManager) updateSearch(score int, bestMove move.Move) {
	// Check if score changed significantly (more than 0.5 pawns)
	tm.scoreChanged = abs(score-tm.lastScore) > 50

	// Check if best move is stable
	tm.searchStable = bestMove == tm.lastBestMove

	tm.lastScore = score
	tm.lastBestMove = bestMove
}

// Helper functions
func minTime(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
