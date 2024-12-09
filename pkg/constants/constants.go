// Package constants contains the shared variabled between the different packages
package constants

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
