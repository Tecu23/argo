# Board Representation with Bitboards

This document describes the board representation strategy used
in this chess engine, focusing on the `Bitboard` type and how
it underpins the move generation and evaluation logic.

## Overview

In a chess engine, the representation of the board plays a crucial role in
performance and complexity. This project uses **bitboards** as a core data structure.
A bitboard is essentially a 64-bit integer where each bit corresponds to a single
square on the chessboard. By leveraging bit-level operations, we can perform
extremely fast computations for move generation, attack detection, and evaluation.

**Key advantages of using bitboards:**

1. **Performance**: Bitwise operations (AND, OR, XOR, shifts) are among the fastest
   available on modern CPUs. This makes operations like checking attacks, occupancy,
   and mobility very efficient.
2. **Simplicity in Calculation**: With precomputed attack masks, determining whether
   a piece can move to a particular square or generate all potential moves can
   be done with a few bitwise operations.
3. **Compactness**: A 64-bit integer is a compact representation of the entire 8x8
   board state for a particular attribute (e.g., piece placement of a certain type).

## What is a Bitboard?

A bitboard is a 64-bit unsigned integer (`uint64` in Go) where the bit positions
map to squares on a chessboard. Typically, we number squares from 0 to 63.
There are many conventions, but a common one is:

- The least significant bit (bit 0) corresponds to one corner
  of the board (e.g., `a1`).
- The most significant bit (bit 63) corresponds to the opposite corner (`h8`).

For example:

```markdown
63.......................................0
'-------------------------------'
8 | bit 63 ... ... bit 56 | Rank 8
7 | bit 55 ... ... bit 48 | Rank 7
6 | bit 47 ... ... bit 40 | Rank 6
5 | bit 39 ... ... bit 32 | Rank 5
4 | bit 31 ... ... bit 24 | Rank 4
3 | bit 23 ... ... bit 16 | Rank 3
2 | bit 15 ... ... bit 08 | Rank 2
1 | bit 07 ... ... bit 00 | Rank 1
'a b c d e f g h'
```

In this coordinate system, each rank and file maps to a known set of bits,
enabling quick operations like shifting a bitboard to represent movement.

## Use in the Chess Engine

Each type of piece and each color may be represented by its own bitboard.
For instance, you might have:

- `whitePawns` (a bitboard with bits set where white pawns are located)
- `blackQueens`
- `allPieces` (a bitboard combining all occupied squares)
- `whitePieces`
- `blackPieces`

This allows quick queries. For example, to see if a certain square
is occupied, you can test a single bit on `allPieces`. To generate
moves for a rook, you can use `rookAttacks[square] & ~whitePieces`
to find all legal moves for a white rook from that square
(assuming `rookAttacks[square]` is a precomputed bitboard of
attack squares from that square).

## The `Bitboard` Type

In this codebase, `Bitboard` is defined as:

```go
type Bitboard uint64
```

It comes with several helper methods to manipulate bits:

### `Set(pos int)`

Set the bit at position `pos` to `1`. If `pos` corresponds to,
for example, square `a1`, then after calling `b.Set(0)`,
bit 0 of the bitboard `b` is now set.

### `Test(pos int) bool`

Check if the bit at position `pos` is set. This quickly tells you if a
particular square is occupied by the piece corresponding to that bitboard.

### `Clear(pos int)`

Clear the bit at position `pos`, setting it to `0`.

### `Count() int`

Return the number of bits set to `1` in the bitboard. Useful for
counting how many pieces are on a certain bitboard, how many moves are
generated, or during evaluation.

### `FirstOne()` and `LastOne()`

These functions find and remove either the least significant set bit
(`FirstOne()`) or the most significant set bit (`LastOne()`), returning
their positions. They are very handy in loops where you repeatedly extract
one piece or one move at a time from a bitboard until it’s empty.

For example, to iterate over all set bits in a `Bitboard`:

```go
bb := someBitboard
for {
    pos := bb.FirstOne()
    if pos == 64 {
        break // no more bits
    }
    // pos is the index of a set bit
}
```

### `String()`

Returns a 64-character binary string representing the bitboard, from the most
significant bit to the least significant bit. This is primarily for
debugging and visualization.

### `PrintBitboard()`

Prints the bitboard in a human-readable 8x8 grid format, showing
which squares are set. This makes it easy to visualize the board state for debugging.

## Example Usage

**Setting and checking bits:**

```go
var b Bitboard
b.Set(0)
// Set square a1
b.Set(63)
// Set square h8
if b.Test(0) {
    fmt.Println("Square a1 is occupied.")
}
if !b.Test(32) {
    fmt.Println("Square e5 is not occupied.")
}

```

**Counting pieces:**

```go
count := b.Count()
// returns how many squares are set
fmt.Printf("There are %d pieces on the board.\n", count)
```

**Extracting positions:**

```go
var b Bitboard
b.Set(10)
b.Set(20)
b.Set(30)
for {
    pos := b.FirstOne()
    if pos == 64 {
        break
    }
    fmt.Printf("Found a piece at position %d\n", pos)
}
```

## Integration With Move Generation

Bitboards become invaluable when generating moves. For sliding pieces
(rooks, bishops, queens), we maintain precomputed tables of moves from
each square. Given a board configuration, we can look up a piece’s possible moves
via a few bitwise operations. For example:

- To generate rook moves from square `s`, start with a precomputed bitboard `rookAttacks[s]`.
- Mask out squares occupied by your own pieces (`& ~whitePieces`).
- Use the resulting bitboard to get all possible moves in one go.

For pawns, you can quickly determine legal captures by shifting a pawn’s bitboard
(to simulate moves) and intersecting with the opponent’s pieces bitboard to see
which capture squares are available.

## Debugging and Visualization

`PrintBitboard()` is a simple way to visualize the bitboard as an 8x8 grid.
It can show you exactly which squares are currently set, making it easier to
debug issues in move generation or evaluation code.

For example:

```go
b.PrintBitboard()
```

Might produce output that looks like:

```css
8   0 0 0 0 0 0 0 1
7   0 0 0 0 0 0 0 0 6
... ...
1   1 0 0 0 0 0 0 0
a b c d e f g h`
```

indicating set bits at `a1` and `h8`.

## Summary

- **What:** Bitboards are 64-bit integers representing the state of a chessboard.
- **Why:** They offer fast, compact, and efficient calculation
  for move generation and evaluation.
- **How:** Each bit corresponds to a square. We use precomputed attacks,
  simple bit operations, and helper methods to manipulate and query these boards.

This bitboard representation forms the foundation of our engine’s move generation
and evaluation strategy, enabling fast and efficient chess computations.
