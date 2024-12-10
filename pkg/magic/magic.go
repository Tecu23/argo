package magic

import "github.com/Tecu23/argov2/pkg/bitboard"

// SetOccupancy generates an occupancy bitboard (a subset of the given attackMask)
// based on the binary pattern of 'index'.
//
// Parameters:
// - index: Used as a bit pattern to decide which squares from the attackMask are included.
// - bitsInMask: The number of set bits (ones) in attackMask. We will iterate over each set bit.
// - attackMask: A bitboard representing a set of squares (e.g., a ray of attacks).
//
// Process:
//  1. Convert the 'attackMask' into an ordered list of squares by repeatedly extracting the
//     least significant set bit using FirstOne().
//  2. For each extracted bit position (square), check the corresponding bit in 'index' (starting from 0).
//  3. If the bit in 'index' is set, include that square in the final occupancy bitboard.
//
// In other words, 'index' acts like a binary index into all subsets of the bits in attackMask.
// For example, if 'attackMask' has N set bits, then 'index' ranges from 0 to 2^N - 1, allowing
// you to select any subset of those bits.
func SetOccupancy(index, bitsInMask int, attackMask bitboard.Bitboard) bitboard.Bitboard {
	// Initialize an empty occupancy bitboard.
	occupancy := bitboard.Bitboard(0)

	// Iterate over the number of bits in the mask.
	// Each iteration removes one LSB (least significant bit) from attackMask.
	for count := 0; count < bitsInMask; count++ {
		// Extract the position of the least significant set bit in attackMask.
		square := attackMask.FirstOne()

		// Check if the corresponding bit in 'index' is set.
		// If so, add that square to the occupancy bitboard.
		if index&(1<<count) != 0 {
			occupancy |= (1 << square)
		}
	}

	// Return the generated occupancy subset.
	return occupancy
}
