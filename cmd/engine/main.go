// Package main is the entry point of the program
package main

import (
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

func main() {
	initHelpers()

	b := bitboard.Bitboard(0)

	b.Set(constants.F4)
	b.Set(constants.C4)
	b.Set(constants.D2)
	att := attacks.GetQueenAttacks(constants.D4, b)

	att.PrintBitboard()
}

func initHelpers() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(constants.Bishop)
	attacks.InitSliderPiecesAttacks(constants.Rook)
}
