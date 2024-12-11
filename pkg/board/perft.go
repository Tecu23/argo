// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"fmt"

	"github.com/Tecu23/argov2/pkg/move"
	"github.com/Tecu23/argov2/pkg/util"
)

// PerftDriver is a recursive function for the performance test (perft).
// Perft counts the number of positions reachable at a given depth.
// It's used to validate move generation correctness.
func PerftDriver(b *Board, depth int, nodes int64) int64 {
	// If depth=0, we've reached a leaf node, count it.
	if depth == 0 {
		nodes++
		return 1
	}

	var mvlst move.Movelist
	b.GenerateMoves(&mvlst)

	nds := int64(0)

	for _, mv := range mvlst {
		copyB := b.CopyBoard()

		if !b.MakeMove(mv, AllMoves) {
			continue
		}

		nds += PerftDriver(b, depth-1, nodes)

		b.TakeBack(copyB)
	}
	nodes += nds
	return nodes
}

// PerftTest runs a perft test at a given depth, printing the move counts from the initial moves,
// and total nodes visited. Useful for debugging and verifying correctness of move generation.
func PerftTest(b *Board, depth int, nodes int64) {
	var whiteMoves move.Movelist

	totalMoves := int64(0)

	b.GenerateMoves(&whiteMoves)

	fmt.Printf("\n  Performance test\n\n")
	start := util.GetTimeInMiliseconds()

	for _, m := range whiteMoves {
		tmp_nodes := int64(0)

		copyB := b.CopyBoard()

		if !b.MakeMove(m, AllMoves) {
			continue
		}
		PerftDriver(b, depth-1, tmp_nodes)

		// take back move
		b.TakeBack(copyB)

		fmt.Printf("%s: %d\n", m, nodes)

		totalMoves += tmp_nodes
	}
	// print results
	fmt.Printf("\n Nodes: %d Time: %d\n\n ", totalMoves, util.GetTimeInMiliseconds()-start)
}
