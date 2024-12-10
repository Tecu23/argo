// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"fmt"

	"github.com/Tecu23/argov2/pkg/move"
	"github.com/Tecu23/argov2/pkg/util"
)

// Nodes is the number of positions reached during the perft
var Nodes int64

// perftDriver is a recursive function for the performance test (perft).
// Perft counts the number of positions reachable at a given depth.
// It's used to validate move generation correctness.
func perftDriver(b *Board, depth int) {
	// If depth=0, we've reached a leaf node, count it.
	if depth == 0 {
		Nodes++
		return
	}

	var mvlst move.Movelist
	b.GenerateMoves(&mvlst)

	for _, mv := range mvlst {
		copyB := b.CopyBoard()

		if !b.MakeMove(mv, AllMoves) {
			continue
		}

		perftDriver(b, depth-1)

		b.TakeBack(copyB)
	}
}

// PerftTest runs a perft test at a given depth, printing the move counts from the initial moves,
// and total nodes visited. Useful for debugging and verifying correctness of move generation.
func PerftTest(b *Board, depth int) {
	var whiteMoves move.Movelist

	totalMoves := int64(0)

	b.GenerateMoves(&whiteMoves)

	fmt.Printf("\n  Performance test\n\n")
	start := util.GetTimeInMiliseconds()

	for _, m := range whiteMoves {
		Nodes = 0

		copyB := b.CopyBoard()

		if !b.MakeMove(m, AllMoves) {
			continue
		}
		perftDriver(b, depth-1)

		// take back move
		b.TakeBack(copyB)

		fmt.Printf("%s: %d\n", m, Nodes)

		totalMoves += Nodes
	}
	// print results
	fmt.Printf("\n Nodes: %d Time: %d\n\n ", totalMoves, util.GetTimeInMiliseconds()-start)
}
