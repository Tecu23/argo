// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGenerateBishopPossibleBlockers tests the generation of possible blocker squares for bishops
func TestGenerateBishopPossibleBlockers(t *testing.T) {
	testCases := []struct {
		name            string
		square          int
		expectedCount   int
		expectedSquares []int
	}{
		{
			name:            "Bishop on D4 (middle board)",
			square:          D4,
			expectedCount:   9,
			expectedSquares: []int{C5, E5, C3, E3, B6, F6, B2, F2},
		},
		{
			name:            "Bishop on A1 (bottom-left corner)",
			square:          A1,
			expectedCount:   6,
			expectedSquares: []int{B2, C3, D4},
		},
		{
			name:            "Bishop on H8 (top-right corner)",
			square:          H8,
			expectedCount:   6,
			expectedSquares: []int{G7, F6, E5},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			blockers := generateBishopPossibleBlockers(tc.square)

			// Count the number of blocker squares
			blockerCount := blockers.Count()
			if blockerCount != tc.expectedCount {
				t.Errorf("Expected %d blocker squares, got %d", tc.expectedCount, blockerCount)
			}

			// Verify specific blocker squares
			for _, sq := range tc.expectedSquares {
				if !blockers.Test(sq) {
					t.Errorf("Expected square %d to be a possible blocker, but it was not", sq)
				}
			}
		})
	}
}

// TestGenerateBishopAttacks tests the generation of bishop attacks with various blocking scenarios
func TestGenerateBishopAttacks(t *testing.T) {
	testCases := []struct {
		name              string
		square            int
		blockSquares      []int
		expectedAttacks   []int
		unexpectedAttacks []int
	}{
		{
			name:         "Bishop on D4 with no blockers",
			square:       D4,
			blockSquares: []int{},
			expectedAttacks: []int{
				C5, B6, A7, E5, F6, G7, H8,
				C3, B2, A1, E3, F2, G1,
			},
			unexpectedAttacks: []int{D4, D5, D3, C4, E4},
		},
		{
			name:         "Bishop on D4 with partial blocking",
			square:       D4,
			blockSquares: []int{F6, B2},
			expectedAttacks: []int{
				C5, B6, A7, E5, F6,
				C3, B2, E3, F2, G1,
			},
			unexpectedAttacks: []int{G7, H8, A1},
		},
		{
			name:              "Bishop on A1 with blocking",
			square:            A1,
			blockSquares:      []int{B2},
			expectedAttacks:   []int{B2},
			unexpectedAttacks: []int{C3, D4},
		},
		{
			name:              "Bishop on H8 with partial blocking",
			square:            H8,
			blockSquares:      []int{F6},
			expectedAttacks:   []int{G7, F6},
			unexpectedAttacks: []int{E5, D4},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create block bitboard from block squares
			block := bitboard.Bitboard(0)
			for _, sq := range tc.blockSquares {
				block.Set(sq)
			}

			// Generate attacks
			attacks := generateBishopAttacks(tc.square, block)

			// Verify expected attacks
			for _, sq := range tc.expectedAttacks {
				if !attacks.Test(sq) {
					t.Errorf("Expected square %d to be attacked, but it was not", sq)
				}
			}

			// Verify unexpected attacks
			for _, sq := range tc.unexpectedAttacks {
				if attacks.Test(sq) {
					t.Errorf("Square %d should NOT be attacked, but it was", sq)
				}
			}
		})
	}
}

// TestGetBishopAttacks tests the GetBishopAttacks function
func TestGetBishopAttacks(t *testing.T) {
	// Test cases with various occupancy scenarios
	testCases := []struct {
		name          string
		square        int
		occupancy     bitboard.Bitboard
		expectedCount int
	}{
		{
			name:          "Bishop on D4 with no occupancy",
			square:        D4,
			occupancy:     bitboard.Bitboard(0),
			expectedCount: 13, // Maximum possible attacks for a bishop in the middle of the board
		},
		{
			name:          "Bishop on D4 with partial occupancy",
			square:        D4,
			occupancy:     bitboard.Bitboard(0x120000004200000), // Some occupied square
			expectedCount: 8,                                    // Reduced attacks due to blocking
		},
		{
			name:          "Bishop on A1 with occupancy",
			square:        A1,
			occupancy:     bitboard.Bitboard(0x2000000000000), // Some occupied square
			expectedCount: 1,                                  // Very limited moves due to board edge
		},
	}

	// Init Magics Before Testing
	InitSliderPiecesAttacks(Bishop)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attacks := GetBishopAttacks(tc.square, tc.occupancy)

			attackCount := attacks.Count()

			if attackCount != tc.expectedCount {
				t.Errorf("Expected %d attacked squares, got %d", tc.expectedCount, attackCount)
			}
		})
	}
}

// Benchmark for bishop attack generation functions
func BenchmarkBishopAttacks(b *testing.B) {
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
				generateBishopAttacks(sq, occ)
				GetBishopAttacks(sq, occ)
			}
		}
	}
}
