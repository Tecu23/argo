// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestRandomNumberGeneration tests the random number generation functions
func TestRandomNumberGeneration(t *testing.T) {
	randomState := uint32(1804289383) // Same seed as used in findMagicNumbers

	testCases := []struct {
		name     string
		numCalls int
	}{
		{"Single number generation", 1},
		{"Multiple number sequence", 100},
		{"Extended sequence", 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seen := make(map[uint64]bool)

			for i := 0; i < tc.numCalls; i++ {
				num := generateRandomUint64Number(&randomState)

				// Verify number isn't all zeros or all ones
				if num == 0 {
					t.Error("Generated zero value")
				}
				if num == 0xFFFFFFFFFFFFFFFF {
					t.Error("Generated all ones")
				}

				// Check for duplicates (though unlikely with 64-bit numbers)
				if seen[num] {
					t.Error("Generated duplicate number")
				}
				seen[num] = true
			}
		})
	}
}

// TestMagicNumberGeneration tests the magic number generation function
func TestMagicNumberGeneration(t *testing.T) {
	randomState := uint32(1804289383)

	testCases := []struct {
		name     string
		numTests int
	}{
		{"Single magic generation", 1},
		{"Multiple magic generations", 100},
		{"Extended magic sequence", 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			seen := make(map[uint64]bool)

			for i := 0; i < tc.numTests; i++ {
				magic := generateMagicNumber(&randomState)

				// Magic numbers should never be zero
				if magic == 0 {
					t.Error("Generated zero magic number")
				}

				// Check density (magic numbers should be relatively sparse)
				setBits := bitboard.Bitboard(magic).Count()
				if setBits > 48 { // Arbitrary threshold, adjust based on actual requirements
					t.Errorf("Magic number too dense: %d bits set", setBits)
				}

				// Check for duplicates
				if seen[magic] {
					t.Error("Generated duplicate magic number")
				}
				seen[magic] = true
			}
		})
	}
}

// TestSetOccupancy tests the occupancy bitboard generation
func TestSetOccupancy(t *testing.T) {
	testCases := []struct {
		name       string
		index      int
		bitsInMask int
		attackMask bitboard.Bitboard
		expected   int // Expected number of set bits in result
	}{
		{
			name:       "Empty occupancy",
			index:      0,
			bitsInMask: 4,
			attackMask: bitboard.Bitboard(0xF), // 1111
			expected:   0,
		},
		{
			name:       "Full occupancy",
			index:      0xF, // All bits set
			bitsInMask: 4,
			attackMask: bitboard.Bitboard(0xF), // 1111
			expected:   4,
		},
		{
			name:       "Partial occupancy",
			index:      5, // 0101
			bitsInMask: 4,
			attackMask: bitboard.Bitboard(0xF), // 1111
			expected:   2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := SetOccupancy(tc.index, tc.bitsInMask, tc.attackMask)
			setBits := result.Count()

			if setBits != tc.expected {
				t.Errorf("Expected %d set bits, got %d", tc.expected, setBits)
			}

			// Verify result is a subset of attack mask
			if result&tc.attackMask != result {
				t.Error("Generated occupancy contains bits outside attack mask")
			}
		})
	}
}

// TestMagicNumbersValidity tests the validity of precomputed magic numbers
func TestMagicNumbersValidity(t *testing.T) {
	testCases := []struct {
		name         string
		magics       [64]bitboard.Bitboard
		relevantBits [64]int
		piece        int
	}{
		{
			name:         "Bishop magic numbers",
			magics:       bishopMagicNumbers,
			relevantBits: bishopRelevantBits,
			piece:        Bishop,
		},
		{
			name:         "Rook magic numbers",
			magics:       rookMagicNumbers,
			relevantBits: rookRelevantBits,
			piece:        Rook,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for square := 0; square < 64; square++ {
				magic := tc.magics[square]

				// Verify magic number isn't zero
				if magic == 0 {
					t.Errorf("Square %d has zero magic number", square)
				}

				// Verify magic number can generate unique indices
				verifyMagicNumber(t, square, magic, tc.relevantBits[square], tc.piece)
			}
		})
	}
}

// Helper function to verify a magic number generates unique indices
func verifyMagicNumber(
	t *testing.T,
	square int,
	magic bitboard.Bitboard,
	relevantBits int,
	piece int,
) {
	var attackMask bitboard.Bitboard
	if piece == Bishop {
		attackMask = generateBishopPossibleBlockers(square)
	} else {
		attackMask = generateRookPossibleBlockers(square)
	}

	occupancyVariations := 1 << relevantBits
	used := make(map[int]bitboard.Bitboard)

	for i := 0; i < occupancyVariations; i++ {
		occupancy := SetOccupancy(i, relevantBits, attackMask)

		// Calculate magic index
		magicIndex := int((occupancy * magic) >> (64 - relevantBits))

		// Generate actual attacks for this occupancy
		var attacks bitboard.Bitboard
		if piece == Bishop {
			attacks = generateBishopAttacks(square, occupancy)
		} else {
			attacks = generateRookAttacks(square, occupancy)
		}

		// Check for collisions
		if prev, exists := used[magicIndex]; exists {
			if prev != attacks {
				t.Errorf("Magic collision detected for square %d", square)
			}
		}
		used[magicIndex] = attacks
	}
}

// BenchmarkMagicNumberGeneration benchmarks the magic number generation process
func BenchmarkMagicNumberGeneration(b *testing.B) {
	randomState := uint32(1804289383)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		generateMagicNumber(&randomState)
	}
}

// BenchmarkSetOccupancy benchmarks the occupancy generation process
func BenchmarkSetOccupancy(b *testing.B) {
	attackMask := bitboard.Bitboard(0xF0F0F0F0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SetOccupancy(i%16, 4, attackMask)
	}
}
