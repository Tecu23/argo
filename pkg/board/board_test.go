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
				b.Side = color.BLACK
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
	if b.Side != copy.Side {
		return fmt.Errorf("side mismatch")
	}
	if b.EnPassant != copy.EnPassant {
		return fmt.Errorf("en passant mismatch")
	}
	if b.Rule50 != copy.Rule50 {
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

func TestPassedPawns(t *testing.T) {
	tests := []struct {
		name         string
		fen          string
		whiteSquares []int  // Squares where white should have passed pawns
		blackSquares []int  // Squares where black should have passed pawns
		desc         string // Description of what we're testing
	}{
		{
			name:         "Empty Board",
			fen:          "8/8/8/8/8/8/8/8 w - - 0 1",
			whiteSquares: []int{},
			blackSquares: []int{},
			desc:         "No pawns on board",
		},
		{
			name:         "Starting Position",
			fen:          "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			whiteSquares: []int{},
			blackSquares: []int{},
			desc:         "Initial position, no passed pawns",
		},
		{
			name:         "Simple Passed Pawn",
			fen:          "8/8/8/4P3/8/8/8/8 w - - 0 1",
			whiteSquares: []int{E5}, // e5
			blackSquares: []int{},
			desc:         "Single white passed pawn with no obstacles",
		},
		{
			name:         "Blocked Pawn",
			fen:          "8/4p3/4P3/8/8/8/8/8 w - - 0 1",
			whiteSquares: []int{},
			blackSquares: []int{},
			desc:         "Pawn blocked directly in front",
		},
		{
			name:         "Diagonal Threats",
			fen:          "8/3p4/4P3/8/8/8/8/8 w - - 0 1",
			whiteSquares: []int{E6},
			blackSquares: []int{D7},
			desc:         "Pawn threatened diagonally",
		},
		{
			name:         "Multiple Passed Pawns",
			fen:          "8/8/2P2P2/8/8/8/8/8 w - - 0 1",
			whiteSquares: []int{C6, F6},
			blackSquares: []int{},
			desc:         "Two white passed pawns",
		},
		{
			name:         "Rook Pawns",
			fen:          "8/8/P6p/8/8/8/8/8 w - - 0 1",
			whiteSquares: []int{A6},
			blackSquares: []int{H6},
			desc:         "Passed pawns on the a and h files",
		},
		{
			name:         "Protected Passed Pawn",
			fen:          "8/8/8/3P4/2P5/8/8/8 w - - 0 1",
			whiteSquares: []int{D5, C4},
			blackSquares: []int{},
			desc:         "Passed pawn protected by another pawn",
		},
		{
			name:         "Complex Position",
			fen:          "8/p1p3pp/P1P5/8/8/8/8/8 w - - 0 1",
			whiteSquares: []int{},
			blackSquares: []int{G7, H7},
			desc:         "Multiple passed pawns for both sides",
		},
		{
			name:         "Mutual Blockade",
			fen:          "8/8/8/2pPp3/8/8/8/8 w - - 0 1",
			whiteSquares: []int{D5},
			blackSquares: []int{C5, E5},
			desc:         "Pawns blocking each other",
		},
		{
			name:         "Distant Blockers",
			fen:          "8/3p4/8/4P3/8/8/8/8 w - - 0 1",
			whiteSquares: []int{},
			blackSquares: []int{},
			desc:         "Pawn blocked by distant enemy pawn",
		},
		{
			name:         "Phalanx Formation",
			fen:          "8/8/8/3PPP2/8/8/8/8 w - - 0 1",
			whiteSquares: []int{D5, E5, F5},
			blackSquares: []int{},
			desc:         "Connected passed pawns",
		},
		{
			name:         "Lever Position",
			fen:          "8/8/8/3P4/2p5/8/8/8 w - - 0 1",
			whiteSquares: []int{D5},
			blackSquares: []int{C4},
			desc:         "Pawn that can be captured by lever",
		},
		{
			name:         "Advanced Passed Pawns",
			fen:          "8/P7/8/8/8/8/p7/8 w - - 0 1",
			whiteSquares: []int{A7},
			blackSquares: []int{A2},
			desc:         "Very advanced passed pawns for both sides",
		},
		{
			name:         "Chain Formation",
			fen:          "8/8/3P4/4P3/5P2/8/8/8 w - - 0 1",
			whiteSquares: []int{F4, E5, D6},
			blackSquares: []int{},
			desc:         "Pawn chain with passed pawn at front",
		},
		{
			name:         "Backward Pawns",
			fen:          "8/8/8/2P1P3/3P4/8/8/8 w - - 0 1",
			whiteSquares: []int{C5, E5, D4},
			blackSquares: []int{},
			desc:         "Formation with backward pawn",
		},
		{
			name:         "Edge Cases",
			fen:          "8/P6p/8/8/8/p6P/8/8 w - - 0 1",
			whiteSquares: []int{A7},
			blackSquares: []int{A3},
			desc:         "Various edge case positions",
		},
		{
			name:         "Doubled Pawns",
			fen:          "8/3P4/3P4/8/8/8/8/8 w - - 0 1",
			whiteSquares: []int{D7},
			blackSquares: []int{},
			desc:         "Doubled pawns with one passed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			board, _ := ParseFEN(test.fen)

			// Test white passed pawns
			for sq := 0; sq < 64; sq++ {
				isExpectedPassed := false
				for _, passedSq := range test.whiteSquares {
					if sq == passedSq {
						isExpectedPassed = true
						break
					}
				}
				isPassed := board.IsPassedPawn(sq)
				if isPassed != isExpectedPassed {
					t.Errorf("White pawn on %s: got passed=%v, want passed=%v\nFEN: %s\nDesc: %s",
						util.Sq2Fen[sq], isPassed, isExpectedPassed, test.fen, test.desc)
				}
			}

			// Test black passed pawns by mirroring
			mirroredBoard := board.Mirror()
			for sq := 0; sq < 64; sq++ {
				mirroredSq := sq ^ 56 // Flip square vertically
				isExpectedPassed := false
				for _, passedSq := range test.blackSquares {
					if sq == passedSq {
						isExpectedPassed = true
						break
					}
				}

				isPassed := mirroredBoard.IsPassedPawn(mirroredSq)
				if isPassed != isExpectedPassed {
					t.Errorf("Black pawn on %s: got passed=%v, want passed=%v\nFEN: %s\nDesc: %s",
						util.Sq2Fen[sq], isPassed, isExpectedPassed, test.fen, test.desc)
				}
			}

			// Test total counts
			whiteCount := board.CandidatePassed(color.WHITE)
			blackCount := mirroredBoard.CandidatePassed(color.WHITE)

			if whiteCount != len(test.whiteSquares) {
				t.Errorf("White passed pawn count: got %d, want %d\nFEN: %s\nDesc: %s",
					whiteCount, len(test.whiteSquares), test.fen, test.desc)
			}

			if blackCount != len(test.blackSquares) {
				t.Errorf("Black passed pawn count: got %d, want %d\nFEN: %s\nDesc: %s",
					blackCount, len(test.blackSquares), test.fen, test.desc)
			}
		})
	}
}
