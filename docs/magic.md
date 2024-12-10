# Overview

**Magic numbers** are a key optimization technique used in bitboard-based chess
engines to generate attacks for sliding pieces (bishops, rooks, and queens) efficiently.
The idea behind magic numbers is to use a carefully chosen bitwise multiplication
(a "magic" factor) followed by a shift to map any possible board occupancy pattern
in a piece's sliding direction(s) to a unique index. This index is then used to quickly
retrieve a precomputed set of attack squares from a lookup table.

In simpler terms:

1. For each square on the board and for each sliding piece type (rook or bishop),
   we define a "mask" of squares that can potentially be attacked by that piece
   from that position if the board were empty.
2. The actual attacked squares depend on which of those masked squares are occupied.
3. To avoid computing this from scratch every time, we use a "magic number" and a
   small lookup table. By combining the occupancy pattern of these masked squares
   with the magic number, we get an index that directly leads us to a
   precomputed attack bitboard.

This approach dramatically speeds up move generation because it reduces what would
otherwise be a complex series of bit operations into a single multiplication, shift,
and array indexing operation.

---

## Terminology and Concepts

### Bitboards

A **bitboard** is a 64-bit integer where each bit corresponds to a square on
the chessboard. Using this representation, set bits represent occupied squares or
attacked squares, depending on the context. They allow extremely fast operations
using standard CPU instructions like shifts and bitwise AND, OR, XOR.

### Sliding Pieces

Bishops, rooks, and queens are considered "sliding" pieces. They move along straight
lines (diagonals for bishops, ranks and files for rooks, both for queens) until they
reach the edge of the board or are blocked by another piece. Determining which squares
they attack given an arbitrary position involves checking along these
lines and stopping at blockers.

### Attack Masks

For each square `S` and piece type (bishop or rook), we predefine an **attack mask**
that represents all the squares along the sliding lines from `S` to the edges of
the board, excluding the square `S` itself. This mask is the "maximum" area that
piece could attack if there were no pieces in the way.

For rooks, the mask is all squares in the same rank and file
(excluding the square itself). For bishops, it’s all squares along the
diagonals passing through that square. For each square, these attack masks are
fixed and can be precomputed once.

### Occupancy and Blockers

**Occupancy** refers to which squares in that attack mask are actually occupied
by pieces. Because pieces can block sliding moves, not all squares in the attack
mask will be attacked. The actual attacked squares depend on which squares along
that line are occupied:

- If a rook on D4 sees no pieces on the same rank or file, it attacks all squares
  along that rank and file.
- If there is a piece blocking on the same file, the rook's attacks stop before
  that piece, and the squares beyond it are not attacked.

To handle this generically, we consider all subsets of the attack mask’s squares.
Each subset represents one possible occupancy scenario along that line.
By enumerating these subsets, we can precompute the resulting attacked squares
for every possible configuration of blockers.

### Magic Numbers

A **magic number** is a 64-bit constant chosen specifically for each square and piece
type (rook or bishop) to create a perfect (collision-free) hash function from the
occupancy pattern to the index in the precomputed attack table.

**How Magic Numbers Work:**

1. **Mask Relevant Squares**: Start with the full attack mask for that square.
   Consider only these "relevant" squares that affect the piece’s sliding attacks.
2. **Occupancy Subsets**: Each subset of these relevant squares corresponds to one
   possible pattern of blockers. There can be up to
   2^(number_of_relevant_squares) such subsets.
3. **Indexing With Magic Multiplication**: For a given occupancy pattern:

   - Extract the subset of relevant squares that are occupied.
   - Convert that subset into an integer index by:
     - Taking the bitboard of occupied squares (just the ones in the relevant mask).
     - Multiplying this occupancy bitboard by a predetermined magic number.
     - Shifting the result right by (64 - number_of_relevant_squares).

   This process yields a unique index (no collisions) into the attack lookup table.

4. **Magic Numbers Are Carefully Chosen**: Magic numbers must be found
   (usually via a brute-force or a precomputation algorithm) so that for every
   possible occupancy pattern, the indexing yields unique results with no overlaps.
   Once a suitable magic number is found for each square and piece type, it is hard-coded.

### Lookup Tables

For each square `S` and piece type (rook or bishop), we have:

- A **mask** of relevant squares (attackMask).
- A **magic number** (RookMagicNumbers[sq] or BishopMagicNumbers[sq]).
- A **relevant bits** count that tells how many bits in the mask are relevant.
- A **lookup table** (an array of bitboards) of size 2^(relevant_bits). Each entry
  in this table corresponds to one subset of these relevant squares. The value stored
  is the attacked squares for that occupancy.

The process at runtime:

1. Given a board position, compute the occupancy along that piece's line by intersecting
   the current position bitboard with the mask.
2. Use the `SetOccupancy` function (or equivalent) to find the subset index of
   that occupancy.
3. Multiply by the magic number, shift, and use the result to index into the
   attack table.
4. Return the precomputed attack bitboard stored there. This bitboard now tells you
   exactly which squares are attacked by that piece from that square given the
   current blockers.

### The `SetOccupancy` Function

The `SetOccupancy` function takes an index (representing which subset of the
mask we want), the number of bits in the mask, and the mask itself. It returns
a bitboard that represents a particular occupancy pattern chosen by `index`.
This is used during initialization to generate all subsets of the attack mask
and fill the lookup table.

During engine initialization or precomputation phase:

1. Enumerate all subsets of the mask. For each `index` from 0 to 2^(bits_in_mask)-1:
   - `SetOccupancy` determines which bits in the mask are included in this subset.
   - Compute the resulting attacked squares by simulating the piece's move
     along lines until blocked by these occupancy squares.
   - Store that attack bitboard in the lookup table at the index
     determined by `(occupancy * magic_number) >> (64 - relevant_bits)`.

After this precomputation, every position's rook or bishop attacks from a
given square can be found in O(1) time.

---

## Example

Consider a rook on square D4 (assume D4 maps to some index `sq`):

1. Identify the rook’s mask for `sq` which includes all squares in D4’s rank and
   file, excluding the square itself.
2. Determine the occupancy pattern of those squares from the current board position.
3. Multiplying this occupancy pattern by the rook’s magic number for `sq` and
   shifting appropriately yields an index.
4. Use this index to quickly lookup the precomputed attack bitboard.
5. The returned bitboard shows all squares the rook attacks.

No loops or complicated logic at runtime—just one multiplication, one shift,
and one array access.

---

## Advantages of Using Magic Numbers

- **Speed**: Magic indexing replaces a complex computation with a constant-time lookup.
- **Simplicity at Runtime**: After the initial precomputation, the engine’s move
  generation is streamlined.
- **Memory vs. Speed Tradeoff**: Magic bitboards use more memory for tables but yield
  faster move generation.

---

## Conclusion

Magic numbers are a clever technique to transform the problem of dynamic sliding
piece move generation into a constant-time lookup by precomputing attack patterns
for every possible occupancy configuration along relevant rays. By pairing carefully
chosen magic numbers with occupancy subsets, this approach avoids collision and ensures
that each unique occupancy pattern maps to a unique entry in a precomputed attack
table, enabling very fast and efficient chess move generation.
