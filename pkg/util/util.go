// Package util contains utility functions for piece conversion, timing, and square indexing.
package util

import (
	"time"

	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
)

// PcFen contains pieces in Fen format for validation
const (
	PcFen = "PpNnBbRrQqKk     "
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
	return constants.Empty
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
	Fen2Sq["a1"] = constants.A1
	Fen2Sq["a2"] = constants.A2
	Fen2Sq["a3"] = constants.A3
	Fen2Sq["a4"] = constants.A4
	Fen2Sq["a5"] = constants.A5
	Fen2Sq["a6"] = constants.A6
	Fen2Sq["a7"] = constants.A7
	Fen2Sq["a8"] = constants.A8

	Fen2Sq["b1"] = constants.B1
	Fen2Sq["b2"] = constants.B2
	Fen2Sq["b3"] = constants.B3
	Fen2Sq["b4"] = constants.B4
	Fen2Sq["b5"] = constants.B5
	Fen2Sq["b6"] = constants.B6
	Fen2Sq["b7"] = constants.B7
	Fen2Sq["b8"] = constants.B8

	Fen2Sq["c1"] = constants.C1
	Fen2Sq["c2"] = constants.C2
	Fen2Sq["c3"] = constants.C3
	Fen2Sq["c4"] = constants.C4
	Fen2Sq["c5"] = constants.C5
	Fen2Sq["c6"] = constants.C6
	Fen2Sq["c7"] = constants.C7
	Fen2Sq["c8"] = constants.C8

	Fen2Sq["d1"] = constants.D1
	Fen2Sq["d2"] = constants.D2
	Fen2Sq["d3"] = constants.D3
	Fen2Sq["d4"] = constants.D4
	Fen2Sq["d5"] = constants.D5
	Fen2Sq["d6"] = constants.D6
	Fen2Sq["d7"] = constants.D7
	Fen2Sq["d8"] = constants.D8

	Fen2Sq["e1"] = constants.E1
	Fen2Sq["e2"] = constants.E2
	Fen2Sq["e3"] = constants.E3
	Fen2Sq["e4"] = constants.E4
	Fen2Sq["e5"] = constants.E5
	Fen2Sq["e6"] = constants.E6
	Fen2Sq["e7"] = constants.E7
	Fen2Sq["e8"] = constants.E8

	Fen2Sq["f1"] = constants.F1
	Fen2Sq["f2"] = constants.F2
	Fen2Sq["f3"] = constants.F3
	Fen2Sq["f4"] = constants.F4
	Fen2Sq["f5"] = constants.F5
	Fen2Sq["f6"] = constants.F6
	Fen2Sq["f7"] = constants.F7
	Fen2Sq["f8"] = constants.F8

	Fen2Sq["g1"] = constants.G1
	Fen2Sq["g2"] = constants.G2
	Fen2Sq["g3"] = constants.G3
	Fen2Sq["g4"] = constants.G4
	Fen2Sq["g5"] = constants.G5
	Fen2Sq["g6"] = constants.G6
	Fen2Sq["g7"] = constants.G7
	Fen2Sq["g8"] = constants.G8

	Fen2Sq["h1"] = constants.H1
	Fen2Sq["h2"] = constants.H2
	Fen2Sq["h3"] = constants.H3
	Fen2Sq["h4"] = constants.H4
	Fen2Sq["h5"] = constants.H5
	Fen2Sq["h6"] = constants.H6
	Fen2Sq["h7"] = constants.H7
	Fen2Sq["h8"] = constants.H8

	// -------------- Sq2Fen
	Sq2Fen[constants.A1] = "a1"
	Sq2Fen[constants.A2] = "a2"
	Sq2Fen[constants.A3] = "a3"
	Sq2Fen[constants.A4] = "a4"
	Sq2Fen[constants.A5] = "a5"
	Sq2Fen[constants.A6] = "a6"
	Sq2Fen[constants.A7] = "a7"
	Sq2Fen[constants.A8] = "a8"

	Sq2Fen[constants.B1] = "b1"
	Sq2Fen[constants.B2] = "b2"
	Sq2Fen[constants.B3] = "b3"
	Sq2Fen[constants.B4] = "b4"
	Sq2Fen[constants.B5] = "b5"
	Sq2Fen[constants.B6] = "b6"
	Sq2Fen[constants.B7] = "b7"
	Sq2Fen[constants.B8] = "b8"

	Sq2Fen[constants.C1] = "c1"
	Sq2Fen[constants.C2] = "c2"
	Sq2Fen[constants.C3] = "c3"
	Sq2Fen[constants.C4] = "c4"
	Sq2Fen[constants.C5] = "c5"
	Sq2Fen[constants.C6] = "c6"
	Sq2Fen[constants.C7] = "c7"
	Sq2Fen[constants.C8] = "c8"

	Sq2Fen[constants.D1] = "d1"
	Sq2Fen[constants.D2] = "d2"
	Sq2Fen[constants.D3] = "d3"
	Sq2Fen[constants.D4] = "d4"
	Sq2Fen[constants.D5] = "d5"
	Sq2Fen[constants.D6] = "d6"
	Sq2Fen[constants.D7] = "d7"
	Sq2Fen[constants.D8] = "d8"

	Sq2Fen[constants.E1] = "e1"
	Sq2Fen[constants.E2] = "e2"
	Sq2Fen[constants.E3] = "e3"
	Sq2Fen[constants.E4] = "e4"
	Sq2Fen[constants.E5] = "e5"
	Sq2Fen[constants.E6] = "e6"
	Sq2Fen[constants.E7] = "e7"
	Sq2Fen[constants.E8] = "e8"

	Sq2Fen[constants.F1] = "f1"
	Sq2Fen[constants.F2] = "f2"
	Sq2Fen[constants.F3] = "f3"
	Sq2Fen[constants.F4] = "f4"
	Sq2Fen[constants.F5] = "f5"
	Sq2Fen[constants.F6] = "f6"
	Sq2Fen[constants.F7] = "f7"
	Sq2Fen[constants.F8] = "f8"

	Sq2Fen[constants.G1] = "g1"
	Sq2Fen[constants.G2] = "g2"
	Sq2Fen[constants.G3] = "g3"
	Sq2Fen[constants.G4] = "g4"
	Sq2Fen[constants.G5] = "g5"
	Sq2Fen[constants.G6] = "g6"
	Sq2Fen[constants.G7] = "g7"
	Sq2Fen[constants.G8] = "g8"

	Sq2Fen[constants.H1] = "h1"
	Sq2Fen[constants.H2] = "h2"
	Sq2Fen[constants.H3] = "h3"
	Sq2Fen[constants.H4] = "h4"
	Sq2Fen[constants.H5] = "h5"
	Sq2Fen[constants.H6] = "h6"
	Sq2Fen[constants.H7] = "h7"
	Sq2Fen[constants.H8] = "h8"
	Sq2Fen[-1] = "-"
}

// GetTimeInMiliseconds returns current time in milliseconds since Unix epoch.
// Used for timing operations like perft tests.
func GetTimeInMiliseconds() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
