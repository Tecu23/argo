// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"strings"

	. "github.com/Tecu23/argov2/pkg/constants"
)

// CastlingRights is an array that helps update castling rights
// when a certain square is moved from or to. Its values represent bit masks
// that modify the current castling availability.
var CastlingRights = []uint{
	7, 15, 15, 15, 3, 15, 15, 11,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15,
	13, 15, 15, 15, 12, 15, 15, 14,
}

/*
	castling  move    in       in
	   right update   binary   decimal

king & rooks didn't move:       1111 & 1111  =  1111    15

	       white king moved:       1111 & 1100  =  1100    12
	white king's rook moved:       1111 & 1110  =  1110    14
   white queen's rook moved:       1111 & 1101  =  1101    13

	       black king moved:       1111 & 0011  =  0011    3
	black king's rook moved:       1111 & 1011  =  1011    11
   black queen's rook moved:       1111 & 0111  =  0111    7
*/

// Castlings represents the current castling rights as a bitfield.
type Castlings uint

// ParseCastlings parses the FEN castling substring (e.g., "KQkq") and returns
// the corresponding castling rights bitfield.
func ParseCastlings(fenCastl string) Castlings {
	c := uint(0)

	if fenCastl == "-" {
		return Castlings(0)
	}

	if strings.Contains(fenCastl, "K") {
		c |= ShortW
	}
	if strings.Contains(fenCastl, "Q") {
		c |= LongW
	}
	if strings.Contains(fenCastl, "k") {
		c |= ShortB
	}
	if strings.Contains(fenCastl, "q") {
		c |= LongB
	}

	return Castlings(c)
}

// String returns a string representation of the current castling rights (e.g., "KQkq" or "-").
func (c Castlings) String() string {
	flags := ""
	if uint(c)&ShortW != 0 {
		flags = "K"
	}
	if uint(c)&LongW != 0 {
		flags += "Q"
	}
	if uint(c)&ShortB != 0 {
		flags += "k"
	}
	if uint(c)&LongB != 0 {
		flags += "q"
	}
	if flags == "" {
		flags = "-"
	}
	return flags
}
