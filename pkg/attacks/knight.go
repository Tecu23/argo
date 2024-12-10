// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

// KnightAttacks is a lookup table where KnightAttacks[sq] gives a bitboard of all squares
// a knight can attack from square sq. It is precomputed for fast lookups during move generation.
var KnightAttacks [64]bitboard.Bitboard

// InitKnightAttacks initializes the KnightAttacks table for all squares.
// Each entry is computed once, allowing O(1) move generation for knights later.
func InitKnightAttacks() {
	for sq := constants.A8; sq <= constants.H1; sq++ {
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
	// We apply file masks (e.g., ^constants.FileA) to prevent moves that
	// would cross file boundaries.
	// The shifts (e.g., constants.NW + constants.N) represent vertical and horizontal moves combined.
	// For example, NW+N means: move one rank up+left, then another rank up.

	// Up-Left moves:
	attacks |= (b & ^constants.FileA) << (constants.NW + constants.N)
	attacks |= (b & ^constants.FileA & ^constants.FileB) << (constants.NW + constants.W)

	// Down-Left moves:
	attacks |= (b & ^constants.FileA & ^constants.FileB) >> (constants.NE + constants.E)
	attacks |= (b & ^constants.FileA) >> (constants.NE + constants.N)

	// Up-Right moves:
	attacks |= (b & ^constants.FileH) >> (constants.NW + constants.N)
	attacks |= (b & ^constants.FileH & ^constants.FileG) >> (constants.NW + constants.W)

	// Down-Right moves:
	attacks |= (b & ^constants.FileH & ^constants.FileG) << (constants.NE + constants.E)
	attacks |= (b & ^constants.FileH) << (constants.NE + constants.N)

	return attacks
}
