# Late Move Reduction (LMR)

## Overview

Late Move Reduction (LMR) is an optimization technique widely used in chess engines
to improve search efficiency and playing strength. Traditional alpha-beta search
tries to explore a vast number of moves. LMR helps prune or reduce the search depth
for moves that are less likely to yield better results, enabling the engine to search
effectively deeper in the more promising lines.

## How LMR Works

1. **Initial Concept**: In alpha-beta search, the engine considers moves in a certain
   order. Moves likely to improve the position (like principal variation moves, captures
   of high-value pieces, or transposition table hits) are examined first. Later moves
   in the ordering are often less interesting, as the engine has already seen several
   better moves.

2. **"Late" Moves**: When the engine reaches these later moves—often after it has
   already tried several more promising candidates—it's increasingly unlikely that
   a late move will surpass the best move found so far. LMR leverages this insight.

3. **Reducing the Depth**: Instead of searching these late moves at the full remaining
   depth, the engine reduces the search depth by a small amount (1-3 plies typically).
   This results in a cheaper, shallower probe into that move's line of play.

4. **Verification**: If the reduced-depth search unexpectedly finds that this "late"
   move might be good (for example, it fails high, indicating it's better than the
   current best), the engine re-searches the move at the original depth. This
   two-step approach ensures that potential gems are not prematurely discarded.

## Example

1. The engine reaches depth 10 and is examining the 10th move in a position.
   The first 9 moves were already considered, and none suggested that exploring
   every branch of this 10th move in depth is necessary.
2. Instead of fully exploring this move at the full 10 plies, the engine only
   searches it at (10 - 1 = 9) plies or less, saving computation time.
3. If this reduced search shows no promise, the engine moves on quickly, having
   saved time. If it looks promising, the engine invests more time and re-searches
   it at full depth.

## Why LMR is Useful

- **Performance Boost**: By reducing the amount of time spent on unpromising moves,
  LMR lowers the average branching factor of the search. This allows the engine to
  reach deeper plies in the same amount of time, improving overall playing strength.
- **Better Resource Allocation**: Resources (time and computation) get concentrated
  on moves more likely to influence the outcome. LMR helps the engine avoid wasting
  effort on lines that aren’t likely to improve the best score.

- **Scalability**: As the engine grows more complex, implementing advanced heuristics
  and deeper searches become increasingly expensive. LMR helps maintain manageability
  by ensuring that the engine doesn’t get bogged down examining every move equally.

## Key Points to Remember

- **Condition-Based**: LMR is often applied only if certain conditions are met—such
  as the move appearing late in the order, sufficient search depth, and not being
  a capture or check move.
- **Tuning-Dependent**: The specific reduction values and conditions for applying
  LMR are highly engine-dependent. They must be tuned via testing and empirical
  evidence.
- **Complementary**: LMR works best alongside other heuristics like good move
  ordering (principal variation and transposition table moves first), history
  heuristics, and killer moves.

---

In summary, Late Move Reduction is a heuristic designed to make the search more
efficient by cutting corners on less likely moves, thus allowing an engine to
focus more time and effort on the most promising parts of the search tree.
This leads to stronger play and better use of available computational resources.
