// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package attacks

import (
	"testing"

	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGenerateKingAttacks tests the generateKingAttacks function for various scenarios
func TestGenerateKingAttacks(t *testing.T) {
	testCases := []struct {
		name           string
		square         int
		expectedBits   []int // Squares that should be attacked
		unexpectedBits []int // Squares that should NOT be attacked
	}{
		{
			name:           "King on A1 (bottom-left corner)",
			square:         A1,
			expectedBits:   []int{A2, B1, B2},
			unexpectedBits: []int{C1, D1},
		},
		{
			name:           "King on H8 (top-right corner)",
			square:         H8,
			expectedBits:   []int{H7, G8, G7},
			unexpectedBits: []int{F8, E8},
		},
		{
			name:           "King on D4 (middle of board)",
			square:         D4,
			expectedBits:   []int{C3, C4, C5, D3, D5, E3, E4, E5},
			unexpectedBits: []int{B2, B3, B4, F2, F3, F4},
		},
		{
			name:           "King on E1 (bottom edge)",
			square:         E1,
			expectedBits:   []int{D1, D2, E2, F1, F2},
			unexpectedBits: []int{G1, G2},
		},
		{
			name:           "King on A8 (top-left corner)",
			square:         A8,
			expectedBits:   []int{A7, B8, B7},
			unexpectedBits: []int{C8, D8},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attacks := generateKingAttacks(tc.square)

			// Check that the expected squares are attacked
			for _, expectedBit := range tc.expectedBits {
				if !attacks.Test(expectedBit) {
					t.Errorf("Expected square %d to be attacked, but it was not", expectedBit)
				}
			}

			// Check that unexpected squares are not attacked
			for _, unexpectedBit := range tc.unexpectedBits {
				if attacks.Test(unexpectedBit) {
					t.Errorf("Square %d should NOT be attacked, but it was", unexpectedBit)
				}
			}

			// Verify total number of attacked squares
			attackedSquaresCount := attacks.Count()
			if tc.square == A1 || tc.square == H1 || tc.square == A8 || tc.square == H8 {
				// Corner squares should have 3 attacked squares
				if attackedSquaresCount != 3 {
					t.Errorf(
						"Corner square %d should have 3 attacked squares, got %d",
						tc.square,
						attackedSquaresCount,
					)
				}
			} else if tc.square%8 == 0 || tc.square%8 == 7 || tc.square < 8 || tc.square >= 56 {
				// Edge squares should have 5 attacked squares
				if attackedSquaresCount != 5 {
					t.Errorf("Edge square %d should have 5 attacked squares, got %d", tc.square, attackedSquaresCount)
				}
			} else {
				// Middle squares should have 8 attacked squares
				if attackedSquaresCount != 8 {
					t.Errorf("Middle square %d should have 8 attacked squares, got %d", tc.square, attackedSquaresCount)
				}
			}
		})
	}
}

// TestInitKingAttacks tests the initialization of the KingAttacks lookup table
func TestInitKingAttacks(t *testing.T) {
	// Ensure table is empty before initialization
	for sq := 0; sq < 64; sq++ {
		if KingAttacks[sq] != 0 {
			t.Fatalf("KingAttacks table should be empty before initialization")
		}
	}

	// Initialize king attacks
	InitKingAttacks()

	// Verify table is populated correctly
	for sq := A8; sq <= H1; sq++ {
		attacks := KingAttacks[sq]

		// Verify non-zero attacks
		if attacks == 0 {
			t.Errorf("KingAttacks for square %d should not be zero", sq)
		}

		// Verify consistency with generateKingAttacks
		expectedAttacks := generateKingAttacks(sq)
		if attacks != expectedAttacks {
			t.Errorf("Inconsistent attacks for square %d", sq)
		}
	}
}

// Benchmark for generateKingAttacks to ensure good performance
func BenchmarkGenerateKingAttacks(b *testing.B) {
	testSquares := []int{A1, D4, H7, B2, E5, F6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, sq := range testSquares {
			generateKingAttacks(sq)
		}
	}
}

// TestKingAttackDirections ensures king attacks cover all 8 possible directions
func TestKingAttackDirections(t *testing.T) {
	middleSquare := D4
	attacks := generateKingAttacks(middleSquare)

	// Directional test cases
	directions := []struct {
		name     string
		offsetFn func(int) int
	}{
		{"Northwest", func(sq int) int { return sq + 7 }},
		{"North", func(sq int) int { return sq + 8 }},
		{"Northeast", func(sq int) int { return sq + 9 }},
		{"East", func(sq int) int { return sq + 1 }},
		{"West", func(sq int) int { return sq - 1 }},
		{"Southwest", func(sq int) int { return sq - 9 }},
		{"South", func(sq int) int { return sq - 8 }},
		{"Southeast", func(sq int) int { return sq - 7 }},
	}

	for _, dir := range directions {
		t.Run(dir.name, func(t *testing.T) {
			attackedSquare := dir.offsetFn(middleSquare)
			if !attacks.Test(attackedSquare) {
				t.Errorf("King should attack %s direction from square %d", dir.name, middleSquare)
			}
		})
	}
}
