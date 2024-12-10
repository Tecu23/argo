// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
)

// PawnAttacks is a lookup table that stores the attack bitboards for pawns
// depending on their color and square. For example, PawnAttacks[WHITE][sq]
// gives a bitboard of all squares attacked by a white pawn placed on square `sq`.
var PawnAttacks [2][64]bitboard.Bitboard

// InitPawnAttacks precomputes the pawn attacks for both white and black pawns
// on every square. This is done at startup for fast lookups during move generation.
func InitPawnAttacks() {
	for sq := constants.A8; sq <= constants.H1; sq++ {
		PawnAttacks[color.WHITE][sq] = generatePawnAttacks(sq, color.WHITE)
		PawnAttacks[color.BLACK][sq] = generatePawnAttacks(sq, color.BLACK)
	}
}

// generatePawnAttacks returns the bitboard of attacks for a pawn of a given side
// on a specific square. White pawns attack diagonally forward-left (NW) and
// forward-right (NE). Black pawns attack diagonally backward-left (SW) and
// backward-right (SE), when viewed from White's perspective.
//
// The function considers file boundaries to ensure pawns don't wrap around the board.
func generatePawnAttacks(square int, side color.Color) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	b := bitboard.Bitboard(0)
	b.Set(square)

	if side == color.WHITE {
		// White pawn attacks move up and diagonally (towards higher ranks)
		// SW attack (if not on file A)
		attacks |= (b & ^constants.FileA) >> constants.NE
		// SE attack (if not on file H)
		attacks |= (b & ^constants.FileH) >> constants.NW
	} else {
		// Black pawn attacks move down and diagonally (towards lower ranks)
		// NW attack (if not on file A)
		attacks |= (b & ^constants.FileA) << constants.NW
		// NE attack (if not on file H)
		attacks |= (b & ^constants.FileH) << constants.NE
	}

	return attacks
}
