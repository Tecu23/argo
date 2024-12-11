// Package main is the entry point of the program
package main

import (
	"fmt"
	"time"

	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/board"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func main() {
	initHelpers()

	b, _ := board.ParseFEN(constants.StartPosition)

	start := time.Now()
	nodes := board.PerftDriver(&b, 7, 0)
	duration := time.Since(start)
	fmt.Println(nodes, duration)
}

func initHelpers() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(constants.Bishop)
	attacks.InitSliderPiecesAttacks(constants.Rook)

	util.InitFen2Sq()
}
