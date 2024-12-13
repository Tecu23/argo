# Chess Engine Search Implementation Documentation

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Search Implementation](#search-implementation)
4. [UCI Integration](#uci-integration)
5. [Time Management](#time-management)
6. [Interaction Flow](#interaction-flow)

## Overview

This chess engine implements an alpha-beta search with iterative deepening,
quiescence search, and time management. The engine communicates via the UCI
protocol and manages its search depth based on available time.

## Architecture

### Core Components

```plaintext
engine/
├── engine.go         # Main engine structure and initialization
├── search.go        # Search implementation
├── evaluate.go      # Position evaluation
└── timemanager.go   # Time management
```

### Key Types and Structures

```go
// SearchInfo represents search results
type SearchInfo struct {
    Depth    int         // Current search depth
    Score    Score       // Position evaluation
    Nodes    uint64      // Nodes searched
    Time     Duration    // Time spent searching
    MainLine []move.Move // Principal variation
}

// Engine structure
type Engine struct {
    nodes     uint64
    mainLine  searchLine
    progress  func(SearchInfo)
    evaluate  func(board.Board) int
}

// Time manager structure
type timeManager struct {
    start             time.Time
    allocated         time.Duration
    maxTime          time.Duration
    limits           LimitsType
    side             bool
    lastScore        int
    scoreChanged     bool
    lastBestMove     move.Move
    searchStable     bool
    done             <-chan struct{}
    cancel           context.CancelFunc
}
```

## Search Implementation

### Iterative Deepening

```go
func (e *Engine) search(ctx context.Context, b board.Board, tm *timeManager) SearchInfo {
    // Initialize search
    e.nodes = 0
    var bestMove move.Move
    var bestMoveAtDepth move.Move
    var bestScoreAtDepth int

    // Iterative deepening loop
    for depth := 1; depth <= MaxDepth; depth++ {
        // Search at current depth
        bestScoreAtDepth = -Infinity
        bestMoveAtDepth = move.Move(0)

        // Root move search
        moves := b.GenerateMoves()
        for _, mv := range moves {
            score := -e.alphaBeta(ctx, b, depth-1, -Infinity, -bestScoreAtDepth)
            if score > bestScoreAtDepth {
                bestScoreAtDepth = score
                bestMoveAtDepth = mv
            }
        }

        // Update best move
        bestMove = bestMoveAtDepth

        // Update search info and check time
        e.updateSearchInfo(depth, bestMove, bestScoreAtDepth)
        if tm.shouldStop() {
            break
        }
    }
    return e.createSearchInfo()
}
```

### Alpha-Beta Search

```go
func (e *Engine) alphaBeta(ctx context.Context, b board.Board, depth, alpha, beta int) int {
    // Base cases
    if depth == 0 {
        return e.quiescence(ctx, b, alpha, beta)
    }

    // Generate and search moves
    moves := b.GenerateMoves()
    for _, mv := range moves {
        score := -e.alphaBeta(ctx, b, depth-1, -beta, -alpha)

        // Alpha-beta pruning
        if score >= beta {
            return beta
        }
        if score > alpha {
            alpha = score
        }
    }
    return alpha
}
```

### Quiescence Search

```go
func (e *Engine) quiescence(ctx context.Context, b board.Board, alpha, beta int) int {
    // Stand-pat evaluation
    standPat := e.evaluate(b)
    if standPat >= beta {
        return beta
    }
    alpha = max(alpha, standPat)

    // Search captures
    moves := b.GenerateCaptures()
    for _, mv := range moves {
        score := -e.quiescence(ctx, b, -beta, -alpha)
        if score >= beta {
            return beta
        }
        alpha = max(alpha, score)
    }
    return alpha
}
```

## UCI Integration

### Protocol Flow

1. UCI command received from GUI
2. Command parsed and validated
3. Search parameters extracted
4. Search initiated with time controls
5. Progress reported back to GUI
6. Best move sent when search completes

### Search Integration

```go
func (uci *Protocol) goCommand(fields []string) error {
    // Parse search limits
    limits := parseLimits(fields)

    // Create context and time manager
    ctx, cancel := context.WithCancel(context.TODO())
    tm := newTimeManager(limits, b.Side, done, cancel)

    // Start search
    go func() {
        result := uci.engine.Search(ctx, SearchParams{
            Boards: uci.boards,
            Limits: limits,
            Progress: func(si SearchInfo) {
                uci.engineOutput <- si
            },
        })
        uci.engineOutput <- result
        close(uci.engineOutput)
    }()

    return nil
}
```

## Time Management

### Time Allocation

```go
func (tm *timeManager) calculateTimeAllocation() {
    if tm.limits.MoveTime > 0 {
        tm.allocated = time.Duration(tm.limits.MoveTime) * time.Millisecond
        return
    }

    // Get base time and increment
    baseTime := tm.limits.WhiteTime
    increment := tm.limits.WhiteIncrement
    if tm.side {
        baseTime = tm.limits.BlackTime
        increment = tm.limits.BlackIncrement
    }

    // Calculate allocation
    movesToGo := tm.limits.MovesToGo
    if movesToGo == 0 {
        movesToGo = 30 // Default estimate
    }

    // Basic time management formula
    tm.allocated = time.Duration(baseTime/movesToGo) +
                  time.Duration(increment/2)
}
```

### Stop Conditions

```go
func (tm *timeManager) shouldStop() bool {
    // Check time usage
    elapsed := time.Since(tm.start)
    if elapsed >= tm.maxTime {
        return true
    }

    // Check search stability
    if elapsed >= tm.allocated &&
       (tm.searchStable || !tm.scoreChanged) {
        return true
    }

    return false
}
```

## Interaction Flow

1. **UCI Protocol Level**

   - Receives commands from GUI
   - Parses search parameters
   - Creates search context and time manager
   - Reports search progress
   - Sends best move when complete

2. **Engine Level**

   - Manages iterative deepening
   - Tracks best moves and scores
   - Updates search statistics
   - Checks time management conditions

3. **Search Level**

   - Implements alpha-beta algorithm
   - Handles move generation
   - Manages quiescence search
   - Evaluates positions

4. **Time Management Level**
   - Calculates time allocations
   - Monitors search stability
   - Decides when to stop search
   - Handles time controls

### Example Flow

1. GUI sends "go wtime 300000 btime 300000 winc 2000 binc 2000"
2. UCI protocol parses time controls
3. Time manager calculates allocation
4. Search begins with iterative deepening
5. Progress reported every iteration
6. Time manager monitors conditions
7. Search stops when criteria met
8. Best move sent to GUI

This implementation provides a balance between search depth and time
management while maintaining UCI protocol compliance.
