// Copyright (C) 2025 Tecu23
// Port of Koivisto evaluation, licensed under GNU GPL v3

// Package nnue keeps the NNUE (Efficiently Updated Neural Network) responsible for
// evaluation the current position
package nnue

import "fmt"

// KingSquareIndices maps board squares to bucketed indices for king positions.
// The board is divided into buckets so that similar king positions share accumulator data.
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

// KingSquareIndex computes the bucket index for a given king square and color.
// It first transforms the square index based on the king's perspective and then retrieves the bucket index.
func KingSquareIndex(kingSquare int, kingColor int) int {
	// Safety check: ensure the king's square is within valid bounds
	if kingSquare < 0 || kingSquare >= 64 {
		fmt.Printf("ERROR: Invalid king square: %d\n", kingSquare)
		return 0 // Return a safe default index
	}

	// Transform the square index based on the king's color perspective
	kingSquare = ((56 * kingColor) ^ kingSquare)

	return KingSquareIndices[kingSquare]
}

// Index computes the feature index for a piece on a square given the piece type, color, and current perspective.
// This index is used to access the correct weights in the input layer.
func Index(pieceType, pieceColor, square, view, kingSquare int) int {
	// The parameters 'square' and 'kingSquare' are preserved for clarity.
	ksIndex := KingSquareIndex(kingSquare, view)
	// Transform the square index from one representation to the NNUE expected format
	square ^= 56 * view
	// If the king is on a specific bucket (as indicated by the 0x4 bit), adjust the square index accordingly
	if kingSquare&0x4 != 0 {
		square ^= 7
	}

	// Compute the overall index by combining square index, piece type and color,
	// and the bucket index for the king.
	return square + pieceType*64 + boolToInt(pieceColor == view)*64*6 + ksIndex*64*6*2
}

// FeatureIndex represents a combination of a piece type, its color, and its square.
// It also stores the king positions for both sides for proper indexing.
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

	// Use the global Index function with the proper king square
	return Index(f.PieceType, f.PieceColor, f.Square, side, kingSq)
}

// boolToInt converts a boolean value to its integer equivalent (1 if true, 0 if false).
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
