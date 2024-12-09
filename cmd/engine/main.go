// Package main is the entry point of the program
package main

import (
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

func main() {
	b := bitboard.Bitboard(0)

	b.Set(constants.E2)

	b.PrintBitboard()
}
