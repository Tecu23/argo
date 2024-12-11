// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	util.InitFen2Sq() // Make sure square mappings are initialized
}

func TestParseFEN(t *testing.T) {
	var b Board
	b.ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	if b.Side != color.WHITE {
		t.Errorf("Expected side to move: WHITE, got %v", b.Side)
	}
	if b.EnPassant != -1 {
		t.Errorf("Expected no en passant, got %d", b.EnPassant)
	}
	if b.Rule50 != 0 {
		t.Errorf("Expected rule50=0, got %d", b.Rule50)
	}

	// Check a piece square:
	// A1 should be a White Rook (WR=3)
	if !b.Bitboards[constants.WR].Test(constants.A1) {
		t.Error("Expected White Rook at A1")
	}
}

func TestIsSquareAttacked(t *testing.T) {
	var b Board
	b.ParseFEN("rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR w KQkq - 0 1")

	// TODO: Implement this
}

func TestMakeMove(t *testing.T) {
	var b Board
	b.ParseFEN("rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 1")
	// Make move e4-e5:
	m := move.EncodeMove(constants.E4, constants.E5, constants.WP, 0, 0, 0, 0, 0)
	if !b.MakeMove(m, AllMoves) {
		t.Errorf("Expected move e4-e5 to be made successfully")
	}
	if !b.Bitboards[constants.WP].Test(constants.E5) {
		t.Error("Expected a white pawn at E5 after move")
	}
}
