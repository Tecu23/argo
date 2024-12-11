package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGeneratePawnAttacks tests the generatePawnAttacks function for various scenarios
func TestGeneratePawnAttacks(t *testing.T) {
	testCases := []struct {
		name         string
		square       int
		side         color.Color
		expectedBits []int // Squares that should be attacked
	}{
		// White pawn attack tests
		{
			name:         "White pawn on A7 (corner)",
			square:       A7,
			side:         color.WHITE,
			expectedBits: []int{B8},
		},
		{
			name:         "White pawn on H7 (corner)",
			square:       H7,
			side:         color.WHITE,
			expectedBits: []int{G8},
		},
		{
			name:         "White pawn on D4 (middle board)",
			square:       D4,
			side:         color.WHITE,
			expectedBits: []int{C5, E5},
		},
		{
			name:         "White pawn on E2 (near promotion rank)",
			square:       E2,
			side:         color.WHITE,
			expectedBits: []int{D3, F3},
		},

		// Black pawn attack tests
		{
			name:         "Black pawn on A2 (corner)",
			square:       A2,
			side:         color.BLACK,
			expectedBits: []int{B1},
		},
		{
			name:         "Black pawn on H2 (corner)",
			square:       H2,
			side:         color.BLACK,
			expectedBits: []int{G1},
		},
		{
			name:         "Black pawn on D5 (middle board)",
			square:       D5,
			side:         color.BLACK,
			expectedBits: []int{C4, E4},
		},
		{
			name:         "Black pawn on E7 (near promotion rank)",
			square:       E7,
			side:         color.BLACK,
			expectedBits: []int{D6, F6},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			attacks := generatePawnAttacks(tc.square, tc.side)

			// Check that the expected squares are attacked
			for _, expectedBit := range tc.expectedBits {
				if !attacks.Test(expectedBit) {
					t.Errorf("Expected square %d to be attacked, but it was not", expectedBit)
				}
			}

			if attacks.Count() != len(tc.expectedBits) {
				t.Errorf(
					"Expected to be only %d attacks when attacking from %d",
					len(tc.expectedBits),
					tc.square,
				)
			}
		})
	}
}

// TestInitPawnAttacks tests the initialization of the PawnAttacks lookup table
func TestInitPawnAttacks(t *testing.T) {
	// Ensure table is empty before initialization
	for side := 0; side < 2; side++ {
		for sq := 0; sq < 64; sq++ {
			if PawnAttacks[side][sq] != 0 {
				t.Fatalf("PawnAttacks table should be empty before initialization")
			}
		}
	}

	// Initialize pawn attacks
	InitPawnAttacks()

	// Verify table is populated correctly
	for side := range []color.Color{color.WHITE, color.BLACK} {
		for sq := A8; sq <= H1; sq++ {
			attacks := PawnAttacks[side][sq]

			// Verify non-zero attacks
			if (color.Color(side) == color.WHITE && sq > H8) ||
				(color.Color(side) == color.BLACK && sq < A1) {
				if attacks == 0 {
					t.Errorf("PawnAttacks for %v on square %d should not be zero",
						[]string{"WHITE", "BLACK"}[side], sq)
				}
			}

			// Verify consistency with generatePawnAttacks
			expectedAttacks := generatePawnAttacks(
				sq,
				[]color.Color{color.WHITE, color.BLACK}[side],
			)
			if attacks != expectedAttacks {
				t.Errorf("Inconsistent attacks for %v on square %d",
					[]string{"WHITE", "BLACK"}[side], sq)
			}
		}
	}
}

// Benchmark for generatePawnAttacks to ensure good performance
func BenchmarkGeneratePawnAttacks(b *testing.B) {
	testSquares := []int{A1, D4, H7, B2, E5, F6}
	testColors := []color.Color{color.WHITE, color.BLACK}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, sq := range testSquares {
			for _, c := range testColors {
				generatePawnAttacks(sq, c)
			}
		}
	}
}
