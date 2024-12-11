// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/bitboard"
	. "github.com/Tecu23/argov2/pkg/constants"
)

// TestGenerateBishopAttacks verifies the bishop's diagonal moves on an empty board (no blockers).
func TestGenerateBishopAttacks(t *testing.T) {
	// Let's pick a square in the middle: D4
	bb := generateBishopPossibleBlockers(D4)
	if bb.Count() == 0 {
		t.Error("Expected some attacks from D4, got none")
	}

	// Let's pick a corner like A1 (bottom-left corner)
	bbA1 := generateBishopPossibleBlockers(A1)
	if bbA1.Count() == 0 {
		t.Error("Expected some diagonal squares from A1, got none")
	}
}

// TestGenerateBishopAttacksOnTheFly checks bishop attacks with blockers.
func TestGenerateBishopAttacksOnTheFly(t *testing.T) {
	// Place a bishop at D4 (27) on an empty board except one blocker at F6 (which is rank 2 down from top?)
	// We'll block one diagonal direction and ensure the bishop stops there.
	block := bitboard.Bitboard(0)
	block.Set(F6) // place blocker at F6

	bb := generateBishopAttacks(D4, block)

	// Check that squares on the diagonal beyond F6 are not included:
	// G7 would be further up-right from F6 (but since we hit a blocker at F6, we must not have G7).
	if bb.Test(G7) {
		t.Error("Expected bishop attacks to be blocked at F6 and not include G7")
	}

	// Check that F6 itself is included (the bishop attacks the blocker square):
	if !bb.Test(F6) {
		t.Error("Expected bishop to include the blocked square F6 in its attacks")
	}
}
