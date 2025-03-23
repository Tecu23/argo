# Time Management System

## Overview

Effective time management is a critical aspect of modern chess engines that
directly impacts playing strength. ArGO implements a sophisticated time
management system that balances the need for deep analysis with the practical
constraints of chess time controls. This document details the implementation
and decision-making process of ArGO's time manager.

## Time Manager Architecture

The time management system is encapsulated in the `timeManager` struct:

```go
type timeManager struct {
    start        time.Time          // When search started
    limits       LimitsType         // Time control parameters
    side         bool               // Side to move (true = White, false = Black)
    difficulty   float64            // Position complexity scaling factor
    lastScore    int                // Score from previous iteration
    lastBestMove move.Move          // Best move from previous iteration
    done         <-chan struct{}    // Signal channel for termination
    cancel       context.CancelFunc // Function to cancel context
}
```

## Initialization

The time manager is initialized at the beginning of each search:

```go
func newTimeManager(
    ctx context.Context,
    start time.Time,
    limits LimitsType,
    b *board.Board,
) *timeManager {
    tm := &timeManager{
        start:      start,
        limits:     limits,
        side:       b.SideToMove == color.WHITE,
        difficulty: 1, // Start with neutral difficulty
    }

    var cancel context.CancelFunc

    // Set up deadline based on time controls
    if limits.MoveTime > 0 || limits.WhiteTime > 0 || limits.BlackTime > 0 {
        var maximum time.Duration
        if limits.MoveTime > 0 {
            maximum = time.Duration(limits.MoveTime) * time.Millisecond
        } else {
            maximum = tm.calculateTimeLimit(maxDifficulty, maxBranchFactor)
        }

        ctx, cancel = context.WithDeadline(ctx, start.Add(maximum))
    } else {
        ctx, cancel = context.WithCancel(ctx)
    }

    tm.done = ctx.Done()
    tm.cancel = cancel
    return tm
}
```

## Time Calculation

The core of the time management system is the time allocation algorithm:

```go
func (tm *timeManager) calculateTimeLimit(difficulty, branchFactor float64) time.Duration {
    const (
        DefaultMovesToGo = 40                     // Assume 40 moves remain if not specified
        MoveOverhead     = 300 * time.Millisecond // Time buffer for communication
        MinTimeLimit     = 1 * time.Millisecond   // Minimum allocation
    )

    // Determine main time and increment
    var main, inc time.Duration
    if tm.side {
        main = time.Duration(tm.limits.WhiteTime) * time.Millisecond
        inc = time.Duration(tm.limits.WhiteIncrement) * time.Millisecond
    } else {
        main = time.Duration(tm.limits.BlackTime) * time.Millisecond
        inc = time.Duration(tm.limits.BlackIncrement) * time.Millisecond
    }

    // Deduct overhead
    main -= MoveOverhead
    if main < MinTimeLimit {
        main = MinTimeLimit
    }

    // Determine moves remaining
    moves := tm.limits.MovesToGo
    if moves == 0 || moves > DefaultMovesToGo {
        moves = DefaultMovesToGo
    }

    // Calculate total available thinking time
    total := float64(main) + float64(moves-1)*float64(inc)

    // Primary allocation formula
    timeLimit := time.Duration(
        difficulty * branchFactor * total / (difficulty*maxBranchFactor + float64(moves-1)),
    )

    // Scale based on remaining time
    timeScaleFactor := 1.0
    if main < 30*time.Second {
        timeScaleFactor = 0.3 // Critical time situation
    } else if main < 60*time.Second {
        timeScaleFactor = 0.5 // Low time situation
    }
    timeLimit = time.Duration(float64(timeLimit) * timeScaleFactor)

    // Adjust for game phase
    if moves > 30 {
        // Early game - be more conservative
        timeLimit = timeLimit * 2 / 3
    }

    // Apply constraints
    if timeLimit > main {
        timeLimit = main
    }
    if timeLimit < MinTimeLimit {
        timeLimit = MinTimeLimit
    }

    return timeLimit
}
```

## Adaptive Difficulty

A key innovation in ArGO's time management is the concept of "difficulty" -
a dynamic scaling factor that adjusts based on position characteristics:

```go
// Called after each iteration completes
func (tm *timeManager) OnIterationComplete(line mainLine) {
    // Skip time management in infinite analysis mode
    if tm.limits.Infinite {
        return
    }

    // Once we have some reasonable depth, adjust difficulty
    if line.depth >= 5 {
        scoreDrop := tm.lastScore - line.score
        if scoreDrop > 50 {
            // Score dropped significantly - increase difficulty
            tm.difficulty = math.Min(tm.difficulty*1.3, maxDifficulty)
        } else if line.moves[0] != tm.lastBestMove {
            // Best move changed - increase difficulty slightly
            tm.difficulty = math.Min(tm.difficulty*1.2, 1.5)
        } else {
            // Best move stable - reduce difficulty
            tm.difficulty = math.Max(0.8, tm.difficulty*0.8)
        }
    }

    // Update last known best score and move
    tm.lastScore = line.score
    tm.lastBestMove = line.moves[0]

    // Check if we've spent enough time
    optimum := tm.calculateTimeLimit(tm.difficulty, minBranchFactor)
    if time.Since(tm.start) >= optimum {
        tm.cancel()
        return
    }
}
```

