// Package attacks contains the pre-computed attack tables for all pieces.
// For sliding pieces (bishop, rook, queens) it uses magic numbers for indexing
package attacks

import (
	"testing"

	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/constants"
)

func TestQueenAttacksCenter(t *testing.T) {
	// Initialize slider piece attacks before testing (assuming you have a function that does so).
	// For simplicity, assume BishopAttacks, RookAttacks, bishopMasks, rookMasks, and magic numbers are set.
	// If not, call InitSliderPiecesAttacks(constants.Bishop) and InitSliderPiecesAttacks(constants.Rook) here.
	InitSliderPiecesAttacks(constants.Bishop)
	InitSliderPiecesAttacks(constants.Rook)

	// Place no blockers (occupancy=0) and get queen attacks from D4 (27)
	sq := constants.D4
	occupancy := bitboard.Bitboard(0)

	att := GetQueenAttacks(sq, occupancy)
	if att.Count() == 0 {
		t.Errorf("Expected queen at D4 on empty board to have many attacks, got zero")
	}
}

func TestQueenAttacksWithBlockers(t *testing.T) {
	InitSliderPiecesAttacks(constants.Bishop)
	InitSliderPiecesAttacks(constants.Rook)
	sq := constants.D4
	// Occupancy with a blocker directly in line: place a blocker at D6=11?
	// Adjust indexing as needed based on your board indexing.
	blocker := bitboard.Bitboard(0)
	blocker.Set(constants.D6) // D6=11 if indexing top-left A8=0

	att := GetQueenAttacks(sq, blocker)
	// Ensure that squares beyond D6 in that direction aren't included:
	// Just a sanity check:
	if att.Test(constants.D7) {
		t.Error("Expected D7 not to be attacked due to blocker at D6")
	}
}
