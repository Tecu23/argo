package attacks

import (
	"testing"

	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGenerateKnightAttacks tests the generateKnightAttacks function for various scenarios
func TestGenerateKnightAttacks(t *testing.T) {
	testCases := []struct {
		name           string
		square         int
		expectedBits   []int // Squares that should be attacked
		unexpectedBits []int // Squares that should NOT be attacked
		expectedCount  int   // Total number of attacked squares
	}{
		{
			name:           "Knight on A1 (bottom-left corner)",
			square:         A1,
			expectedBits:   []int{B3, C2},
			unexpectedBits: []int{A2, A3, B1, B2, C1, D1},
			expectedCount:  2,
		},
		{
			name:           "Knight on H8 (top-right corner)",
			square:         H8,
			expectedBits:   []int{F7, G6},
			unexpectedBits: []int{H7, H6, G8, F8, E8, D8},
			expectedCount:  2,
		},
		{
			name:           "Knight on D4 (middle of board)",
			square:         D4,
			expectedBits:   []int{B3, B5, C2, C6, E2, E6, F3, F5},
			unexpectedBits: []int{D3, D5, C4, E4, B4, F4, A4, G4},
			expectedCount:  8,
		},
		{
			name:           "Knight on E2 (bottom edge)",
			square:         E2,
			expectedBits:   []int{C1, C3, D4, F4, G1, G3},
			unexpectedBits: []int{E1, E3, D2, F2, A2, B2, H2},
			expectedCount:  6,
		},
		{
			name:           "Knight on A8 (top-left corner)",
			square:         A8,
			expectedBits:   []int{B6, C7},
			unexpectedBits: []int{A7, A6, B8, C8, D8},
			expectedCount:  2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attacks := generateKnightAttacks(tc.square)

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
			if attackedSquaresCount != tc.expectedCount {
				t.Errorf(
					"Expected %d attacked squares, got %d",
					tc.expectedCount,
					attackedSquaresCount,
				)
			}
		})
	}
}

// TestInitKnightAttacks tests the initialization of the KnightAttacks lookup table
func TestInitKnightAttacks(t *testing.T) {
	// Ensure table is empty before initialization
	for sq := 0; sq < 64; sq++ {
		if KnightAttacks[sq] != 0 {
			t.Fatalf("KnightAttacks table should be empty before initialization")
		}
	}

	// Initialize knight attacks
	InitKnightAttacks()

	// Verify table is populated correctly
	for sq := A8; sq <= H1; sq++ {
		attacks := KnightAttacks[sq]

		// Verify non-zero attacks
		if attacks == 0 {
			t.Errorf("KnightAttacks for square %d should not be zero", sq)
		}

		// Verify consistency with generateKnightAttacks
		expectedAttacks := generateKnightAttacks(sq)
		if attacks != expectedAttacks {
			t.Errorf("Inconsistent attacks for square %d", sq)
		}
	}
}

// Benchmark for generateKnightAttacks to ensure good performance
func BenchmarkGenerateKnightAttacks(b *testing.B) {
	testSquares := []int{A1, D4, H7, B2, E5, F6}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, sq := range testSquares {
			generateKnightAttacks(sq)
		}
	}
}

// TestKnightAttackDirections ensures knight attacks cover all 8 possible L-shaped moves
func TestKnightAttackDirections(t *testing.T) {
	middleSquare := D4
	attacks := generateKnightAttacks(middleSquare)

	// Directional test cases for knight's L-shaped moves
	directions := []struct {
		name     string
		offsetFn func(int) int
	}{
		{"North-Northwest", func(sq int) int { return sq + 15 }}, // 2 up, 1 left
		{"North-Northeast", func(sq int) int { return sq + 17 }}, // 2 up, 1 right
		{"East-Northeast", func(sq int) int { return sq + 10 }},  // 1 up, 2 right
		{"East-Southeast", func(sq int) int { return sq - 6 }},   // 1 down, 2 right
		{"South-Southeast", func(sq int) int { return sq - 15 }}, // 2 down, 1 right
		{"South-Southwest", func(sq int) int { return sq - 17 }}, // 2 down, 1 left
		{"West-Southwest", func(sq int) int { return sq - 10 }},  // 1 down, 2 left
		{"West-Northwest", func(sq int) int { return sq + 6 }},   // 1 up, 2 left
	}

	for _, dir := range directions {
		t.Run(dir.name, func(t *testing.T) {
			attackedSquare := dir.offsetFn(middleSquare)
			if !attacks.Test(attackedSquare) {
				t.Errorf("Knight should attack %s direction from square %d", dir.name, middleSquare)
			}
		})
	}
}

// TestCrossFileBoundaries ensures knight attacks respect file boundaries
func TestCrossFileBoundaries(t *testing.T) {
	testCases := []struct {
		name           string
		square         int
		expectedBits   []int
		unexpectedBits []int
	}{
		{
			name:           "Knight on A-file",
			square:         A4,
			expectedBits:   []int{B6, C5, C3, B2},
			unexpectedBits: []int{A5, A3, A6, A2, A7, A1},
		},
		{
			name:           "Knight on H-file",
			square:         H4,
			expectedBits:   []int{F5, G6, F3, G2},
			unexpectedBits: []int{H5, H3, H6, H2, H7, H1},
		},
		{
			name:           "Knight on B-file (left edge)",
			square:         B3,
			expectedBits:   []int{A5, C5, D4, D2, A1, C1},
			unexpectedBits: []int{B4, B2, B5, B1},
		},
		{
			name:           "Knight on G-file (right edge)",
			square:         G3,
			expectedBits:   []int{E4, E2, F5, H5, F1, H1},
			unexpectedBits: []int{G4, G2, G5, G1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attacks := generateKnightAttacks(tc.square)

			for _, expectedBit := range tc.expectedBits {
				if !attacks.Test(expectedBit) {
					t.Errorf("Expected square %d to be attacked, but it was not", expectedBit)
				}
			}

			for _, unexpectedBit := range tc.unexpectedBits {
				if attacks.Test(unexpectedBit) {
					t.Errorf("Square %d should NOT be attacked, but it was", unexpectedBit)
				}
			}
		})
	}
}
