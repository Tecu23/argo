// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// KnightAttacks is a lookup table where KnightAttacks[sq] gives a bitboard of all squares
// a knight can attack from square sq. It is precomputed for fast lookups during move generation.
var KnightAttacks [64]bitboard.Bitboard

// InitKnightAttacks initializes the KnightAttacks table for all squares.
// Each entry is computed once, allowing O(1) move generation for knights later.
func InitKnightAttacks() {
	for sq := A8; sq <= H1; sq++ {
		KnightAttacks[sq] = generateKnightAttacks(sq)
	}
}

// generateKnightAttacks computes a bitboard of attacks for a knight placed on the given square.
// Knights move in an "L" shape: two squares in one direction and one square perpendicular.
// The code uses bit-shifts and file masks to ensure moves don't wrap around the board.
func generateKnightAttacks(square int) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	b := bitboard.Bitboard(0)
	b.Set(square)

	// Each of these lines represents a possible knight move direction.
	// We apply file masks (e.g., ^FileA) to prevent moves that
	// would cross file boundaries.
	// The shifts (e.g., NW + N) represent vertical and horizontal moves combined.
	// For example, NW+N means: move one rank up+left, then another rank up.

	// Up-Left moves:
	attacks |= (b & ^FileA) << (7 + 8)
	attacks |= (b & ^FileA & ^FileB) << (7 + -1)

	// Down-Left moves:
	attacks |= (b & ^FileA & ^FileB) >> (9 + 1)
	attacks |= (b & ^FileA) >> (9 + 8)

	// Up-Right moves:
	attacks |= (b & ^FileH) >> (7 + 8)
	attacks |= (b & ^FileH & ^FileG) >> (7 + -1)

	// Down-Right moves:
	attacks |= (b & ^FileH & ^FileG) << (9 + 1)
	attacks |= (b & ^FileH) << (9 + 8)

	return attacks
}
