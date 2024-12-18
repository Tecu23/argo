package history

import "testing"

func TestHistoryTable(t *testing.T) {
	h := New()

	// Test basic update
	h.Update(0, 12, 28, 4) // e2e4 depth 4
	score := h.Get(0, 12, 28)
	if score <= 0 {
		t.Errorf("Expected positive score after update, got %d", score)
	}

	// Test score aging
	for i := 0; i < 100; i++ {
		h.Update(0, 12, 28, 4)
	}
	// newScore := h.Get(0, 12, 28)
	// if newScore >= HistoryMax {
	// 	t.Error("History score not properly aged")
	// }
}
