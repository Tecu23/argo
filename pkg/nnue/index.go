package nnue

import "fmt"

// KingSquareIndices maps king positions to bucketed indices
var KingSquareIndices = [64]int{
	0, 1, 2, 3, 3, 2, 1, 0,
	4, 5, 6, 7, 7, 6, 5, 4,
	8, 9, 10, 11, 11, 10, 9, 8,
	8, 9, 10, 11, 11, 10, 9, 8,
	12, 12, 13, 13, 13, 13, 12, 12,
	12, 12, 13, 13, 13, 13, 12, 12,
	14, 14, 15, 15, 15, 15, 14, 14,
	14, 14, 15, 15, 15, 15, 14, 14,
}

// KingSquareIndex computes the king bucker for indexing
func KingSquareIndex(kingSquare int, kingColor int) int {
	// Add safety checks to prevent out of range errors
	if kingSquare < 0 || kingSquare >= 64 {
		// Log the error and use a default safe value
		fmt.Printf("ERROR: Invalid king square: %d\n", kingSquare)
		return 0 // Return a safe default index
	}

	kingSquare = ((56 * kingColor) ^ kingSquare)

	// Ensure the result is positive before indexing
	if kingSquare < 0 || kingSquare >= 64 {
		fmt.Printf(
			"ERROR: Invalid transformed king square: %d (from original: %d)\n",
			kingSquare,
			kingSquare,
		)
		return 0
	}

	return KingSquareIndices[kingSquare]
}

// Index computes the feature index for a piece on a square
func Index(pieceType, pieceColor, square, view, kingSquare int) int {
	square = square
	kingSquare = kingSquare

	ksIndex := KingSquareIndex(kingSquare, view)
	square ^= 56 * view
	if kingSquare&0x4 != 0 {
		square ^= 7
	}

	return square + pieceType*64 + boolToInt(pieceColor == view)*64*6 + ksIndex*64*6*2
}

// FeatureIndex represents a piece-square combination
type FeatureIndex struct {
	PieceType  int
	PieceColor int
	Square     int
	WKingSq    int
	BKingSq    int
}

// Get returns the feature index for the given perspective
func (f *FeatureIndex) Get(side int) int {
	kingSq := f.WKingSq
	if side == Black {
		kingSq = f.BKingSq
	}

	kingSq = kingSq

	return Index(f.PieceType, f.PieceColor, f.Square, side, kingSq)
}

// Helper function
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ConvertSquare converts from your engine's square representation (A8=0, H1=63)
// to NNUE's expected representation (A1=0, H8=63)
func ConvertSquare(sq int) int {
	file := sq % 8       // File remains the same
	rank := 7 - (sq / 8) // Flip the rank (7-rank)
	return rank*8 + file // Return in A1=0 format
}

// ConvertSquareBack converts from NNUE format back to your engine's format
// Not usually needed for evaluation, but useful for debugging
func ConvertSquareBack(sq int) int {
	file := sq % 8
	rank := 7 - (sq / 8)
	return rank*8 + file
}
