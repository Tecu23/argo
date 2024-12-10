package attacks

import "github.com/Tecu23/argov2/pkg/bitboard"

// func GetRookAttacks(sq int, occupancy Bitboard) Bitboard {
// 	// calculate magic index
// 	occupancy &= RookMasks[sq]
// 	occupancy *= RookMagicNumbers[sq]
// 	occupancy >>= 64 - RookRelevantBits[sq]
//
// 	return RookAttacks[sq][occupancy]
// }

// GenerateRookAttacks computes potential rook attacks from a given square assuming
// an *almost* empty board (no blockers), but this function stops at rank/file <= 6 and >= 1,
// which might be an incomplete logic or a partial version of a function. Normally, you would
// iterate fully to the edge of the board (0 to 7). Here, it seems to stop earlier.
//
// The board indexing convention: A8=0 at the top-left, and H1=63 at the bottom-right.
// Each rank increases going down (A7=8, A6=16, ...), and each file increases going right.
//
// Rooks move vertically and horizontally, so this function adds bits along the same
// rank and file from the starting square, but only within certain bounds.
func GenerateRookAttacks(square int) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// mask relevant bishop occupancy bits
	for r = tr + 1; r <= 6; r++ {
		attacks |= 1 << (r*8 + tf)
	}

	for f = tf + 1; f <= 6; f++ {
		attacks |= 1 << (tr*8 + f)
	}

	for r = tr - 1; r >= 1; r-- {
		attacks |= 1 << (r*8 + tf)
	}

	for f = tf - 1; f >= 1; f-- {
		attacks |= 1 << (tr*8 + f)
	}

	return attacks
}

// GenerateRookAttacksOnTheFly computes rook attacks for a given square and a given
// occupancy bitboard "block" that represents which squares are occupied. The rook
// attack rays stop when they encounter a block.
//
// Unlike GenerateRookAttacks, this function goes fully to the edges of the board (0 to 7).
// When a blocker is hit on a particular square, that direction is stopped.
//
// Parameters:
// - square: the position of the rook (0-based index on an 8x8 board with A8=0).
// - block: bitboard of occupied squares that can block the rook.
func GenerateRookAttacksOnTheFly(square int, block bitboard.Bitboard) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// mask relevant bishop occupancy bits
	for r = tr + 1; r <= 7; r++ {
		attacks |= 1 << (r*8 + tf)
		if (1<<(r*8+tf))&block != 0 {
			break
		}
	}

	for f = tf + 1; f <= 7; f++ {
		attacks |= 1 << (tr*8 + f)
		if (1<<(tr*8+f))&block != 0 {
			break
		}
	}

	for r = tr - 1; r >= 0; r-- {
		attacks |= 1 << (r*8 + tf)
		if (1<<(r*8+tf))&block != 0 {
			break
		}
	}

	for f = tf - 1; f >= 0; f-- {
		attacks |= 1 << (tr*8 + f)
		if (1<<(tr*8+f))&block != 0 {
			break
		}
	}

	return attacks
}
