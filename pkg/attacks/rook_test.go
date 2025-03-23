// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGenerateRookPossibleBlockers tests the generation of possible blocker squares for rooks
func TestGenerateRookPossibleBlockers(t *testing.T) {
	testCases := []struct {
		name            string
		square          int
		expectedCount   int
		expectedSquares []int
	}{
		{
			name:          "Rook on D4 (middle board)",
			square:        D4,
			expectedCount: 10,
			expectedSquares: []int{
				D5, D6, D7, D3, D2, // Vertical
				E4, F4, G4, C4, B4, // Horizontal
			},
		},
		{
			name:          "Rook on A1 (corner)",
			square:        A1,
			expectedCount: 12,
			expectedSquares: []int{
				A2, A3, A4, A5, A6, A7, // Vertical
				B1, C1, D1, E1, F1, G1, // Horizontal
			},
		},
		{
			name:          "Rook on H8 (corner)",
			square:        H8,
			expectedCount: 12,
			expectedSquares: []int{
				H7, H6, H5, H4, H3, H2, // Vertical
				G8, F8, E8, D8, C8, B8, // Horizontal
			},
		},
	}

	InitSliderPiecesAttacks(Rook)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			blockers := generateRookPossibleBlockers(tc.square)

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

// TestGenerateRookAttacks tests the generation of rook attacks with various blocking scenarios
func TestGenerateRookAttacks(t *testing.T) {
	testCases := []struct {
		name              string
		square            int
		blockSquares      []int
		expectedAttacks   []int
		unexpectedAttacks []int
	}{
		{
			name:         "Rook on D4 with no blockers",
			square:       D4,
			blockSquares: []int{},
			expectedAttacks: []int{
				D1, D2, D3, D5, D6, D7, D8, // Vertical
				A4, B4, C4, E4, F4, G4, H4, // Horizontal
			},
			unexpectedAttacks: []int{
				E5, C5, E3, C3, // Diagonal squares
			},
		},
		{
			name:         "Rook on D4 with partial blocking",
			square:       D4,
			blockSquares: []int{D6, F4},
			expectedAttacks: []int{
				D6, D5, D3, D2, D1, // Vertical (blocked at D6)
				C4, B4, A4, E4, F4, // Horizontal (blocked at F4)
			},
			unexpectedAttacks: []int{
				D7, D8, // Beyond vertical block
				G4, H4, // Beyond horizontal block
			},
		},
		{
			name:         "Rook on A1 with blocking",
			square:       A1,
			blockSquares: []int{A3, C1},
			expectedAttacks: []int{
				A2, A3, // Vertical (blocked at A3)
				B1, C1, // Horizontal (blocked at C1)
			},
			unexpectedAttacks: []int{
				A4, A5, // Beyond vertical block
				D1, E1, // Beyond horizontal block
			},
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
			attacks := generateRookAttacks(tc.square, block)

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

// TestGetRookAttacks tests the GetRookAttacks function with magic bitboards
func TestGetRookAttacks(t *testing.T) {
	testCases := []struct {
		name          string
		square        int
		occupancy     bitboard.Bitboard
		expectedCount int
	}{
		{
			name:          "Rook on D4 with no occupancy",
			square:        D4,
			occupancy:     bitboard.Bitboard(0),
			expectedCount: 14, // Maximum possible attacks for a rook
		},
		{
			name:          "Rook on D4 with partial occupancy",
			square:        D4,
			occupancy:     bitboard.Bitboard(0x8002108000000), // Some occupied square
			expectedCount: 8,                                  // Reduced attacks due to blocking
		},
		{
			name:          "Rook on A1 with occupancy",
			square:        A1,
			occupancy:     bitboard.Bitboard(0x1000010000000000), // Block at A2
			expectedCount: 6,                                     // Limited by edge and blocking
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attacks := GetRookAttacks(tc.square, tc.occupancy)
			attackCount := attacks.Count()

			if attackCount != tc.expectedCount {
				t.Errorf("Expected %d attacked squares, got %d", tc.expectedCount, attackCount)
			}
		})
	}
}

// Benchmark for rook attack generation functions
func BenchmarkRookAttacks(b *testing.B) {
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
				generateRookAttacks(sq, occ)
				GetRookAttacks(sq, occ)
			}
		}
	}
}
