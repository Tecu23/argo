// Package util contains utility functions for piece conversion, timing, and square indexing.
package util

import (
	"time"

	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// PcFen contains pieces in Fen format for validation
const (
	PcFen    = "PpNnBbRrQqKk     "
	PieceFen = "PNBRQK     "
)

// ASCIIPieces maps piece indices to ASCII representation for printing.
const (
	ASCIIPieces = "PNBRQKpnbrqk"
)

// Fen2pc converts a Fen piece character to an internal piece index.
func Fen2pc(c string) int {
	for p, x := range ASCIIPieces {
		if string(x) == c {
			return p
		}
	}
	return Empty
}

// PcColor returns WHITE if pc < 6, otherwise BLACK, identifying piece color by its index.
func PcColor(pc int) color.Color {
	if pc < 6 {
		return color.WHITE
	}
	return color.BLACK
}

// Fen2Sq and Sq2Fen convert between algebraic notation (e.g., "e4") and internal square indices.
var (
	Fen2Sq = make(map[string]int)
	Sq2Fen = make(map[int]string)
)

// InitFen2Sq initializes the Fen2Sq and Sq2Fen maps for all squares. This allows easy
// conversion between numeric indices and coordinate notation.
func InitFen2Sq() {
	Fen2Sq["a1"] = A1
	Fen2Sq["a2"] = A2
	Fen2Sq["a3"] = A3
	Fen2Sq["a4"] = A4
	Fen2Sq["a5"] = A5
	Fen2Sq["a6"] = A6
	Fen2Sq["a7"] = A7
	Fen2Sq["a8"] = A8

	Fen2Sq["b1"] = B1
	Fen2Sq["b2"] = B2
	Fen2Sq["b3"] = B3
	Fen2Sq["b4"] = B4
	Fen2Sq["b5"] = B5
	Fen2Sq["b6"] = B6
	Fen2Sq["b7"] = B7
	Fen2Sq["b8"] = B8

	Fen2Sq["c1"] = C1
	Fen2Sq["c2"] = C2
	Fen2Sq["c3"] = C3
	Fen2Sq["c4"] = C4
	Fen2Sq["c5"] = C5
	Fen2Sq["c6"] = C6
	Fen2Sq["c7"] = C7
	Fen2Sq["c8"] = C8

	Fen2Sq["d1"] = D1
	Fen2Sq["d2"] = D2
	Fen2Sq["d3"] = D3
	Fen2Sq["d4"] = D4
	Fen2Sq["d5"] = D5
	Fen2Sq["d6"] = D6
	Fen2Sq["d7"] = D7
	Fen2Sq["d8"] = D8

	Fen2Sq["e1"] = E1
	Fen2Sq["e2"] = E2
	Fen2Sq["e3"] = E3
	Fen2Sq["e4"] = E4
	Fen2Sq["e5"] = E5
	Fen2Sq["e6"] = E6
	Fen2Sq["e7"] = E7
	Fen2Sq["e8"] = E8

	Fen2Sq["f1"] = F1
	Fen2Sq["f2"] = F2
	Fen2Sq["f3"] = F3
	Fen2Sq["f4"] = F4
	Fen2Sq["f5"] = F5
	Fen2Sq["f6"] = F6
	Fen2Sq["f7"] = F7
	Fen2Sq["f8"] = F8

	Fen2Sq["g1"] = G1
	Fen2Sq["g2"] = G2
	Fen2Sq["g3"] = G3
	Fen2Sq["g4"] = G4
	Fen2Sq["g5"] = G5
	Fen2Sq["g6"] = G6
	Fen2Sq["g7"] = G7
	Fen2Sq["g8"] = G8

	Fen2Sq["h1"] = H1
	Fen2Sq["h2"] = H2
	Fen2Sq["h3"] = H3
	Fen2Sq["h4"] = H4
	Fen2Sq["h5"] = H5
	Fen2Sq["h6"] = H6
	Fen2Sq["h7"] = H7
	Fen2Sq["h8"] = H8

	// -------------- Sq2Fen
	Sq2Fen[A1] = "a1"
	Sq2Fen[A2] = "a2"
	Sq2Fen[A3] = "a3"
	Sq2Fen[A4] = "a4"
	Sq2Fen[A5] = "a5"
	Sq2Fen[A6] = "a6"
	Sq2Fen[A7] = "a7"
	Sq2Fen[A8] = "a8"

	Sq2Fen[B1] = "b1"
	Sq2Fen[B2] = "b2"
	Sq2Fen[B3] = "b3"
	Sq2Fen[B4] = "b4"
	Sq2Fen[B5] = "b5"
	Sq2Fen[B6] = "b6"
	Sq2Fen[B7] = "b7"
	Sq2Fen[B8] = "b8"

	Sq2Fen[C1] = "c1"
	Sq2Fen[C2] = "c2"
	Sq2Fen[C3] = "c3"
	Sq2Fen[C4] = "c4"
	Sq2Fen[C5] = "c5"
	Sq2Fen[C6] = "c6"
	Sq2Fen[C7] = "c7"
	Sq2Fen[C8] = "c8"

	Sq2Fen[D1] = "d1"
	Sq2Fen[D2] = "d2"
	Sq2Fen[D3] = "d3"
	Sq2Fen[D4] = "d4"
	Sq2Fen[D5] = "d5"
	Sq2Fen[D6] = "d6"
	Sq2Fen[D7] = "d7"
	Sq2Fen[D8] = "d8"

	Sq2Fen[E1] = "e1"
	Sq2Fen[E2] = "e2"
	Sq2Fen[E3] = "e3"
	Sq2Fen[E4] = "e4"
	Sq2Fen[E5] = "e5"
	Sq2Fen[E6] = "e6"
	Sq2Fen[E7] = "e7"
	Sq2Fen[E8] = "e8"

	Sq2Fen[F1] = "f1"
	Sq2Fen[F2] = "f2"
	Sq2Fen[F3] = "f3"
	Sq2Fen[F4] = "f4"
	Sq2Fen[F5] = "f5"
	Sq2Fen[F6] = "f6"
	Sq2Fen[F7] = "f7"
	Sq2Fen[F8] = "f8"

	Sq2Fen[G1] = "g1"
	Sq2Fen[G2] = "g2"
	Sq2Fen[G3] = "g3"
	Sq2Fen[G4] = "g4"
	Sq2Fen[G5] = "g5"
	Sq2Fen[G6] = "g6"
	Sq2Fen[G7] = "g7"
	Sq2Fen[G8] = "g8"

	Sq2Fen[H1] = "h1"
	Sq2Fen[H2] = "h2"
	Sq2Fen[H3] = "h3"
	Sq2Fen[H4] = "h4"
	Sq2Fen[H5] = "h5"
	Sq2Fen[H6] = "h6"
	Sq2Fen[H7] = "h7"
	Sq2Fen[H8] = "h8"
	Sq2Fen[-1] = "-"
}

// GetTimeInMiliseconds returns current time in milliseconds since Unix epoch.
// Used for timing operations like perft tests.
func GetTimeInMiliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// OppositeColorPiece returns the opposite color piece.
func OppositeColorPiece(piece int) int {
	// Handle the empty case
	if piece == Empty {
		return Empty
	}

	// Determine if the piece is white (0-5) or black (6-11)
	if piece >= WP && piece <= WK {
		return piece + 6 // Convert white to black
	} else if piece >= BP && piece <= BK {
		return piece - 6 // Convert black to white
	}

	// Return Empty for invalid input
	return Empty
}

func GetPieceType(piece int) int {
	if piece == Empty {
		return -1
	}

	if piece == WK || piece == BK {
		return Pawn
	} else if piece == WN || piece == BN {
		return Knight
	} else if piece == WB || piece == BB {
		return Bishop
	} else if piece == WR || piece == BR {
		return Rook
	} else if piece == WQ || piece == BQ {
		return Queen
	}

	return King
}
