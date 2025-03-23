// Package board contains the board representation and all board helper functions.
// This package will handle move generation
package board

import (
	"fmt"
	"testing"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/bitboard"
	"github.com/Tecu23/argov2/pkg/color"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	util.InitFen2Sq() // Make sure square mappings are initialized
	hash.Init()
}

func TestSetSq(t *testing.T) {
	tests := []struct {
		name     string
		piece    int
		square   int
		expected struct {
			pieceBoard  bitboard.Bitboard
			occupancies [3]bitboard.Bitboard
		}
	}{
		{
			name:   "Set white pawn on e4",
			piece:  WP,
			square: E4,
			expected: struct {
				pieceBoard  bitboard.Bitboard
				occupancies [3]bitboard.Bitboard
			}{
				pieceBoard: bitboard.Bitboard(1) << uint(E4),
				occupancies: [3]bitboard.Bitboard{
					color.WHITE: bitboard.Bitboard(1) << uint(E4),
					color.BLACK: 0,
					color.BOTH:  bitboard.Bitboard(1) << uint(E4),
				},
			},
		},
		{
			name:   "Replace black piece with white piece",
			piece:  WQ,
			square: D5,
			expected: struct {
				pieceBoard  bitboard.Bitboard
				occupancies [3]bitboard.Bitboard
			}{
				pieceBoard: bitboard.Bitboard(1) << uint(D5),
				occupancies: [3]bitboard.Bitboard{
					color.WHITE: bitboard.Bitboard(1) << uint(D5),
					color.BLACK: 0,
					color.BOTH:  bitboard.Bitboard(1) << uint(D5),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Board{}
			b.SetSq(tt.piece, tt.square)

			// Check piece bitboard
			if b.Bitboards[tt.piece] != tt.expected.pieceBoard {
				t.Errorf("Piece bitboard mismatch: got %v, want %v",
					b.Bitboards[tt.piece], tt.expected.pieceBoard)
			}

			// Check occupancy bitboards
			for clr := color.WHITE; clr <= color.BLACK; clr++ {
				if b.Occupancies[clr] != tt.expected.occupancies[clr] {
					t.Errorf("Occupancy bitboard for color %d mismatch: got %v, want %v",
						clr, b.Occupancies[clr], tt.expected.occupancies[clr])
				}
			}
		})
	}
}

func TestCopyAndTakeBack(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(*Board)
		modify func(*Board)
	}{
		{
			name:  "Copy empty board",
			setup: func(b *Board) {},
			modify: func(b *Board) {
				b.SetSq(WP, E4)
			},
		},
		{
			name: "Copy with pieces",
			setup: func(b *Board) {
				b.SetSq(WK, E1)
				b.SetSq(BK, E8)
				b.SetSq(WP, E2)
			},
			modify: func(b *Board) {
				b.SetSq(WP, E4)
				b.SideToMove = color.BLACK
				b.EnPassant = E3
			},
		},
		{
			name: "Copy with en passant",
			setup: func(b *Board) {
				b.SetSq(WP, E4)
				b.SetSq(BP, D4)
				b.EnPassant = E3
			},
			modify: func(b *Board) {
				b.SetSq(Empty, E4)
				b.EnPassant = 0
			},
		},
		{
			name: "Copy with castling rights",
			setup: func(b *Board) {
				b.SetSq(WK, E1)
				b.SetSq(WR, H1)
				b.Castlings = Castlings(ShortW)
			},
			modify: func(b *Board) {
				b.SetSq(WK, G1)
				b.SetSq(WR, F1)
				b.Castlings = 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create and setup original board
			original := &Board{}
			tt.setup(original)

			// Create copy and verify it matches
			copy := original.CopyBoard()

			// Verify initial copy matches
			if err := verifyBoardsMatch(original, &copy); err != nil {
				t.Errorf("Initial copy failed: %v", err)
			}

			// Modify original board
			tt.modify(original)

			// Verify boards are different
			if !boardsAreDifferent(original, &copy) {
				t.Error("Boards should be different after modification")
			}

			// Take back to copied position
			original.TakeBack(copy)

			// Verify boards match again
			if err := verifyBoardsMatch(original, &copy); err != nil {
				t.Errorf("After takeback: %v", err)
			}
		})
	}
}

// Helper function to verify if two board positions match
func verifyBoardsMatch(b *Board, copy *Board) error {
	for i := range b.Bitboards {
		if b.Bitboards[i] != copy.Bitboards[i] {
			return fmt.Errorf("piece bitboard %d mismatch", i)
		}
	}
	for i := range b.Occupancies {
		if b.Occupancies[i] != copy.Occupancies[i] {
			return fmt.Errorf("occupancy bitboard %d mismatch", i)
		}
	}
	if b.SideToMove != copy.SideToMove {
		return fmt.Errorf("side mismatch")
	}
	if b.EnPassant != copy.EnPassant {
		return fmt.Errorf("en passant mismatch")
	}
	if b.HalfMoveClock != copy.HalfMoveClock {
		return fmt.Errorf("rule 50 mismatch")
	}
	if b.Castlings != copy.Castlings {
		return fmt.Errorf("castling rights mismatch")
	}
	return nil
}

// Helper function to verify if two boards are different
func boardsAreDifferent(b *Board, copy *Board) bool {
	return verifyBoardsMatch(b, copy) != nil
}

// func TestHashConsistency(t *testing.T) {
// 	tests := []struct {
// 		name         string
// 		moves        []string // Sequence of moves
// 		expectedHash uint64   // Known correct hash
// 	}{
// 		{
// 			name:  "Starting Position",
// 			moves: []string{},
// 			// Calculate expectedHash offline for start position
// 		},
// 		{
// 			name:  "e4 e5",
// 			moves: []string{"e2e4", "e7e5"},
// 		},
// 		{
// 			name:  "Same Position Different Move Order",
// 			moves: []string{"e2e4", "e7e5", "g1f3", "b8c6"},
// 			// Should match hash of: g1f3, b8c6, e2e4, e7e5
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b := NewBoard()
//
// 			// Make moves
// 			for _, move := range tt.moves {
// 				b2, ok := b.ParseMove(move)
// 				if !ok {
// 					t.Fatalf("Failed to make move: %s", move)
// 				}
// 				b = &b2
// 			}
//
// 			// Test make/unmake doesn't change hash
// 			moves := b.GenerateMoves()
// 			originalHash := b.Hash()
//
// 			for _, move := range moves {
// 				copyB := b.CopyBoard()
// 				if b.MakeMove(move, AllMoves) {
// 					b.TakeBack(copyB)
// 					if b.Hash() != originalHash {
// 						t.Errorf("Hash changed after make/unmake: %v", move)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

func BenchmarkPerftDriver(b *testing.B) {
	board, _ := ParseFEN(StartPosition)

	for i := 0; i < b.N; i++ {
		PerftDriver(&board, 7)
	}
}
