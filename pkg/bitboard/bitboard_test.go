package bitboard

import (
	"testing"
)

func TestBitboardSetAndTest(t *testing.T) {
	var b Bitboard
	if b.Test(0) {
		t.Error("Expected bit 0 to be clear initially")
	}

	b.Set(0)
	if !b.Test(0) {
		t.Error("Expected bit 0 to be set after calling Set(0)")
	}

	b.Set(63)
	if !b.Test(63) {
		t.Error("Expected bit 63 to be set after calling Set(63)")
	}

	if b.Test(32) {
		t.Error("Expected bit 32 to be clear, never set it")
	}
}

func TestBitboardClear(t *testing.T) {
	var b Bitboard
	b.Set(10)
	b.Set(20)
	if !b.Test(10) || !b.Test(20) {
		t.Error("Expected bits 10 and 20 to be set")
	}

	b.Clear(10)
	if b.Test(10) {
		t.Error("Expected bit 10 to be cleared")
	}

	if !b.Test(20) {
		t.Error("Expected bit 20 to remain set")
	}
}

func TestBitboardCount(t *testing.T) {
	var b Bitboard
	for i := 0; i < 64; i += 2 {
		b.Set(i) // set every other bit
	}
	count := b.Count()
	if count != 32 {
		t.Errorf("Expected 32 bits set, got %d", count)
	}
}

func TestBitboardFirstOne(t *testing.T) {
	var b Bitboard
	b.Set(5)
	b.Set(10)
	b.Set(63)

	first := b.FirstOne()
	if first != 5 {
		t.Errorf("Expected first one at position 5, got %d", first)
	}

	// After removing first one, next first should be at 10
	next := b.FirstOne()
	if next != 10 {
		t.Errorf("Expected next first one at position 10, got %d", next)
	}

	// After removing that, next first should be at 63
	last := b.FirstOne()
	if last != 63 {
		t.Errorf("Expected last first one at position 63, got %d", last)
	}

	// Now it should be empty
	empty := b.FirstOne()
	if empty != 64 {
		t.Errorf("Expected no bits left, got position %d", empty)
	}
}

func TestBitboardLastOne(t *testing.T) {
	var b Bitboard
	b.Set(2)
	b.Set(40)
	b.Set(63)

	last := b.LastOne()
	if last != 63 {
		t.Errorf("Expected last one at position 63, got %d", last)
	}

	// After removing that, last one should be at 40
	next := b.LastOne()
	if next != 40 {
		t.Errorf("Expected next last one at position 40, got %d", next)
	}

	// After removing that, last one should be at 2
	lowest := b.LastOne()
	if lowest != 2 {
		t.Errorf("Expected last one at position 2, got %d", lowest)
	}

	// Now empty
	empty := b.LastOne()
	if empty != 64 {
		t.Errorf("Expected no bits left, got position %d", empty)
	}
}

func TestBitboardString(t *testing.T) {
	var b Bitboard
	b.Set(0)
	b.Set(63)

	strRep := b.String()
	if len(strRep) != 64 {
		t.Errorf("Expected string of length 64, got %d", len(strRep))
	}

	// bit 0 set means least significant bit set
	// bit 63 set means most significant bit set
	// The string prints with bit 63 as first char and bit 0 as last char if read normally.
	// But here we rely on the direct binary representation:
	if strRep[63] != '1' {
		t.Error("Expected the least significant bit (strRep[63]) to be '1'")
	}

	if strRep[0] != '1' {
		t.Error("Expected the most significant bit (strRep[0]) to be '1'")
	}
}
