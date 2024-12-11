// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"testing"

	. "github.com/Tecu23/argov2/pkg/constants"
)

func TestInitKnightAttacks(t *testing.T) {
	InitKnightAttacks()

	// Test a known position:
	// A knight on B1 (square B1 = 1) traditionally can move to A3, C3 and D2.
	knightAtB1 := KnightAttacks[B1]
	if !knightAtB1.Test(A3) {
		t.Errorf("Knight at B1 should attack A3, but does not.")
	}
	if !knightAtB1.Test(C3) {
		t.Errorf("Knight at B1 should attack C3, but does not.")
	}
	if !knightAtB1.Test(D2) {
		t.Errorf("Knight at B1 should attack D2, but does not.")
	}
	if knightAtB1.Count() != 3 {
		t.Errorf("Expected knight at B1 to have 3 attacks, got %d", knightAtB1.Count())
	}

	// Knight on D4
	// A knight on D4 can move to: B5, C6, E6, F5, F3, E2, C2, B3
	knightAtD4 := KnightAttacks[D4]
	if !knightAtD4.Test(B5) {
		t.Error("Knight at D4 should attack B5")
	}
	if !knightAtD4.Test(C6) {
		t.Error("Knight at D4 should attack C6")
	}
	if !knightAtD4.Test(E6) {
		t.Error("Knight at D4 should attack E6")
	}
	if !knightAtD4.Test(F5) {
		t.Error("Knight at D4 should attack F5")
	}
	if !knightAtD4.Test(F3) {
		t.Error("Knight at D4 should attack F3")
	}
	if !knightAtD4.Test(E2) {
		t.Error("Knight at D4 should attack E2")
	}
	if !knightAtD4.Test(C2) {
		t.Error("Knight at D4 should attack C2")
	}
	if !knightAtD4.Test(B3) {
		t.Error("Knight at D4 should attack B3")
	}
}

func TestKnightAttacksCorner(t *testing.T) {
	InitKnightAttacks()

	// Knight on A1 should attack C2 and B3.
	knightAtA1 := KnightAttacks[A1]
	if !knightAtA1.Test(C2) {
		t.Error("Knight at A1 should attack C2")
	}
	if !knightAtA1.Test(B3) {
		t.Error("Knight at A1 should attack B3")
	}
	if knightAtA1.Count() != 2 {
		t.Errorf("Expected knight at A1 to have 2 attacks, got %d", knightAtA1.Count())
	}
}
