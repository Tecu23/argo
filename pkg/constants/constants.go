// Package constants contains the shared variabled between the different packages
package constants

import "github.com/Tecu23/argov2/pkg/bitboard"

// Piece constants are an enumeration for piece types.
// They help identify which piece's attack mask or magic to use.
const (
	Pawn int = iota
	Knight
	Bishop
	Rook
	Queen
	King
)

const (
	WP = iota
	WN
	WB
	WR
	WQ
	WK
	BP
	BN
	BB
	BR
	BQ
	BK
	Empty = 15
)

// The constants below define row (rank) and file bitboards.
// Each bitboard highlights all squares on a particular rank or file.
//
// For example, Row1 is 0x00000000000000FF, which in binary sets bits
// corresponding to the bottom rank (from White's perspective).
// Similarly, FileA is 0x0101010101010101, which sets all squares in the 'a' file.
const (
	Row1 = bitboard.Bitboard(0x00000000000000FF) // Represents rank 1 (a1-h1)
	Row2 = bitboard.Bitboard(0x000000000000FF00) // Rank 2 (a2-h2)
	Row3 = bitboard.Bitboard(0x0000000000FF0000) // Rank 3 (a3-h3)
	Row4 = bitboard.Bitboard(0x00000000FF000000) // Rank 4 (a4-h4)
	Row5 = bitboard.Bitboard(0x000000FF00000000) // Rank 5 (a5-h5)
	Row6 = bitboard.Bitboard(0x0000FF0000000000) // Rank 6 (a6-h6)
	Row7 = bitboard.Bitboard(0x00FF000000000000) // Rank 7 (a7-h7)
	Row8 = bitboard.Bitboard(0xFF00000000000000) // Rank 8 (a8-h8)

	FileA = bitboard.Bitboard(0x0101010101010101) // File a (a1, a2, ..., a8)
	FileB = bitboard.Bitboard(0x0202020202020202) // File b
	FileC = bitboard.Bitboard(0x0404040404040404) // File c
	FileD = bitboard.Bitboard(0x0808080808080808) // File d
	FileE = bitboard.Bitboard(0x1010101010101010) // File e
	FileF = bitboard.Bitboard(0x2020202020202020) // File f
	FileG = bitboard.Bitboard(0x4040404040404040) // File g
	FileH = bitboard.Bitboard(0x8080808080808080) // File h
)

// The constants below define directional offsets used for moving
// pieces or generating attacks on the board.
// Each direction represents a shift in terms of the bit index:
//
// N  (North) = +8: move up one rank (higher-indexed rank)
// S  (South) = -8: move down one rank (lower-indexed rank)
// E  (East)  = +1: move one file to the right
// W  (West)  = -1: move one file to the left
//
// NW = +7 and NE = +9 represent diagonal moves. Similarly, SW and SE are defined
// relative to these directions. Shifting a bitboard by these constants simulates
// piece movement or attack generation in the given direction.
const (
	E  = +1  // Move one square to the right (east)
	W  = -1  // Move one square to the left (west)
	N  = 8   // Move one rank up (north)
	S  = -8  // Move one rank down (south)
	NW = +7  // Move one square diagonally northwest
	NE = +9  // Move one square diagonally northeast
	SW = -NE // Move diagonally southwest (negative of northeast)
	SE = -NW // Move diagonally southeast (negative of northwest)
)

// The constants below represent indices for each square on a chessboard.
// They follow a top-left (A8) to bottom-right (H1) indexing, starting with A8 as 0.
// The first rank (from White's perspective) is at the bottom (A1 to H1) and the eighth rank at the top (A8 to H8).
// Each subsequent constant corresponds to a square moving left-to-right, top-to-bottom.
//
// For reference:
//
//	Ranks are numbered 1 to 8 from White's side (1 is closest to White, 8 is closest to Black).
//	Files are labeled a to h from left (White's left) to right.
//
// The indexing scheme here starts from A8 (top-left corner) and moves horizontally first, then down to the next rank.
// Example: A8=0, B8=1, ..., H8=7, A7=8, B7=9, ... , H1=63.
//
// This layout corresponds to a bitboard representation where bit 0 could be A8 and bit 63 is H1.
const (
	A8 = iota
	B8
	C8
	D8
	E8
	F8
	G8
	H8

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A1
	B1
	C1
	D1
	E1
	F1
	G1
	H1
)

/*
*           CASTLINGS CONSTANTS
 */
/*
- 0001 -- 1 -> white king can castle to the king side
- 0010 -- 2 -> white king can castle to the queen side
- 0100 -- 4 -> black king can castle to the king side
- 1000 -- 8 -> black king can castle to the queen side

ex.
	    1111 -> both side can castle in both directons

		1001 -> black king => queen side
		     -> white king => king side
*/
const (
	ShortW = uint(0x1)
	LongW  = uint(0x2)
	ShortB = uint(0x4)
	LongB  = uint(0x8)
)
