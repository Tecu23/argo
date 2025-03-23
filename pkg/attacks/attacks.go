// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// InitSliderPiecesAttacks initializes the magic-based attacks for sliding pieces (either bishop or rook).
// This function:
// 1. Precomputes masks and occupancy variations for each square.
// 2. Generates all possible attack configurations for each occupancy variation.
// 3. Populates the BishopAttacks or RookAttacks lookup tables accordingly.
//
// 'piece' parameter should be Bishop or Rook.
// After this initialization, attack lookups for those pieces are O(1).
func InitSliderPiecesAttacks(piece int) {
	// Loop over all 64 squares on the board.
	for sq := A8; sq <= H1; sq++ {

		// Initialize bishop & rook masks for the square.
		// These masks represent all squares that can potentially affect the bishop/rook's moves from 'sq'.
		bishopMasks[sq] = generateBishopPossibleBlockers(sq)
		rookMasks[sq] = generateRookPossibleBlockers(sq)

		// Determine which mask to use based on the piece type.
		// For bishops, use the bishop mask; for rooks, use the rook mask.
		var attackMask bitboard.Bitboard
		if piece == Bishop {
			attackMask = generateBishopPossibleBlockers(sq)
		} else {
			attackMask = generateRookPossibleBlockers(sq)
		}

		// Count how many bits are set in the attackMask (i.e., how many squares are relevant).
		// This determines the number of occupancy variations we must consider.
		bitCount := attackMask.Count()

		// The number of occupancy variations is 2^(bitCount).
		// Each variation represents a different subset of those squares being occupied.
		occupancyVariations := 1 << bitCount

		// For each occupancy variation, we:
		// 1. Generate the occupancy subset with SetOccupancy().
		// 2. Compute a magic index to index into our precomputed attack table.
		// 3. Store the resulting attack bitboard in BishopAttacks or RookAttacks.
		for count := 0; count < occupancyVariations; count++ {
			// Build an occupancy bitboard for this particular subset of the attackMask.
			occupancy := SetOccupancy(count, bitCount, attackMask)
			if piece == Bishop {

				// Compute the magic index for the bishop based on occupancy.
				magicIndex := occupancy * bishopMagicNumbers[sq] >> (64 - bishopRelevantBits[sq])
				// Store the generated attacks in the bishop attack table.
				BishopAttacks[sq][magicIndex] = generateBishopAttacks(sq, occupancy)
			} else {

				// Compute the magic index for the rook based on occupancy.
				magicIndex := occupancy * rookMagicNumbers[sq] >> (64 - rookRelevantBits[sq])
				// Store the generated attacks in the rook attack table.
				RookAttacks[sq][magicIndex] = generateRookAttacks(sq, occupancy)
			}
		}
	}
}
