// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	util.InitFen2Sq()
}

func TestPerftInitialPosition(t *testing.T) {
	var b Board
	// Standard initial position
	b.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	// TODO: Update perft functions to return the nodes and properly tests the data
}