## Branch Factor

The branch factor represents how many nodes the engine needs to search
to increase its depth by one ply. ArGO uses two branch factor values:

```go
const (
    minBranchFactor = 0.75 // Conservative estimate
    maxBranchFactor = 1.5  // Aggressive estimate
)
```

- `maxBranchFactor` is used to calculate the maximum allowed time
- `minBranchFactor` is used to calculate the optimal time (when the engine
  should consider stopping)

This dual approach provides a balance between allocating enough time for
complex positions while not overspending on simpler positions.

## Search Monitoring

During search, the time manager monitors several conditions that might
trigger early termination:

### 1. Fixed Limits Check

```go
func (tm *timeManager) OnIterationComplete(line mainLine) {
    // Check fixed constraints like depth
    if tm.limits.Depth != 0 && line.depth >= tm.limits.Depth {
        tm.cancel()
        return
    }

    // Check for found mate
    if line.score >= winIn(line.depth-5) || line.score <= lossIn(line.depth-5) {
        tm.cancel()
        return
    }

    // Other time-based checks
    // ...
}
```

### 2. Nodes Limit

```go
func (tm *timeManager) OnNodesChanged(nodes int) {
    if tm.limits.Nodes > 0 && nodes >= tm.limits.Nodes {
        tm.cancel() // Stop if node limit reached
    }
}
```

### 3. Cancellation Check

The engine periodically checks if it should stop searching:

```go
// In alphaBeta function
if (e.nodes & 1023) == 0 {
    if ctx.Err() != nil || tm.IsDone() {
        // Return appropriate score to stop search
    }
}
```

## Emergency Time Management

ArGO implements special handling for critical time situations:

```go
const (
    EmergencyTime    = 10 * time.Second
    CriticalTime     = 5 * time.Second
)

// When time is critically low, use simplified allocation
if main < EmergencyTime {
    // Use fixed percentage of remaining time
    return main / 5 // Use 20% of remaining time
}

if main < CriticalTime {
    // Use minimum safe allocation
    return MinTimeLimit * 10 // Minimal search
}
```

## Special Time Controls

### 1. Fixed Move Time

When a specific move time is specified:

```go
if limits.MoveTime > 0 {
    // Simply use the specified time
    maximum = time.Duration(limits.MoveTime) * time.Millisecond
}
```

### 2. Infinite Analysis

For infinite analysis mode:

```go
if tm.limits.Infinite {
    return // Never stop based on time
}
```

### 3. Tournament Controls

For standard tournament time controls:

```go
// For classical time control with increment
if tm.limits.WhiteTime > 0 || tm.limits.BlackTime > 0 {
    // Detailed time management using the formulas described above
}
```

## Time Usage Strategy by Game Phase

ArGO adjusts its time usage based on the game phase:

```go
// Add progressive moves estimation
if moves > 30 {
    // Early game - be more conservative
    timeLimit = timeLimit * 2 / 3
} else if moves < 10 {
    // Endgame - potentially be more aggressive
    timeLimit = timeLimit * 4 / 3
}
```

## Stability Detection

A critical aspect of time management is detecting when search results have stabilized:

```go
// Position seems stable (same best move for consecutive iterations)
if line.moves[0] == tm.lastBestMove &&
   math.Abs(float64(line.score - tm.lastScore)) < 10 &&
   line.depth >= 10 {

    // Reduce allocated time
    optimum = optimum * 3 / 4
}
```

## Mate Search Optimization

When a mate is found, the engine can often stop searching earlier:

```go
// If we found a forced mate, no need to search deeper
if bestScore > MateScore-MaxDepth || bestScore < -MateScore+MaxDepth {
    // If the mate distance is short enough, we can stop
    mateInN := MateScore - abs(bestScore)
    if mateInN <= 5 { // Mate in 5 or less
        tm.cancel()
        return
    }
}
```

## Tactical Positions

For highly tactical positions (many captures, checks, etc.), the engine may
allocate extra time:

```go
// This is detected in the search by measuring:
// 1. The number of tactical moves in the position
// 2. Large score fluctuations
// 3. Frequent changes in the best move

if tacticalPosition {
    // Increase difficulty to allocate more time
    tm.difficulty = math.Min(tm.difficulty*1.5, maxDifficulty)
}
```

## UCI Integration

The time manager interfaces with the UCI protocol through the `go` command parameters:

```go
// In UCI protocol
func (uci *Protocol) goCommand(fields []string) error {
    limits := parseLimits(fields)
    ctx, cancel := context.WithCancel(context.TODO())

    // Create time manager with these limits
    tm := newTimeManager(ctx, time.Now(), limits, &currentBoard)

    // Start search with the time manager
    // ...
}
```

## Performance Impact

Effective time management can significantly impact engine strength:

1. **Optimal Allocation**: Spending more time on critical positions
2. **Time Saving**: Moving quickly in forced or simple positions
3. **Avoiding Time Pressure**: Maintaining a time buffer for late-game complications
4. **Consistent Playing Strength**: Balancing time usage across the game

The time manager's adaptive approach allows ArGO to maintain strong play
under various time controls, from bullet to classical games.
