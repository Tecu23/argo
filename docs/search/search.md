# Chess Engine Search Implementation Documentation

## Table of Contents

1. [Overview](#overview)
2. [Search Implementation](#search-implementation)

## Overview

This chess engine implements an alpha-beta search with iterative deepening,
quiescence search, and time management. The engine communicates via the UCI
protocol and manages its search depth based on available time.

## Search Implementation

### Alpha-Beta Search

Alpha-Beta pruning is a search algorithm. It enhances the basic **minimax**
algorithm by eliminating branches in the game tree that do not need to be
explored because they cannot effect the final decision

#### Minimax Overview

Before diving into the alpha-beta pruning, we need to explore the minimax
algorithm. This is used to find the optimal move for a player assuming the
opponent also plays optimally. In a game tree:

- **Maximizer** tries to maximize the score,
- **Minimizer** tries to minimize the score.

The algorithm recursively explores all possible moves to find the optimal
strategy

#### Alpha-Beta Pruning

Alpha-beta pruning improves the minimax by "pruning" branches of the game
tree that cannot influence the final decision. This reduces the number of
nodes evaluated, making the algorithm much more efficient.

**Key Concepts:**

- **Alpha**: The best value that the maximizer can guarantee so far.
- **Beta**: The best value that the minimizer can guarantee so far.

As the algorithm explores the game tree:

- Alpha starts at −∞ and increases as the maximizer finds better options.
- Beta starts at +∞ and decreases as the minimizer finds better options.
- If at any point, α ≥ β, further exploration of the branch is unnecessary
  because it cannot affect the outcome

##### **Steps of the algorithm**

1. Start with Minimax Logic
   - Traverse the game tree recursively, alternating between maximizing
     and minimizing at each level.
2. Track Alpha and Beta:
   - At each node:
     - If it is a maximizer's turn, update α as the highest value found so far.
     - If it is a minimizer's turn, update β as the lowest value found so far.
3. Prune Unnecessary Branches:
   - If the value of α at a maximizer nodes becomes greater than or equal to β
     at a minimizer node, stop the exploring that branch.

##### Example

```mathematica
          Max
       /       \
     Min       Min
    / | \     / | \
   3  5  2   1  9  8
```

- the maximizer starts at the root.
- the minimizer controls the second level.

Without pruning:

- Minimax evaluates all 6 leaf nodes (3, 5, 2, 1, 9, 8)

With Alpha-Beta Pruning:

1. Start at the root node (Max).
2. Evaluate the first child (Min).
   - Explore its children: 3, 5, 2.
   - The best value for Min is 2.
3. Move to the second child (Min).
   - While evaluating its first child (1), the alorithm realizes:
     - For Max, 2 is already better than any potential value from this subtree
       (as it is the minimizer's turn, and 1 ≤ 2)
     - The branch is pruned, and values 9 and 8 are never evaluated.

**Result**: Only 4 leaf nodes are evaluated instead of 6.

#### Negamax Algoritm

The negamax algorithm is a simplified variant of the Minimax algorithm that
leverages symmetry between the maximizing and minimizing players in 2-player
zero-sum games.

Instead of handling separate logic for the maximizer and minimizer, Negamax
assumes that the score for one player is the negative of the score for the
opponent. This simplifies the algorithm, as we can maximize one score (the player's)
and negate the opponent's score in a single recursive function.

##### Negamax with Alpha-Beta Pruning

Negamax can be combined with alpha-beta pruning to prune unnecessary branches
efficiently. The same alpha (α) and beta (β) concepts from alpha-beta pruning
are used:

- α: The best score the maximizing player (or current player) is guaranteed.
- β: The best score the minimizing player (or opponent) is guaranteed.

##### **Steps of Negamax with Alpha-Beta Pruning**

1. Recursive Structure:
   - At each node, recursively evaluate child nodes and return the maximum
     negated value.
2. Alpha-Beta Bounds:
   - Pass -β and α- to the recursive calls for child nodes.
   - Swap α and β, and negate them because the opponent's bounds are symmetric
     to the current player's bounds

##### Algorithm

```go
func negamax(b *board.Board, depth int, alpha int, beta int) {
    if depth == 0 {
        return evaluate(b)
    }

    moves = generateMoves(b)
    for _, mv := range moves {
        copyB := b.copyBoard()
        if !copyB.MakeMove(mv, board.AllMoves) {
            continue
        }
        // Negate the value from the opponent's perspective
        score := -negamax(copyB, depth - 1, -beta, -alpha)
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

### Iterative Deepening

**Iterative Deepening** is a search strategy often used in chess engines to
combine the benefits of depth-limited search and time efficiency. It involves
performing a depth-first search (DFS) with increasing depth limits, starting
from depth 1 and incrementally going deeper until a given time constraint is
reached or a desired search depth is achieved.

#### How It Works

1. **Start with a Shallow Depth:**
   - Begin by searching the game tree to a depth of 1 ply (one move).
2. Increase Depth:
   - Repeatedly increase the depth by 1 and search the game tree again,
     reusing information from previous searches to optimize the process.
3. Use Move Ordering:
   - After each iteration, sort the moves based on their scores to
     prioritize the most promising moves in the next iteration. This makes
     alpha-beta pruning more effective.
4. Stop on Time or Depth:
   - If the allocated time runs out or the desired maximum depth is reached,
     stop and return the best move found so far.

#### Example

```go

func search(ctx context.Context, b *board.Board, tm *timeManager) move.Move {
    var bestMove move.Move
    var bestScore int

    maxDepth := max(MaxDepth, tm.limits.Depth)

    for depth := 1; depth <= maxDepth; depth++ {
        score, mv := searchRoot(ctx, b, depth)

        if mv != move.NoMove {
            bestMove = mv
            bestScore = score
        }

        // Check if time is up
        if tm.IsDone() {
            break
        }

        // Update time manager
        tm.OnNodesChanged(int(e.nodes))
    }

    return bestMove
}

```

### Quiescence Search

### Principal Variation Search

### Transposition Table

### Late Move Reduction

### Null Move Pruning

### Aspiration Windows
