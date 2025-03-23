// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import "github.com/Tecu23/argov2/pkg/bitboard"

// GetQueenAttacks combines bishop-like and rook-like attacks to produce queen attacks.
// A queen moves as both a bishop and a rook. We use previously computed bishop and rook masks,
// magic numbers, and lookup tables. The occupancy is processed to find the correct index,
// and we OR the bishop and rook attack bitboards to get the full set of queen moves.
func GetQueenAttacks(sq int, occupancy bitboard.Bitboard) bitboard.Bitboard {
	queenAttacks := bitboard.Bitboard(0)

	bishopOccupancies := occupancy
	rookOccupancies := occupancy

	// Derive bishop attacks from occupancy using magic indexing
	bishopOccupancies &= bishopMasks[sq]
	bishopOccupancies *= bishopMagicNumbers[sq]
	bishopOccupancies >>= 64 - bishopRelevantBits[sq]

	// Derive rook attacks from occupancy using magic indexing
	rookOccupancies &= rookMasks[sq]
	rookOccupancies *= rookMagicNumbers[sq]
	rookOccupancies >>= 64 - rookRelevantBits[sq]

	// Combine bishop and rook attacks for the queen
	queenAttacks = BishopAttacks[sq][bishopOccupancies] | RookAttacks[sq][rookOccupancies]

	return queenAttacks
}
