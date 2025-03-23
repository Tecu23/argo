// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import "github.com/Tecu23/argov2/pkg/bitboard"

// GetBishopAttacks returns the attack bitboard for a bishop placed on square 'sq'
// given a particular occupancy bitboard of the board.
//
// Steps are similar to GetRookAttacks but use bishop-specific masks, magic numbers,
// and relevant bits count:
// 1. Intersection with BishopMasks[sq] to consider only relevant occupancy bits.
// 2. Multiply by bishopMagicNumbers[sq] and shift to obtain a unique index.
// 3. Use that index to look up the precomputed attack bitboard from BishopAttacks[sq].
//
// This provides O(1) bishop move generation after precomputation.
func GetBishopAttacks(sq int, occupancy bitboard.Bitboard) bitboard.Bitboard {
	// calculate magic index
	occupancy &= bishopMasks[sq]
	occupancy *= bishopMagicNumbers[sq]
	occupancy >>= 64 - bishopRelevantBits[sq]

	return BishopAttacks[sq][occupancy]
}

// benerateBishopPossibleBlockers computes the potential blocker squares of a bishop placed on a given square,
// We only iterate till the second to last in each direction, because a blocker in the last square will not
// affect the possible moves generation. This gice 9 possible blockers with 2^9 possible blocker combinations
func generateBishopPossibleBlockers(square int) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// Move diagonally (up-right)
	for r, f = tr+1, tf+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		attacks |= 1 << (r*8 + f)
	}

	// Move diagonally (up-left)
	for r, f = tr-1, tf+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		attacks |= 1 << (r*8 + f)
	}

	// Move diagonally (down-right)
	for r, f = tr+1, tf-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		attacks |= 1 << (r*8 + f)
	}

	// Move diagonally (down-left)
	for r, f = tr-1, tf-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		attacks |= 1 << (r*8 + f)
	}

	return attacks
}

// generateBishopAttacks computes the potential attack squares of a bishop placed on a given square,
// considering the given bitboard "block" which represents occupied squares. If a blocked square is encountered,
// the bishop's attack ray stops in that direction.
//
// Parameters:
// - square: The position of the bishop on a 0-63 indexed board.
// - block:  A bitboard representing occupied squares that might block the bishop.
func generateBishopAttacks(square int, block bitboard.Bitboard) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// up-right (increasing rank and file)
	for r, f = tr+1, tf+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	// up-left (decreasing rank, increasing file)
	for r, f = tr-1, tf+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	// down-right (increasing rank, decreasing file)
	for r, f = tr+1, tf-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	// down-left (decreasing rank and file)
	for r, f = tr-1, tf-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	return attacks
}
