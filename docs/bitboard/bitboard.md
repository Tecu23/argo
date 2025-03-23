# Bitboard Implementation

## Overview

The `bitboard` package is a fundamental building block of the ArGO chess engine.
It provides a highly efficient representation of the chess board using a 64-bit
integer, where each bit corresponds to a specific square on the board.

## Bitboard Representation

A bitboard is a 64-bit unsigned integer (`uint64`) where each bit represents one
square on a chess board. This representation allows for extremely fast board
manipulation using bitwise operations.

```txt
8 | 0  1  2  3  4  5  6  7
7 | 8  9  10 11 12 13 14 15
6 | 16 17 18 19 20 21 22 23
5 | 24 25 26 27 28 29 30 31
4 | 32 33 34 35 36 37 38 39
3 | 40 41 42 43 44 45 46 47
2 | 48 49 50 51 52 53 54 55
1 | 56 57 58 59 60 61 62 63
  -------------------------
    a  b  c  d  e  f  g  h
```

In this layout, bit 0 corresponds to a8, bit 7 to h8, and bit 63 to h1. This allows
for efficient mapping between bit indices and algebraic chess notation.

## Key Operations

### Bit Manipulation

The `bitboard` package provides several methods for manipulating individual bits:

- `Set(pos int)`: Sets the bit at position `pos` to 1
- `Clear(pos int)`: Sets the bit at position `pos` to 0
- `Test(pos int)`: Tests if the bit at position `pos` is set (returns a boolean)

### Bitboard Scanning

Efficient scanning operations are critical for move generation:

- `FirstOne()`: Returns the position of the least significant set bit (LSB)
  and clears it
- `LastOne()`: Returns the position of the most significant set bit (MSB)
  and clears it
- `Count()`: Returns the number of set bits in the bitboard

The `FirstOne()` implementation is particularly important as it's used extensively
in move generation loops. It uses Go's built-in `bits.TrailingZeros64` function
for optimal performance, which maps to processor-specific instructions on
supported platforms (like the BSF instruction on x86).

## Performance Considerations

Bitboards allow for extremely efficient chess operations:

1. **Move Generation**: Multiple potential moves can be generated simultaneously
   using bitwise operations.
2. **Board Evaluation**: Piece positions and patterns can be quickly
   assessed using bitwise logic.
3. **Attack Detection**: Determining attacked squares can be done using
   pre-computed attack tables combined with bitboard operations.

## Example Usage

```go
// Create a bitboard with a white pawn on e2
bb := bitboard.Bitboard(0)
bb.Set(E2) // Set the bit corresponding to e2

// Check if e4 is empty
if !bb.Test(E4) {
    // e4 is empty, can potentially move there
}

// Count pawns on the board
pawnCount := bb.Count()

// Iterate through all pieces
for bb != 0 {
    square := bb.FirstOne() // Get and clear the LSB
    // Process piece at square
}
```

## Bitboard Constants

The engine defines several useful bitboard constants:

- File masks (FileA through FileH)
- Rank masks (Row1 through Row8)
- Special masks for board regions (center, kingside, queenside, etc.)

These constants are used for pattern recognition and move generation logic.

## Advantages Over Traditional Board Representation

Compared to traditional array-based board representations, bitboards offer:

1. Much faster move generation
2. More efficient board state evaluation
3. Lower memory usage
4. Better cache locality
5. Natural parallelism through bitwise operations

These advantages make bitboards the standard representation in high-performance
chess engines like ArGO.
