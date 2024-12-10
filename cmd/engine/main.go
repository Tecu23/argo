// Package main is the entry point of the program
package main

import (
	"github.com/Tecu23/argov2/pkg/attacks"
)

func main() {
	initHelpers()

	for sq := 0; sq < 64; sq++ {
		b := attacks.GenerateBishopAttacks(sq)
		b.PrintBitboard()
	}
}

func initHelpers() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
}
