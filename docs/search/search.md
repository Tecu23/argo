# Search Implementation

## Overview

The search system is the core of the ArGO chess engine, responsible for
exploring the game tree to find the best move. ArGO implements a
state-of-the-art alpha-beta search with numerous enhancements for
performance and accuracy.

## Alpha-Beta with Principal Variation Search

The main search algorithm is Principal Variation Search (PVS), a refinement
of alpha-beta pruning that's optimized for when the first move is
likely to be the best.

### Key Components

```go
func (e *Engine) alphaBeta(
    ctx context.Context,
    b *board.Board,
    depth, alpha, beta, ply int,
    tm *timeManager,
) int {
    // Transposition table lookup
    // Check extensions
    // Quiescence search for leaf nodes
    // Move generation and ordering
    // Principal Variation Search logic
    // Move loops with various pruning techniques
    // Transposition table storage
}
```

### Primary Search Optimizations

1. **Transposition Table**: Caches search results to avoid redundant work
2. **Check Extensions**: Increases search depth when the king is in check
3. **Quiescence Search**: Extends search in tactical positions
4. **Move Ordering**: Orders moves to maximize pruning efficiency
5. **Late Move Reduction**: Reduces search depth for unlikely moves
6. **Null Move Pruning**: Skips a turn to quickly detect strong positions
7. **Killer Move Heuristic**: Remembers moves that cause cutoffs
8. **History Heuristic**: Tracks effectiveness of quiet moves

## Root Search and Iterative Deepening

ArGO uses iterative deepening, starting from depth 1 and gradually increasing:

```go
func (e *Engine) search(ctx context.Context, b *board.Board, tm *timeManager) SearchInfo {
    // Initialize search variables

    // Iterative deepening loop
    for depth := 1; depth <= maxDepth; depth++ {
        // Root search at current depth
        score, mv := e.searchRoot(ctx, b, depth, tm)

        // Update best move and score
        bestMove = mv
        bestScore = score

        // Check if search should stop
        if tm.IsDone() || ctx.Err() != nil {
            break
        }

        // If found a forced mate, no need to search deeper
        if bestScore > MateScore-MaxDepth || bestScore < -MateScore+MaxDepth {
            break
        }
    }

    // Return search results
}
```

## Quiescence Search

The quiescence search prevents horizon effects by continuing to search tactical
positions even at the nominal end of the search depth:

```go
func (e *Engine) quiescence(
    ctx context.Context,
    b *board.Board,
    alpha, beta, ply int,
    tm *timeManager,
) int {
    // Stand-pat score
    score := e.evaluator.Evaluate(b)

    // Beta cutoff
    if score >= beta {
        return beta
    }

    // Update alpha if needed
    alpha = max(alpha, score)

    // Generate and order captures
    moves := b.GenerateCaptures()
    moves = e.orderMoves(moves, b, move.NoMove, ply)

    // Search capture moves
    for _, mv := range moves {
        // Make capture and search recursively
        // Update alpha if better move found
        // Beta cutoff if score too high
    }

    return alpha
}
```

## Move Ordering

Effective move ordering is critical for search efficiency. ArGO implements
multiple ordering heuristics:

```go
func (e *Engine) orderMoves(
    moves []move.Move,
    b *board.Board,
    ttMove move.Move,
    ply int,
) []move.Move {
    // Score and sort moves by priority:
    // 1. Transposition table move
    // 2. Captures (sorted by MVV-LVA)
    // 3. Killer moves
    // 4. History heuristic scores
}
```

The MVV-LVA (Most Valuable Victim - Least Valuable Aggressor) scoring prioritizes
capturing high-value pieces with low-value ones.

## Late Move Reduction

For moves analyzed later in the move list (which are less likely to be good
based on move ordering), the engine reduces the search depth:

```go
// Calculate reduction with adjustments
reduct := e.reductionTable.GetWithAdjustments(depth, moveCount, isPV, historyScore)

// Additional dynamic adjustments
if mv == e.killerMoves[ply][0] || mv == e.killerMoves[ply][1] {
    reduct = max(0, reduct-1)
}

// Ensure we don't reduce into the quiescence search
if depth-reduct <= 0 {
    reduct = max(0, depth-1)
}
```

The reduction amount is stored in a pre-computed table and adjusted based on
various factors like position type, move characteristics, and search history.

## Time Management

ArGO implements sophisticated time management to make the best use of allocated time:

```go
func (tm *timeManager) calculateTimeLimit(difficulty, branchFactor float64) time.Duration {
    // Calculate time allocation based on:
    // - Remaining time
    // - Increment
    // - Position complexity (difficulty)
    // - Game phase
    // - Move characteristics
}
```

The time manager dynamically adjusts the allocated time based on:

1. **Position Complexity**: Allocates more time for complex positions
2. **Best Move Stability**: If the best move keeps changing, allocates more time
3. **Score Fluctuation**: Allocates more time when the evaluation changes significantly
4. **Game Phase**: Different time allocation strategies for opening,
   middlegame, and endgame
5. **Move Characteristics**: Special handling for forced moves, captures, etc.

## Transposition Table

The transposition table caches search results to avoid redundant work:

```go
type TTEntry struct {
    Key      uint64    // Zobrist hash of the position
    Depth    int       // How deep we searched
    Score    int       // Position evaluation
    Flag     TTFlag    // Type of score (exact/upper/lower bound)
    BestMove move.Move // Best move found
    Age      uint8     // When this entry was created
}
```

Entry replacement uses a combination of depth and age criteria:

```go
func (tt *TranspositionTable) Store(key uint64, score, depth int, flag TTFlag, bestMove move.Move) {
    // Replacement strategy based on:
    // - Depth (prefer deeper searches)
    // - Age (prefer current search)
    // - Exact scores vs bounds
}
```

## Search Progress Monitoring

ArGO reports search progress to the UCI interface:

```go
func (e *Engine) createSearchInfo() SearchInfo {
    return SearchInfo{
        Score: UciScore{
            Centipawns: e.mainLine.score,
            Mate:       0,
        },
        Depth:    e.mainLine.depth,
        Nodes:    e.mainLine.nodes,
        Time:     time.Since(e.start),
        MainLine: e.mainLine.moves,
    }
}
```

This information is used by the UCI interface to display search statistics
and the principal variation.

## Special Search Features

1. **Mate Distance Pruning**: Prunes positions that can't improve the best mate found
2. **Draw Recognition**: Early detection of draws by insufficient material,
   repetition, etc.
3. **Internal Iterative Deepening**: Falls back to a shorter search if no
   best move is available
4. **Search Instability Detection**: Allocates more time when the search is unstable
5. **Adaptive Null Move Pruning**: Adjusts null move reduction based on position
   characteristics

## Endgame Tablebases

ArGO can be extended to support endgame tablebases for perfect play in positions
with few pieces. The search would need to be modified to:

1. Detect positions with few enough pieces to use tablebases
2. Query the tablebase instead of searching
3. Adjust search parameters based on tablebase information

## Performance Metrics

The search performance can be measured in:

1. **Nodes per Second**: Typically millions of nodes per second on modern hardware
2. **Effective Branching Factor**: Ideally close to 2 for well-pruned alpha-beta
3. **Depth Reached**: Typical depth for middle game positions in a few seconds
