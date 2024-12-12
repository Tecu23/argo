package board

import (
	"fmt"
	"testing"

	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/constants"
	. "github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/move"
	"github.com/Tecu23/argov2/pkg/util"
)

// TestPosition represents a test case for move generation
type TestPosition struct {
	name          string
	fen           string
	expectedMoves int
	specific      []move.Move // specific moves to check for
}

func init() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(constants.Bishop)
	attacks.InitSliderPiecesAttacks(constants.Rook)

	util.InitFen2Sq()
}

func TestMoveGeneration(t *testing.T) {
	testCases := []TestPosition{
		// Initial position tests
		{
			name:          "Initial position",
			fen:           "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			expectedMoves: 20, // 16 pawn moves + 4 knight moves
		},

		// Pawn move tests
		{
			name:          "White pawns with captures",
			fen:           "rnbqkbnr/ppp1p1pp/8/3p1p2/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			expectedMoves: 20,
		},
		{
			name:          "Black pawns with en passant",
			fen:           "rnbqkbnr/ppp1p1pp/8/8/3Pp3/8/PPP2PPP/RNBQKBNR b KQkq d3 0 1",
			expectedMoves: 30, // Including en passant capture
			specific: []move.Move{
				move.EncodeMove(E4, D3, BP, 0, 1, 0, 1, 0), // en passant capture
			},
		},
		{
			name:          "White pawn promotions",
			fen:           "8/P7/8/8/8/8/8/8 w - - 0 1",
			expectedMoves: 4, // Queen, Rook, Bishop, Knight promotions
		},

		// Knight move tests
		{
			name:          "Knight in center",
			fen:           "8/8/8/8/3N4/8/8/8 w - - 0 1",
			expectedMoves: 8, // All possible knight moves
		},
		{
			name:          "Knight on edge",
			fen:           "8/8/8/8/N7/8/8/8 w - - 0 1",
			expectedMoves: 4, // Edge knight moves
		},

		// Bishop move tests
		{
			name:          "Bishop in center",
			fen:           "8/8/8/8/3B4/8/8/8 w - - 0 1",
			expectedMoves: 13, // All diagonal moves
		},
		{
			name:          "Bishop with blocking pieces",
			fen:           "8/8/1p3p2/8/3B4/8/8/8 w - - 0 1",
			expectedMoves: 10, // Including captures
		},

		// Rook move tests
		{
			name:          "Rook in center",
			fen:           "8/8/8/8/3R4/8/8/8 w - - 0 1",
			expectedMoves: 14, // All horizontal and vertical moves
		},
		{
			name:          "Rook with blocking pieces",
			fen:           "8/8/3p4/8/3R4/3P4/8/8 w - - 0 1",
			expectedMoves: 9, // Including friendly blocking
		},

		// Queen move tests
		{
			name:          "Queen in center",
			fen:           "8/8/8/8/3Q4/8/8/8 w - - 0 1",
			expectedMoves: 27, // All queen moves
		},

		// King move tests
		{
			name:          "King normal moves",
			fen:           "8/8/8/8/3K4/8/8/8 w - - 0 1",
			expectedMoves: 8, // All king moves
		},
		{
			name:          "White kingside castling",
			fen:           "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			expectedMoves: 26, // Including castling
			specific: []move.Move{
				move.EncodeMove(E1, G1, WK, 0, 0, 0, 0, 1), // kingside castle
			},
		},
		{
			name: "Black queenside castling blocked",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			specific: []move.Move{
				move.EncodeMove(E8, G8, BK, 0, 0, 0, 0, 1), // only kingside should be possible
			},
		},

		// Complex positions
		{
			name:          "Complex middle game",
			fen:           "r1bqk2r/ppp2ppp/2n2n2/2b1p3/4P3/2N2N2/PPPP1PPP/R1BQK2R w KQkq - 0 1",
			expectedMoves: 26, // Various piece interactions
		},

		// Edge cases
		{
			name:          "Pinned piece",
			fen:           "8/8/8/3k4/8/3q4/3K4/8 w - - 0 1",
			expectedMoves: 4, // King is restricted due to pin
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := ParseFEN(tc.fen)
			if err != nil {
				t.Fatalf("Failed to load FEN: %v", err)
			}

			moves := b.GenerateMoves()

			if tc.name == "Pinned piece" {
				b.PrintBoard()
				for _, mv := range moves {
					fmt.Printf("%s ", mv)
				}
			}

			// Check number of moves if expected count is specified
			if tc.expectedMoves > 0 {
				if len(moves) != tc.expectedMoves {
					t.Errorf("Expected %d moves, got %d", tc.expectedMoves, len(moves))
				}
			}

			// Check for specific moves if specified
			if len(tc.specific) > 0 {
				for _, expectedMove := range tc.specific {
					found := false
					for _, actualMove := range moves {
						if moveEqual(expectedMove, actualMove) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected move not found: %v", expectedMove)
					}
				}
			}
		})
	}
}

// Helper function to compare moves
func moveEqual(a, b move.Move) bool {
	return a.GetSource() == b.GetSource() &&
		a.GetTarget() == b.GetTarget() &&
		a.GetPiece() == b.GetPiece() &&
		a.GetPromoted() == b.GetPromoted() &&
		a.GetCapture() == b.GetCapture() &&
		a.GetDoublePush() == b.GetDoublePush() &&
		a.GetEnpassant() == b.GetEnpassant() &&
		a.GetCastling() == b.GetCastling()
}

// Additional test for potential move generation bugs
func TestMoveGenerationEdgeCases(t *testing.T) {
	testCases := []struct {
		name      string
		setupFunc func(*Board)
		verify    func(*testing.T, []move.Move, Board)
	}{
		{
			name: "Avoid self capture",
			setupFunc: func(b *Board) {
				// Setup position where self capture might be possible
				*b, _ = ParseFEN("8/8/8/3p4/3P4/8/8/8 w - - 0 1")
			},
			verify: func(t *testing.T, moves []move.Move, b Board) {
				for _, m := range moves {
					if m.GetCapture() == 1 {
						// Verify capture is not of own piece
						source := b.Bitboards[m.GetPiece()].Test(m.GetSource())
						target := b.Bitboards[m.GetPiece()].Test(m.GetTarget())
						if source && target {
							t.Error("Found self capture move")
						}
					}
				}
			},
		},
		{
			name: "En passant only valid for one move",
			setupFunc: func(b *Board) {
				*b, _ = ParseFEN("rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1")
			},
			verify: func(t *testing.T, moves []move.Move, _ Board) {
				enPassantFound := false
				for _, m := range moves {
					if m.GetEnpassant() == 1 {
						if enPassantFound {
							t.Error("Multiple en passant moves found")
						}
						enPassantFound = true
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := Board{}
			tc.setupFunc(&b)
			moves := b.GenerateMoves()
			tc.verify(t, moves, b)
		})
	}
}

// Benchmark move generation
func BenchmarkMoveGeneration(b *testing.B) {
	brd, _ := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		brd.GenerateMoves()
	}
}
