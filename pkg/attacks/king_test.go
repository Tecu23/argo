// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/constants"
)

func TestInitKingAttacks(t *testing.T) {
	InitKingAttacks()

	// King on D4
	// The king moves one step in any of the 8 directions, if available.
	// From D4, king can move to: C5, D5, E5, E4, C4, C3, D3, E3
	kingAtD4 := KingAttacks[constants.D4]

	expectedSquares := []int{
		constants.C5,
		constants.D5,
		constants.E5,
		constants.E4,
		constants.C4,
		constants.C3,
		constants.D3,
		constants.E3,
	}
	for _, sq := range expectedSquares {
		if !kingAtD4.Test(sq) {
			t.Errorf("King at D4 should attack square %d", sq)
		}
	}
	if kingAtD4.Count() != 8 {
		t.Errorf("Expected king at D4 to have 8 attacks, got %d", kingAtD4.Count())
	}

	// King on A1
	// From A1, king can only move right B1, up A2, and up-right B2
	kingAtA1 := KingAttacks[constants.A1]
	if !kingAtA1.Test(constants.B1) {
		t.Error("King at A1 should attack B1")
	}
	if !kingAtA1.Test(constants.A2) {
		t.Error("King at A1 should attack A2")
	}
	if !kingAtA1.Test(constants.B2) {
		t.Error("King at A1 should attack B2")
	}
	if kingAtA1.Count() != 3 {
		t.Errorf("Expected king at A1 to have 3 moves, got %d", kingAtA1.Count())
	}

	// King on H8
	// H8 is top-right corner, can only move left G8, down H7, and down-left G7
	kingAtH8 := KingAttacks[constants.H8]
	if !kingAtH8.Test(constants.G8) {
		t.Error("King at H8 should attack G8")
	}
	if !kingAtH8.Test(constants.H7) {
		t.Error("King at H8 should attack H7")
	}
	if !kingAtH8.Test(constants.G7) {
		t.Error("King at H8 should attack G7")
	}
	if kingAtH8.Count() != 3 {
		t.Errorf("Expected king at H8 to have 3 moves, got %d", kingAtH8.Count())
	}
}
