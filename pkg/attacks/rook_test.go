// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGenerateRookAttacks checks that some basic cases produce expected results.
func TestGenerateRookAttacks(t *testing.T) {
	// Consider a rook on D4
	ra := generateRookPossibleBlockers(D4)
	if ra.Count() == 0 {
		t.Error("Expected some rook attacks from D4, got none.")
	}

	// Check a corner like A1:
	raA1 := generateBishopPossibleBlockers(A1)
	if raA1.Count() == 0 {
		t.Error("Expected some rook attacks from A1, got none.")
	}
}

// TestGenerateRookAttacksOnTheFly checks the rook attack generation with blockers.
func TestGenerateRookAttacksOnTheFly(t *testing.T) {
	// Place a rook at D4 and create a blocker directly in line of one direction.
	block := bitboard.Bitboard(0)
	block.Set(D7)

	ra := generateRookAttacks(D4, block)

	// Check that D6 square is attacked:
	if !ra.Test(D6) {
		t.Error("Expected rook attacks to include D6.")
	}
	// Check that the blocker square (11) is attacked:
	if !ra.Test(D7) {
		t.Error("Expected rook attacks to include the blocker at D7.")
	}
	// Check a square above that line, say D7=3 if we continue going up, should not be included:
	if ra.Test(D8) {
		t.Error("Expected rook attacks not to include squares beyond the blocker.")
	}

	// Test a horizontal blockage:
	// Place a blocker at F4
	block = bitboard.Bitboard(0)
	block.Set(F4)

	ra = generateRookAttacks(D4, block)

	// Check that E4 square is attacked:
	if !ra.Test(E4) {
		t.Error("Expected rook attacks to include E4.")
	}
	if !ra.Test(F4) {
		t.Error("Expected rook attacks to include the blocker at F4.")
	}
	if ra.Test(G4) {
		t.Error("Expected rook attacks not to include squares beyond the horizontal blocker.")
	}
}
