// Package board contains the board representation and all board functions
package board

import (
	"fmt"

	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

type Board struct {
	Bitboards   [12]bitboard.Bitboard
	Occupancies [3]bitboard.Bitboard
	Side        color.Color
	EnPassant   int
	Castlings
}

func (b *Board) Reset() {
	b.Side = color.WHITE
	b.EnPassant = -1
	b.Castlings = 0

	for i := 0; i < 12; i++ {
		b.Bitboards[i] = 0
	}

	for i := 0; i < 3; i++ {
		b.Occupancies[i] = 0
	}
}

// PrintBoard should print the current position of the board
func (b Board) PrintBoard() {
	for rank := 7; rank >= 0; rank-- {
		for file := 0; file < 8; file++ {
			if file == 0 {
				fmt.Printf("%d  ", rank+1)
			}
			piece := -1

			// loop over all piece bitboards
			for bb := constants.WP; bb <= constants.BK; bb++ {
				if b.Bitboards[bb].Test(rank*8 + file) {
					piece = bb
				}
			}

			if piece == -1 {
				fmt.Printf(" %c", '.')
			} else {
				fmt.Printf(" %c", util.ASCIIPieces[piece])
			}

		}
		fmt.Println()
	}

	fmt.Printf("\n    a b c d e f g h\n\n")

	// print side to move
	if b.Side == color.WHITE {
		fmt.Printf(" Side:     %s\n", "white")
	} else {
		fmt.Printf(" Side:     %s\n", "black")
	}

	// print enpassant square
	fmt.Printf(" Enpassant:   %s\n", util.Sq2Fen[b.EnPassant])

	// print castling rights
	fmt.Printf(" Castling:  %s\n\n", b.Castlings.String())

	// fmt.Printf(" HashKey: 0x%X\n\n", b.Key)
}
