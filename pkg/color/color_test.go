package color

import (
	"testing"
)

// TestColorConstants verifies that the Color constants are assigned correctly.
func TestColorConstants(t *testing.T) {
	if WHITE != Color(0) {
		t.Errorf("Expected WHITE to be 0, got %d", WHITE)
	}
	if BLACK != Color(1) {
		t.Errorf("Expected BLACK to be 1, got %d", BLACK)
	}
	if BOTH != Color(2) {
		t.Errorf("Expected BOTH to be 2, got %d", BOTH)
	}
}

// TestColorOpp checks that the Opp function returns the expected opposite colors.
func TestColorOpp(t *testing.T) {
	// WHITE -> BLACK
	if WHITE.Opp() != BLACK {
		t.Errorf("WHITE.Opp() should return BLACK, got %v", WHITE.Opp())
	}
	// BLACK -> WHITE
	if BLACK.Opp() != WHITE {
		t.Errorf("BLACK.Opp() should return WHITE, got %v", BLACK.Opp())
	}
	// BOTH -> BLACK (due to XOR operation)
	if BOTH.Opp() != BOTH {
		t.Errorf("BOTH.Opp() should return BOTH, got %v", BOTH.Opp())
	}
}

// TestColorString checks the String method for correctness.
func TestColorString(t *testing.T) {
	// WHITE should return "W"
	if WHITE.String() != "W" {
		t.Errorf("Expected WHITE.String() to be 'W', got '%s'", WHITE.String())
	}
	// BLACK should return "B"
	if BLACK.String() != "B" {
		t.Errorf("Expected BLACK.String() to be 'B', got '%s'", BLACK.String())
	}
	// BOTH also returns "B" as it falls under the default case
	if BOTH.String() != "B" {
		t.Errorf("Expected BOTH.String() to be 'B', got '%s'", BOTH.String())
	}
}
