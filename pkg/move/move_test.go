// Package move contains the move and move list representation and all move helper functions.
// Move is represented as 64 bit unsigned integer where some bits represent some part of the move
// The first 6 bits keep the source square, the next 6 bits keep the target square and so on...
package move

import (
	"testing"
)

func TestEncodeDecodeMove(t *testing.T) {
	m := EncodeMove(12, 28, 3, 5, 1, 0, 0, 1) // random values

	if m.GetSource() != 12 {
		t.Errorf("Expected source=12, got %d", m.GetSource())
	}
	if m.GetTarget() != 28 {
		t.Errorf("Expected target=28, got %d", m.GetTarget())
	}
	if m.GetPiece() != 3 {
		t.Errorf("Expected piece=3, got %d", m.GetPiece())
	}
	if m.GetPromoted() != 5 {
		t.Errorf("Expected promoted=5, got %d", m.GetPromoted())
	}
	if m.GetCapture() != 1 {
		t.Errorf("Expected capture=1, got %d", m.GetCapture())
	}
	if m.GetCastling() != 1 {
		t.Errorf("Expected castling=1, got %d", m.GetCastling())
	}
}
