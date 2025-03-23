// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package reduction

import (
	"testing"
)

func TestReductionTable(t *testing.T) {
	table := New()

	tests := []struct {
		depth      int
		moveNumber int
		want       int
	}{
		{2, 5, 0},   // Too shallow depth
		{3, 3, 0},   // Too early move
		{6, 10, 1},  // Normal reduction
		{10, 20, 2}, // Deeper reduction
		{20, 30, 3}, // Max reduction
		{63, 63, 4}, // Maximum values
	}

	for _, tt := range tests {
		got := table.Get(tt.depth, tt.moveNumber)
		if got != tt.want {
			t.Errorf("Reduction(%d, %d) = %d; want %d",
				tt.depth, tt.moveNumber, got, tt.want)
		}
	}
}

func TestReductionLimits(t *testing.T) {
	table := New()

	// Test out of bounds
	if r := table.Get(64, 0); r != 0 {
		t.Errorf("Out of bounds depth should return 0, got %d", r)
	}
	if r := table.Get(0, 64); r != 0 {
		t.Errorf("Out of bounds moveNumber should return 0, got %d", r)
	}

	// Test maximum reduction
	maxReduction := 0
	for d := 0; d < 64; d++ {
		for m := 0; m < 64; m++ {
			r := table.Get(d, m)
			if r > maxReduction {
				maxReduction = r
			}
			if r > MaxReduction {
				t.Errorf("Reduction(%d, %d) = %d exceeds MaxReduction %d",
					d, m, r, MaxReduction)
			}
		}
	}
}
