// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

// KingAttacks is a lookup table where KingAttacks[sq] gives a bitboard of all squares
// a king on square sq can move to. This is precomputed for all 64 squares for O(1) lookup.
var KingAttacks [64]bitboard.Bitboard

// InitKingAttacks initializes the KingAttacks array by generating the attack bitboard
// for a king placed on each square. After this, KingAttacks can be used directly
// to find king moves without recomputing them each time.
func InitKingAttacks() {
	for sq := constants.A8; sq <= constants.H1; sq++ {
		KingAttacks[sq] = generateKingAttacks(sq)
	}
}

// generateKingAttacks computes all possible king moves from the given square.
// The king moves one square in any direction (8 possibilities), but we must
// mask out moves that go off the board using file boundaries.
func generateKingAttacks(square int) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	b := bitboard.Bitboard(0)
	b.Set(square)

	// Each line corresponds to a direction the king can move:
	// NW, N, NE, E, W, SW, S, SE
	// File masks (^constants.FileA, ^constants.FileH) prevent wrapping around edges.
	attacks |= (b & ^constants.FileA) << constants.NW
	attacks |= b << constants.N
	attacks |= (b & ^constants.FileH) << constants.NE

	attacks |= (b & ^constants.FileH) << constants.E
	attacks |= (b & ^constants.FileA) >> constants.E

	attacks |= (b & ^constants.FileH) >> constants.NW
	attacks |= b >> constants.N
	attacks |= (b & ^constants.FileA) >> constants.NE

	return attacks
}
