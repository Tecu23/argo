// Package main is the entry point of the program
package main

import (
	"bufio"
	"os"

	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/magic"
)

func main() {
	initHelpers()

	mask := attacks.GenerateRookAttacks(constants.A1)

	for i := 0; i < 4096; i++ {
		occ := magic.SetOccupancy(i, mask.Count(), mask)
		occ.PrintBitboard()
		input := bufio.NewScanner(os.Stdin)
		input.Scan()

	}
}

func initHelpers() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
}
