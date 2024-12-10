package board

import (
	"strings"

	"github.com/Tecu23/argov2/pkg/constants"
)

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

// Castlings represents the castling possibilities for a given position
type Castlings uint

func ParseCastlings(fenCastl string) Castlings {
	c := uint(0)

	if fenCastl == "-" {
		return Castlings(0)
	}

	if strings.Contains(fenCastl, "K") {
		c |= constants.ShortW
	}
	if strings.Contains(fenCastl, "Q") {
		c |= constants.LongW
	}
	if strings.Contains(fenCastl, "k") {
		c |= constants.ShortB
	}
	if strings.Contains(fenCastl, "q") {
		c |= constants.LongB
	}

	return Castlings(c)
}

// String returns the string represantion of the castling
func (c Castlings) String() string {
	flags := ""
	if uint(c)&constants.ShortW != 0 {
		flags = "K"
	}
	if uint(c)&constants.LongW != 0 {
		flags += "Q"
	}
	if uint(c)&constants.ShortB != 0 {
		flags += "k"
	}
	if uint(c)&constants.LongB != 0 {
		flags += "q"
	}
	if flags == "" {
		flags = "-"
	}
	return flags
}
