// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package attacks

import (
	"fmt"
	"testing"

	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGetQueenAttacks tests the GetQueenAttacks function with various scenarios
func TestGetQueenAttacks(t *testing.T) {
	testCases := []struct {
		name            string
		square          int
		occupancy       bitboard.Bitboard
		expectedSquares []int
		blockedSquares  []int
	}{
		{
			name:      "Queen on D4 with no occupancy",
			square:    D4,
			occupancy: bitboard.Bitboard(0),
			expectedSquares: []int{
				// Horizontal moves
				A4, B4, C4, E4, F4, G4, H4,
				// Vertical moves
				D1, D2, D3, D5, D6, D7, D8,
				// Diagonal moves
				C5, B6, A7, E5, F6, G7, H8,
				C3, B2, A1, E3, F2, G1,
			},
			blockedSquares: []int{
				D4,             // Own square
				A5, B5, E6, F7, // Off-board diagonals
			},
		},
		{
			name:   "Queen on D4 with partial occupancy",
			square: D4,
			occupancy: func() bitboard.Bitboard {
				occ := bitboard.Bitboard(0)
				occ.Set(D6) // Blocker on vertical
				occ.Set(F4) // Blocker on horizontal
				occ.Set(F6) // Blocker on diagonal
				return occ
			}(),
			expectedSquares: []int{
				// Horizontal (blocked at F4)
				C4, B4, A4, E4, F4,
				// Vertical (blocked at D6)
				D3, D2, D1, D5, D6,
				// Diagonal (blocked at F6)
				C5, B6, E5, F6,
				C3, B2, A1, E3, F2, G1,
			},
			blockedSquares: []int{
				D7, D8, // Beyond vertical blocker
				G4, H4, // Beyond horizontal blocker
				G7, H8, // Beyond diagonal blocker
			},
		},
		{
			name:   "Queen on A1 with occupancy",
			square: A1,
			occupancy: func() bitboard.Bitboard {
				occ := bitboard.Bitboard(0)
				occ.Set(A3) // Vertical blocker
				occ.Set(C1) // Horizontal blocker
				occ.Set(C3) // Diagonal blocker
				return occ
			}(),
			expectedSquares: []int{
				A2, A3, // Vertical until blocker
				B1, C1, // Horizontal until blocker
				B2, C3, // Diagonal until blocker
			},
			blockedSquares: []int{
				A4, A5, A6, A7, A8, // Beyond vertical blocker
				D1, E1, F1, G1, H1, // Beyond horizontal blocker
				D4, E5, F6, G7, H8, // Beyond diagonal blocker
			},
		},
		{
			name:   "Queen on H8 with occupancy",
			square: H8,
			occupancy: func() bitboard.Bitboard {
				occ := bitboard.Bitboard(0)
				occ.Set(H6) // Vertical blocker
				occ.Set(F8) // Horizontal blocker
				occ.Set(F6) // Diagonal blocker
				return occ
			}(),
			expectedSquares: []int{
				H7, H6, // Vertical until blocker
				G8, F8, // Horizontal until blocker
				G7, F6, // Diagonal until blocker
			},
			blockedSquares: []int{
				H5, H4, H3, H2, H1, // Beyond vertical blocker
				E8, D8, C8, B8, A8, // Beyond horizontal blocker
				E5, D4, C3, B2, A1, // Beyond diagonal blocker
			},
		},
	}

	InitSliderPiecesAttacks(Bishop)
	InitSliderPiecesAttacks(Rook)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attacks := GetQueenAttacks(tc.square, tc.occupancy)

			// Check expected attacked squares
			for _, sq := range tc.expectedSquares {
				if !attacks.Test(sq) {
					t.Errorf("Expected square %d to be attacked, but it was not", sq)
				}
			}

			// Check squares that should be blocked
			for _, sq := range tc.blockedSquares {
				if attacks.Test(sq) {
					t.Errorf("Square %d should NOT be attacked, but it was", sq)
				}
			}
		})
	}
}

// TestQueenAttacksConsistency verifies that queen attacks are consistent with
// the combination of bishop and rook attacks
func TestQueenAttacksConsistency(t *testing.T) {
	testSquares := []int{A1, D4, H8, E2, B7}
	testOccupancies := []bitboard.Bitboard{
		bitboard.Bitboard(0),
		bitboard.Bitboard(1) << 10,
		bitboard.Bitboard(1) << 30,
		bitboard.Bitboard(0xFFFF), // Some complex occupancy
	}

	for _, sq := range testSquares {
		for _, occ := range testOccupancies {
			t.Run(fmt.Sprintf("Square-%d-Occupancy-%v", sq, occ), func(t *testing.T) {
				queenAttacks := GetQueenAttacks(sq, occ)
				bishopAttacks := GetBishopAttacks(sq, occ)
				rookAttacks := GetRookAttacks(sq, occ)

				// Queen attacks should be exactly bishop OR rook attacks
				expectedAttacks := bishopAttacks | rookAttacks
				if queenAttacks != expectedAttacks {
					t.Errorf("Queen attacks inconsistent with bishop|rook attacks at square %d", sq)
				}

				// Verify attack count
				queenCount := queenAttacks.Count()
				combinedCount := (bishopAttacks | rookAttacks).Count()
				if queenCount != combinedCount {
					t.Errorf("Attack count mismatch: queen=%d, bishop|rook=%d",
						queenCount, combinedCount)
				}
			})
		}
	}
}

// BenchmarkGetQueenAttacks measures performance of queen attack generation
func BenchmarkGetQueenAttacks(b *testing.B) {
	testSquares := []int{A1, D4, H7, B2, E5, F6}
	testOccupancies := []bitboard.Bitboard{
		bitboard.Bitboard(0),
		bitboard.Bitboard(1) << 10,
		bitboard.Bitboard(1) << 30,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, sq := range testSquares {
			for _, occ := range testOccupancies {
				GetQueenAttacks(sq, occ)
			}
		}
	}
}

// TestQueenAttacksWithMagicBitboards specifically tests the magic bitboard lookup
// functionality for queen moves
func TestQueenAttacksWithMagicBitboards(t *testing.T) {
	testCases := []struct {
		name    string
		square  int
		setupFn func() bitboard.Bitboard
		verify  func(t *testing.T, attacks bitboard.Bitboard)
	}{
		{
			name:   "Magic index calculation - Empty board",
			square: D4,
			setupFn: func() bitboard.Bitboard {
				return bitboard.Bitboard(0)
			},
			verify: func(t *testing.T, attacks bitboard.Bitboard) {
				expectedCount := 27 // Maximum queen moves from D4
				if count := attacks.Count(); count != expectedCount {
					t.Errorf("Expected %d attacks, got %d", expectedCount, count)
				}
			},
		},
		{
			name:   "Magic index calculation - Complex position",
			square: E4,
			setupFn: func() bitboard.Bitboard {
				occ := bitboard.Bitboard(0)
				// Set up complex position with multiple blockers
				blockers := []int{E6, G4, C4, E2, G6, C2}
				for _, sq := range blockers {
					occ.Set(sq)
				}
				return occ
			},
			verify: func(t *testing.T, attacks bitboard.Bitboard) {
				// Verify specific blocked squares
				blockedSquares := []int{E7, E8, H4, A4, E1, H7, A1}
				for _, sq := range blockedSquares {
					if attacks.Test(sq) {
						t.Errorf("Square %d should be blocked", sq)
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			occupancy := tc.setupFn()
			attacks := GetQueenAttacks(tc.square, occupancy)
			tc.verify(t, attacks)
		})
	}
}
