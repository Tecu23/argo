// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/color"
	"github.com/Tecu23/argov2/pkg/constants"
)

func TestInitPawnAttacks(t *testing.T) {
	InitPawnAttacks()

	// White pawn on A2
	// White pawns move "up" from White's perspective, so A2 attacks B3 if possible.
	whiteA2Attacks := PawnAttacks[color.WHITE][constants.A2]
	if whiteA2Attacks.Count() != 1 {
		t.Errorf("Expected A2 (white) to have 1 attack (B3), got %d", whiteA2Attacks.Count())
	}
	if !whiteA2Attacks.Test(constants.B3) {
		t.Error("Expected A2 (white) to attack B3")
	}

	// White pawn on H2
	// Similar logic, G3
	whiteH2Attacks := PawnAttacks[color.WHITE][constants.H2]
	if whiteH2Attacks.Count() != 1 {
		t.Errorf("Expected H2 (white) to have 1 attack (G3), got %d", whiteH2Attacks.Count())
	}
	if !whiteH2Attacks.Test(constants.G3) {
		t.Error("Expected H2 (white) to attack G3")
	}

	// Black pawn on A7 should attack only B6
	blackA7Attacks := PawnAttacks[color.BLACK][constants.A7]
	if blackA7Attacks.Count() != 1 {
		t.Errorf("Expected A7 (black) to have 1 attack (B6), got %d", blackA7Attacks.Count())
	}
	if !blackA7Attacks.Test(constants.B6) {
		t.Error("Expected A7 (black) to attack B6")
	}

	// Black pawn on H7
	blackH7Attacks := PawnAttacks[color.BLACK][constants.H7]
	if blackH7Attacks.Count() != 1 {
		t.Errorf("Expected H7 (black) to have 1 attack (G6), got %d", blackH7Attacks.Count())
	}
	// G6
	if !blackH7Attacks.Test(constants.G6) {
		t.Error("Expected H7 (black) to attack G6")
	}
}

func TestGeneratePawnAttacksEdges(t *testing.T) {
	// Test edge cases: White pawn on A8 and Black pawn on H1, for example.
	// A8 for White: can't attack NW (off board), can't attack NE (no square above)
	whiteA8 := generatePawnAttacks(constants.A8, color.WHITE)
	if whiteA8.Count() != 0 {
		t.Errorf("Expected A8 (white) to have no attacks, got %d", whiteA8.Count())
	}

	// H1 for Black: can't attack SW or SE (off the board)
	blackH1 := generatePawnAttacks(constants.H1, color.BLACK)
	if blackH1.Count() != 0 {
		t.Errorf("Expected H8 (black) to have no attacks, got %d", blackH1.Count())
	}
}

func TestGeneratePawnAttacksMiddle(t *testing.T) {
	// White pawn in the middle of the board, say D4:
	// White pawn at D4 attacks C5 and E5
	whiteD4 := generatePawnAttacks(constants.D4, color.WHITE)
	if whiteD4.Count() != 2 {
		t.Errorf("Expected D4 (white) to have 2 attacks (C5 and E5), got %d", whiteD4.Count())
	}
	if !whiteD4.Test(constants.C5) {
		t.Error("Expected D4 (white) to attack C5")
	}
	if !whiteD4.Test(constants.E5) {
		t.Error("Expected D4 (white) to attack E5")
	}
}
