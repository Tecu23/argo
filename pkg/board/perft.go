// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"fmt"

	"github.com/Tecu23/argov2/pkg/util"
)

// PerftDriver is a recursive function for the performance test (perft).
// Perft counts the number of positions reachable at a given depth.
// It's used to validate move generation correctness.
func PerftDriver(b *Board, depth int) int64 {
	// If depth=0, we've reached a leaf node, count it.
	if depth == 0 {
		return 1
	}

	nds := int64(0)
	mvlst := b.GenerateMoves()

	for _, mv := range mvlst {
		copyB := b.CopyBoard()

		if !b.MakeMove(mv, AllMoves) {
			continue
		}

		nds += PerftDriver(b, depth-1)

		b.TakeBack(copyB)
	}
	return nds
}

// PerftTest runs a perft test at a given depth, printing the move counts from the initial moves,
// and total nodes visited. Useful for debugging and verifying correctness of move generation.
func PerftTest(b *Board, depth int) int64 {
	totalMoves := int64(0)
	whiteMoves := b.GenerateMoves()
	fmt.Printf("\n  Performance test\n\n")
	start := util.GetTimeInMiliseconds()

	for _, m := range whiteMoves {
		moveNodes := int64(0)
		copyB := b.CopyBoard()

		if !b.MakeMove(m, AllMoves) {
			continue
		}
		moveNodes = PerftDriver(b, depth-1)

		// take back move
		b.TakeBack(copyB)

		fmt.Printf("%s: %d\n", m, moveNodes)
		totalMoves += moveNodes
	}
	// print results
	fmt.Printf(
		"\n Nodes: %d Time: %d, with: %.0f nodes/s\n\n ",
		totalMoves,
		util.GetTimeInMiliseconds()-start,
		float64(totalMoves)/(float64(util.GetTimeInMiliseconds()-start)/1000.0),
	)

	return totalMoves
}
