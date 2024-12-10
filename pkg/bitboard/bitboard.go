// Package bitboard contains the Bitboard type and all helper functions
// for manipulating and querying bitboards. A bitboard is a 64-bit integer
// where each bit represents a square on a chessboard (or any 8x8 grid).
// Conventionally, bit 0 represents one corner and bit 63 the opposite corner.
package bitboard

import (
	"fmt"
	"math/bits"
)

// Bitboard represents a 64-bit board where each bit corresponds to a square.
type Bitboard uint64

// Count returns the number of set bits (ones) in the Bitboard.
// This is often useful for counting how many pieces exist on a board.
func (b Bitboard) Count() int {
	return bits.OnesCount64(uint64(b))
}

// Set sets the bit at the given position pos to 1.
// Positions are typically 0-based, ranging from 0 to 63.
func (b *Bitboard) Set(pos int) {
	*b |= Bitboard(uint64(1) << uint(pos))
}

// Test returns true if the bit at position pos is set to 1, false otherwise.
func (b *Bitboard) Test(pos int) bool {
	return *b&Bitboard(uint64(1)<<uint(pos)) != 0
}

// Clear sets the bit at position pos to 0.
func (b *Bitboard) Clear(pos int) {
	*b &= Bitboard(^(uint64(1) << uint(pos)))
}

// FirstOne returns the position of the least significant set bit (LSB).
// It also modifies the Bitboard by clearing all bits up to and including that position.
// If no bits are set, it returns 64
func (b *Bitboard) FirstOne() int {
	bit := bits.TrailingZeros64(uint64(*b))
	if bit == 64 {
		return 64
	}

	*b = (*b >> uint(bit+1)) << uint(bit+1)
	return bit
}

// LastOne returns the position of the most significant set bit (MSB).
// It also modifies the Bitboard by clearing all bits from that position upwards.
// If no bits are set, it returns 64.
func (b *Bitboard) LastOne() int {
	bit := bits.LeadingZeros64(uint64(*b))

	if bit == 64 {
		return 64
	}

	*b = (*b << uint(bit+1)) >> uint(bit+1)
	return 63 - bit
}

// String returns a string representation
func (b Bitboard) String() string {
	zeroes := ""
	for i := 0; i < 64; i++ {
		zeroes += "0"
	}

	bits := zeroes + fmt.Sprintf("%b", b)
	return bits[len(bits)-64:]
}

// PrintBitboard prints the Bitboard in an 8x8 grid with row and file indicators,
// similar to a chessboard layout. The top row is rank 8 and the bottom row is rank 1.
// '1' bits are shown, and '0' bits show empty squares.
func (b Bitboard) PrintBitboard() {
	s := b.String()
	row := [8]string{}
	row[7] = s[0:8]
	row[6] = s[8:16]
	row[5] = s[16:24]
	row[4] = s[24:32]
	row[3] = s[32:40]
	row[2] = s[40:48]
	row[1] = s[48:56]
	row[0] = s[56:]
	fmt.Println()
	for i, r := range row {
		fmt.Printf(
			"%d   %v %v %v %v %v %v %v %v\n", 8-i,
			r[7:8],
			r[6:7],
			r[5:6],
			r[4:5],
			r[3:4],
			r[2:3],
			r[1:2],
			r[0:1],
		)
	}
	fmt.Print("\n")
	fmt.Printf("    a b c d e f g h\n\n")

	fmt.Printf("Bitboard: %X\n\n", b)
}
