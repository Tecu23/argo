// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import "github.com/Tecu23/argov2/pkg/bitboard"

// GetRookAttacks returns the attack bitboard for a rook placed on square 'sq'
// given a particular occupancy bitboard of the board.
//
// Steps:
//  1. Intersection with RookMasks[sq] to isolate relevant occupancy bits.
//     This focuses on the subset of squares that affect rook moves from 'sq'.
//  2. Multiply by rookMagicNumbers[sq] (the magic number), which together with a shift
//     will map this occupancy pattern to a unique index.
//  3. Shift right by (64 - rookRelevantBits[sq]) to get the final index.
//  4. Use this index to retrieve the precomputed attack bitboard from RookAttacks[sq].
//
// This approach uses magic bitboards to achieve O(1) rook move generation.
func GetRookAttacks(sq int, occupancy bitboard.Bitboard) bitboard.Bitboard {
	// calculate magic index
	occupancy &= rookMasks[sq]
	occupancy *= rookMagicNumbers[sq]
	occupancy >>= 64 - rookRelevantBits[sq]

	return RookAttacks[sq][occupancy]
}

// generateRookPossibleBlockers computes the potential blocker squares of a rook placed on a given square,
// We only iterate till the second to last in each direction, because a blocker in the last square will not
// affect the possible moves generation. This gice 12 possible blockers with 2^12 possible blocker combinations
func generateRookPossibleBlockers(square int) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// mask relevant bishop occupancy bits
	for r = tr + 1; r <= 6; r++ {
		attacks |= 1 << (r*8 + tf)
	}

	for f = tf + 1; f <= 6; f++ {
		attacks |= 1 << (tr*8 + f)
	}

	for r = tr - 1; r >= 1; r-- {
		attacks |= 1 << (r*8 + tf)
	}

	for f = tf - 1; f >= 1; f-- {
		attacks |= 1 << (tr*8 + f)
	}

	return attacks
}

// generateRookAttacks computes rook attacks for a given square and a given
// occupancy bitboard "block" that represents which squares are occupied. The rook
// attack rays stop when they encounter a block.
//
// Parameters:
// - square: the position of the rook (0-based index on an 8x8 board with A8=0).
// - block: bitboard of occupied squares that can block the rook.
func generateRookAttacks(square int, block bitboard.Bitboard) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// mask relevant bishop occupancy bits
	for r = tr + 1; r <= 7; r++ {
		attacks |= 1 << (r*8 + tf)
		if (1<<(r*8+tf))&block != 0 {
			break
		}
	}

	for f = tf + 1; f <= 7; f++ {
		attacks |= 1 << (tr*8 + f)
		if (1<<(tr*8+f))&block != 0 {
			break
		}
	}

	for r = tr - 1; r >= 0; r-- {
		attacks |= 1 << (r*8 + tf)
		if (1<<(r*8+tf))&block != 0 {
			break
		}
	}

	for f = tf - 1; f >= 0; f-- {
		attacks |= 1 << (tr*8 + f)
		if (1<<(tr*8+f))&block != 0 {
			break
		}
	}

	return attacks
}
