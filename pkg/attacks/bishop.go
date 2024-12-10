package attacks

import "github.com/Tecu23/argov2/pkg/bitboard"

// GenerateBishopAttacks computes the potential attack squares of a bishop placed on a given square,
// assuming an *empty board* (no blocking pieces).
//
// The board is indexed from 0 to 63, where A8 = 0, B8 = 1, ..., H8 = 7, A7 = 8, ..., H1 = 63.
// This means row (rank) = square / 8 and file = square % 8.
//
// This function stops generating attacks before the last rank/file (uses <= 6 checks).
// It seems intended to produce a subset or is an older version of code. Normally, you'd iterate
// until the edge of the board (0 <= r,f <= 7). Here, it stops at rank/file 6 for some reason,
// which may be intentional or a leftover from partial logic.
func GenerateBishopAttacks(square int) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// Move diagonally (up-right)
	for r, f = tr+1, tf+1; r <= 6 && f <= 6; r, f = r+1, f+1 {
		attacks |= 1 << (r*8 + f)
	}

	// Move diagonally (up-left)
	for r, f = tr-1, tf+1; r >= 1 && f <= 6; r, f = r-1, f+1 {
		attacks |= 1 << (r*8 + f)
	}

	// Move diagonally (down-right)
	for r, f = tr+1, tf-1; r <= 6 && f >= 1; r, f = r+1, f-1 {
		attacks |= 1 << (r*8 + f)
	}

	// Move diagonally (down-left)
	for r, f = tr-1, tf-1; r >= 1 && f >= 1; r, f = r-1, f-1 {
		attacks |= 1 << (r*8 + f)
	}

	return attacks
}

// GenerateBishopAttacksOnTheFly computes the potential attack squares of a bishop placed on a given square,
// considering the given bitboard "block" which represents occupied squares. If a blocked square is encountered,
// the bishop's attack ray stops in that direction.
//
// This function uses full 0 to 7 range checks for rank/file, which aligns with the entire board dimension.
//
// Parameters:
// - square: The position of the bishop on a 0-63 indexed board.
// - block:  A bitboard representing occupied squares that might block the bishop.
//
// The bishop attacks diagonally in all four directions. As soon as a blocked square is found, that ray stops.
func GenerateBishopAttacksOnTheFly(square int, block bitboard.Bitboard) bitboard.Bitboard {
	attacks := bitboard.Bitboard(0)

	// init rank & files
	r, f := 0, 0

	// init target rank & files
	tr := square / 8
	tf := square % 8

	// up-right (increasing rank and file)
	for r, f = tr+1, tf+1; r <= 7 && f <= 7; r, f = r+1, f+1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	// up-left (decreasing rank, increasing file)
	for r, f = tr-1, tf+1; r >= 0 && f <= 7; r, f = r-1, f+1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	// down-right (increasing rank, decreasing file)
	for r, f = tr+1, tf-1; r <= 7 && f >= 0; r, f = r+1, f-1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	// down-left (decreasing rank and file)
	for r, f = tr-1, tf-1; r >= 0 && f >= 0; r, f = r-1, f-1 {
		attacks |= 1 << (r*8 + f)
		if (1<<(r*8+f))&block != 0 {
			break
		}
	}

	return attacks
}
