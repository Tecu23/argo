// Package board contains the board representation and all board functions
package board

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
)

type Board struct {
	BitBoards   [12]bitboard.Bitboard
	Occupancies [3]bitboard.Bitboard
	Side        color.Color
	EnPassant   int
	Castlings
}
