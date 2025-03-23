# Magic Bitboards

## Overview

The ArGO chess engine uses "magic bitboards" for extremely efficient move
generation, particularly for sliding pieces (bishops, rooks, and queens).
This advanced technique provides O(1) move generation through a clever
combination of perfect hashing and look-up tables.

## The Problem with Sliding Pieces

Generating moves for sliding pieces (bishops, rooks, queens) is challenging
because their movement depends on the current board occupancy. A piece can
move until it hits another piece or the edge of the board. Traditional approaches
require expensive ray-tracing or loop-based algorithms.

## Magic Bitboards Solution

Magic bitboards provide a perfect hash function that maps from the relevant occupancy
bits to an index in a pre-computed attack table. This approach consists of
several key components:

### 1. Relevant Occupancy Masks

For each square and piece type, we first identify the "relevant occupancy" bits -
these are squares that can potentially block the piece's movement. Edge squares
are excluded from relevant occupancy since they can never block a ray (they can
only be destinations).

For example, for a rook on e4, the relevant occupancy squares are all squares on
the same rank and file, excluding e4 itself and the edge squares (a4, h4, e1, e8).

### 2. Magic Numbers

The key innovation in magic bitboards is finding a "magic number" for each square
and piece type combination. A magic number is a 64-bit value with the special
property that:

```txt
index = (occupancy * magic) >> shift
```

This hash function must produce a unique index for each relevant occupancy pattern,
with no collisions. Finding suitable magic numbers is a non-trivial task that
typically involves a randomized search.

The `findMagicNumbers` function in the `attacks` package implements this search process:

```go
func findMagicNumbers(square, relevantBits int, piece int) bitboard.Bitboard {
    // ... (search for a magic number that creates a perfect hash)
}
```

Once found, these magic numbers are hardcoded in the engine for performance.

### 3. Attack Tables

For each square and piece type, the engine pre-computes attack bitboards for all
possible occupancy patterns:

```go
func InitSliderPiecesAttacks(piece int) {
    for sq := A8; sq <= H1; sq++ {
        // ... (initialize masks)
        occupancyVariations := 1 << bitCount

        for count := 0; count < occupancyVariations; count++ {
            occupancy := SetOccupancy(count, bitCount, attackMask)

            if piece == Bishop {
                magicIndex := occupancy * bishopMagicNumbers[sq] >> (64 - bishopRelevantBits[sq])
                BishopAttacks[sq][magicIndex] = generateBishopAttacks(sq, occupancy)
            } else {
                magicIndex := occupancy * rookMagicNumbers[sq] >> (64 - rookRelevantBits[sq])
                RookAttacks[sq][magicIndex] = generateRookAttacks(sq, occupancy)
            }
        }
    }
}
```

### 4. Move Generation at Runtime

During actual gameplay, generating moves for a sliding piece becomes as simple as:

```go
func GetBishopAttacks(sq int, occupancy bitboard.Bitboard) bitboard.Bitboard {
    // Extract only the relevant occupancy bits
    occupancy &= bishopMasks[sq]

    // Compute the magic index
    occupancy *= bishopMagicNumbers[sq]
    occupancy >>= 64 - bishopRelevantBits[sq]

    // Look up the pre-computed attacks
    return BishopAttacks[sq][occupancy]
}
```

For queens, the engine combines bishop and rook attacks:

```go
func GetQueenAttacks(sq int, occupancy bitboard.Bitboard) bitboard.Bitboard {
    return GetBishopAttacks(sq, occupancy) | GetRookAttacks(sq, occupancy)
}
```

## Memory Requirements

Magic bitboards require significant memory for the attack tables:

- Bishop tables: 64 squares × 512 possible occupancy patterns = 32,768 bitboards
- Rook tables: 64 squares × 4,096 possible occupancy patterns = 262,144 bitboards

Each bitboard is 8 bytes, resulting in a total of approximately 2.3 MB for
the attack tables. This is a reasonable trade-off for the performance gained.

## Performance Benefits

The performance advantage of magic bitboards is substantial:

1. **Constant-time move generation**: O(1) regardless of board complexity
2. **Cache efficiency**: Pre-computed tables have good locality
3. **Branch-free code**: Move generation involves minimal branching
4. **SIMD-friendly**: Bitboard operations map well to modern CPU capabilities

## Implementation Details

ArGO's magic bitboard implementation includes:

1. **Pre-computation at startup**: Attack tables are initialized when the
   engine starts
2. **Hardcoded magic numbers**: Pre-computed magic numbers are stored as constants
3. **Separate tables for bishops and rooks**: Queens use a combination of both
4. **Relevant bits optimization**: Each square uses the minimum
   required relevant bits

This approach provides a perfect balance between initialization time, memory usage,
and runtime performance.
