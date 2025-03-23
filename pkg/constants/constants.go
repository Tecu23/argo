// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package constants contains the shared variabled between the different packages
package constants

import "github.com/Tecu23/argov2/pkg/bitboard"

// StartPosition holds the starting position of a chess game
const StartPosition = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 "

// Piece constants are an enumeration for piece types.
// They help identify which piece's attack mask or magic to use.
const (
	Pawn       int = iota // 0
	Knight                // 1
	Bishop                // 2
	Rook                  // 3
	Queen                 // 4
	King                  // 5
	PieceTypes            // 6
)

// Constants for piece encoding, castling, etc. Provided here for completeness.
const (
	WP = iota // 0
	WN        // 1
	WB        // 2
	WR        // 3
	WQ        // 4
	WK        // 5
	// _             // skip (6)
	// _             // skip (7)
	BP          // 8
	BN          // 9
	BB          // 10
	BR          // 11
	BQ          // 12
	BK          // 13
	Pieces      // 14
	Empty  = 15 // 15
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

var FileMasks = [8]bitboard.Bitboard{
	0x0101010101010101, // File A (squares A1, A2, A3, A4, A5, A6, A7, A8)
	0x0202020202020202, // File B (squares B1, B2, B3, B4, B5, B6, B7, B8)
	0x0404040404040404, // File C (squares C1, C2, C3, C4, C5, C6, C7, C8)
	0x0808080808080808, // File D (squares D1, D2, D3, D4, D5, D6, D7, D8)
	0x1010101010101010, // File E (squares E1, E2, E3, E4, E5, E6, E7, E8)
	0x2020202020202020, // File F (squares F1, F2, F3, F4, F5, F6, F7, F8)
	0x4040404040404040, // File G (squares G1, G2, G3, G4, G5, G6, G7, G8)
	0x8080808080808080, // File H (squares H1, H2, H3, H4, H5, H6, H7, H8)
}

var RankMasks = [8]bitboard.Bitboard{
	0x00000000000000FF, // Rank 1 (squares A1, B1, C1, D1, E1, F1, G1, H1)
	0x000000000000FF00, // Rank 2 (squares A2, B2, C2, D2, E2, F2, G2, H2)
	0x0000000000FF0000, // Rank 3 (squares A3, B3, C3, D3, E3, F3, G3, H3)
	0x00000000FF000000, // Rank 4 (squares A4, B4, C4, D4, E4, F4, G4, H4)
	0x000000FF00000000, // Rank 5 (squares A5, B5, C5, D5, E5, F5, G5, H5)
	0x0000FF0000000000, // Rank 6 (squares A6, B6, C6, D6, E6, F6, G6, H6)
	0x00FF000000000000, // Rank 7 (squares A7, B7, C7, D7, E7, F7, G7, H7)
	0xFF00000000000000, // Rank 8 (squares A8, B8, C8, D8, E8, F8, G8, H8)
}

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
	ShortW = uint(0x1) // White short castlings right
	LongW  = uint(0x2) // White long castlings right
	ShortB = uint(0x4) // Black short castlings right
	LongB  = uint(0x8) // Black long castlings right
)
