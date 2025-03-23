// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

// Package hash
package hash

import "math/rand"

const (
	PieceSquareCount = 12 * 64 // 12 piece types * 64 squares
	CastlingCount    = 4       // 4 castling rights
	EnPassantCount   = 8       // 8 files for en passant
	SideCount        = 1       // Side to Move
)

type ZobristHash struct {
	PieceSquare [PieceSquareCount]uint64 // Random numbers for each piece on each square
	Castling    [CastlingCount]uint64    // Random numbers for castling rights
	EnPassant   [EnPassantCount]uint64   // Random numbers for en passant files
	Side        uint64                   // Random number for side to move
}

var HashTable *ZobristHash

func Init() {
	// use a fixed seed for reproducibility
	const seed = 23

	HashTable = NewZobristHash(seed)
}

func NewZobristHash(seed int64) *ZobristHash {
	rng := rand.New(rand.NewSource(seed))

	z := &ZobristHash{}

	// Initialize piece square values
	for i := 0; i < PieceSquareCount; i++ {
		z.PieceSquare[i] = rng.Uint64()
	}

	// Initialize castlings values
	for i := 0; i < CastlingCount; i++ {
		z.Castling[i] = rng.Uint64()
	}

	// Initialize en passant values
	for i := 0; i < EnPassantCount; i++ {
		z.EnPassant[i] = rng.Uint64()
	}

	// Initialize side to move value
	z.Side = rng.Uint64()

	return z
}
