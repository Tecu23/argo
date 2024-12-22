package evaluation

import (
	"strings"
	"testing"

	"github.com/Tecu23/argov2/pkg/board"
	. "github.com/Tecu23/argov2/pkg/constants"
)

func TestPhase(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		expected int
		desc     string
	}{
		{
			name:     "Starting Position",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			expected: 128, // Full middlegame (all pieces present)
			desc:     "Initial position with all pieces",
		},
		{
			name:     "Pure Endgame",
			fen:      "8/8/8/k7/8/8/K7/8 w - - 0 1",
			expected: 0, // Pure endgame (only kings)
			desc:     "Kings only",
		},
		{
			name:     "Near Endgame",
			fen:      "8/8/8/k7/4r3/8/K7/8 w - - 0 1",
			expected: 0,
			desc:     "Kings and one rook",
		},
		{
			name:     "Light Pieces Only",
			fen:      "8/8/2n1b3/k7/4b3/2n5/K7/8 w - - 0 1",
			expected: 0,
			desc:     "Kings and minor pieces",
		},
		{
			name:     "Queens No Minors",
			fen:      "3q4/8/8/k7/8/8/K7/3Q4 w - - 0 1",
			expected: 13, // Middle endgame (kings + queens)
			desc:     "Kings and both queens",
		},
		{
			name:     "Full Major Pieces",
			fen:      "r2q1rk1/8/8/8/8/8/8/R2Q1RK1 w - - 0 1",
			expected: 70, // Late middlegame (kings + queens + rooks)
			desc:     "Kings, queens, and rooks",
		},
		{
			name:     "Rich Position",
			fen:      "r1bq1rk1/8/2n5/8/8/2N5/8/R1BQ1RK1 w - - 0 1",
			expected: 106, // Early middlegame
			desc:     "Many pieces but not all",
		},
		{
			name:     "Asymmetric Material",
			fen:      "8/8/8/k7/8/8/K7/R2Q4 w - - 0 1",
			expected: 0, // Unbalanced material
			desc:     "One side has much more material",
		},
		{
			name:     "Only Knights",
			fen:      "8/8/2n5/k7/8/2N5/K7/8 w - - 0 1",
			expected: 0, // Early endgame
			desc:     "Kings and two knights",
		},
		{
			name:     "Only Bishops",
			fen:      "8/8/3b4/k7/8/3B4/K7/8 w - - 0 1",
			expected: 0, // Early endgame
			desc:     "Kings and two bishops",
		},
		{
			name:     "Only Rooks",
			fen:      "8/8/3r4/k7/8/3R4/K7/8 w - - 0 1",
			expected: 0, // Early endgame
			desc:     "Kings and two rooks",
		},
		{
			name:     "Heavy Pieces",
			fen:      "3qr3/8/8/k7/8/8/K7/3QR3 w - - 0 1",
			expected: 41, // Late middlegame
			desc:     "Kings, queens, and rooks",
		},
		{
			name:     "Almost Full Material",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBN1 w Qkq - 0 1",
			expected: 128, // Nearly full middlegame
			desc:     "Missing just one rook",
		},
		{
			name:     "One Side Strong",
			fen:      "k7/8/8/8/8/8/PPPPPPPP/RNBQKBNR w - - 0 1",
			expected: 49, // Middle phase
			desc:     "One side has all pieces",
		},
		{
			name:     "Mixed Endgame",
			fen:      "8/8/2nb4/k7/8/2NB4/K7/8 w - - 0 1",
			expected: 0, // Early endgame
			desc:     "Kings and mixed minor pieces",
		},
		{
			name:     "Queen vs Minors",
			fen:      "8/8/2nb4/k7/8/8/K7/3Q4 w - - 0 1",
			expected: 2, // Unbalanced endgame
			desc:     "Queen vs minor pieces",
		},
		{
			name:     "Rooks and Bishops",
			fen:      "2r1r3/8/3b4/k7/8/3B4/K7/2R1R3 w - - 0 1",
			expected: 32, // Late middlegame
			desc:     "Mixed major and minor pieces",
		},
		{
			name:     "Maximum Material",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			expected: 128, // Full middlegame
			desc:     "Maximum possible material",
		},
		{
			name:     "Minimum Material",
			fen:      "4k3/8/8/8/8/8/8/4K3 w - - 0 1",
			expected: 0, // Pure endgame
			desc:     "Minimum possible material",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := board.ParseFEN(tt.fen)
			got := Phase(&b)

			// Allow for small numerical differences due to rounding
			tolerance := 2
			if abs(got-tt.expected) > tolerance {
				t.Errorf("\nPosition: %s\nDescription: %s\nGot phase: %d\nWant phase: %d\nFEN: %s",
					tt.name, tt.desc, got, tt.expected, tt.fen)
			}

			// Additional checks
			if got < 0 {
				t.Errorf("Phase should never be negative, got: %d for position: %s",
					got, tt.name)
			}
			if got > 128 {
				t.Errorf("Phase should never exceed 128, got: %d for position: %s",
					got, tt.name)
			}

			// Test phase monotonicity by removing pieces
			if strings.Contains(tt.name, "Starting Position") {
				variations := generateMaterialVariations(&b)
				lastPhase := Phase(&b)

				for _, varBoard := range variations {
					currentPhase := Phase(varBoard)
					if currentPhase > lastPhase {
						t.Errorf("Phase increased after removing material: %d -> %d\n",
							lastPhase, currentPhase)
					}
					lastPhase = currentPhase
				}
			}
		})
	}
}

// Helper function to generate variations with less material
func generateMaterialVariations(b *board.Board) []*board.Board {
	variations := make([]*board.Board, 0)
	pieceSquares := []int{
		// Example squares for different pieces
		A1, B1, C1, E1, F1, G1, H1,
		A8, B8, C8, E8, F8, G8, H8,
	}

	for _, sq := range pieceSquares {
		if b.GetPieceAt(sq) != Empty {
			newBoard := b.CopyBoard()
			newBoard.SetSq(Empty, sq)
			variations = append(variations, &newBoard)
		}
	}

	return variations
}

// func TestScaleFactor(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		fen           string
// 		eg            int // endgame score
// 		expectedWhite int
// 		expectedBlack int
// 		expected      int // expected scale factor
// 		desc          string
// 	}{
// 		// Basic material combinations
// 		{
// 			name:          "Starting Position",
// 			fen:           "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 64,
// 			expectedBlack: 64,
// 			desc:          "Initial position with full material",
// 		},
//
// 		// Pawnless endgames
// 		{
// 			name:          "Bishop vs Knight",
// 			fen:           "8/8/3k4/8/4B3/8/3K4/5n2 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 0,
// 			expectedBlack: 0,
// 			desc:          "Insufficient material to win",
// 		},
// 		{
// 			name:          "Rook vs Bishop",
// 			fen:           "8/8/3k4/8/4R3/8/3K4/5b2 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 4,
// 			expectedBlack: 4,
// 			desc:          "Rook vs minor piece",
// 		},
//
// 		// Opposite bishops positions
// 		{
// 			name:          "Pure Opposite Bishops",
// 			fen:           "8/3k4/8/8/2B5/8/3K4/6b1 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 0,
// 			expectedBlack: 0,
// 			desc:          "Opposite colored bishops, no pawns",
// 		},
// 		{
// 			name:          "Opposite Bishops with Pawns",
// 			fen:           "8/3k1p2/8/2p5/2B5/2P5/3K1P2/6b1 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 22,
// 			expectedBlack: 22,
// 			desc:          "Opposite colored bishops with pawns",
// 		},
//
// 		// Rook endings
// 		{
// 			name:          "Equal Rooks One Pawn Up",
// 			fen:           "4k3/4p3/8/8/2R5/8/P3K3/4r3 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 36,
// 			expectedBlack: 36,
// 			desc:          "Rook ending with one extra pawn",
// 		},
// 		{
// 			name:          "Rooks with Flank Pawns",
// 			fen:           "4k3/p3p3/8/8/2R5/8/P7/4K2r w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 50,
// 			expectedBlack: 50,
// 			desc:          "Rook ending with pawns on one flank",
// 		},
//
// 		// Queen vs pieces
// 		{
// 			name:          "Queen vs Two Minor Pieces",
// 			fen:           "4k3/8/8/8/3Q4/8/4K3/5nb1 w - - 0 1",
// 			eg:            100,
// 			expected:      0, // 37 + 3*2
// 			expectedWhite: 43,
// 			expectedBlack: 43,
// 			desc:          "Queen vs knight and bishop",
// 		},
// 		{
// 			name:          "Minor Pieces vs Queen",
// 			fen:           "4k3/8/8/8/3q4/8/4K3/5NB1 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 43,
// 			expectedBlack: 43,
// 			desc:          "Two minor pieces vs queen",
// 		},
//
// 		// Pawn count scaling
// 		{
// 			name:          "Many Pawns",
// 			fen:           "4k3/pppppppp/8/8/8/8/PPPPPPPP/4K3 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedWhite: 64,
// 			expectedBlack: 64,
// 			desc:          "Position with maximum pawns",
// 		},
// 		{
// 			name:          "Few Pawns",
// 			fen:           "4k3/p7/8/8/8/8/P7/4K3 w - - 0 1",
// 			eg:            100,
// 			expected:      0, // 36 + 7*1
// 			expectedWhite: 43,
// 			expectedBlack: 43,
// 			desc:          "Position with minimum pawns",
// 		},
//
// 		// Special cases
// 		{
// 			name:          "Protected King with Pawns",
// 			fen:           "3k4/ppp5/8/8/8/8/PPP5/3K4 w - - 0 1",
// 			eg:            100,
// 			expected:      0, // 36 + 7*3
// 			expectedWhite: 57,
// 			expectedBlack: 57,
// 			desc:          "King protected by pawns",
// 		},
// 		{
// 			name:          "Blocked Position",
// 			fen:           "4k3/pppp4/4p3/4P3/4PPPP/8/8/4K3 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedBlack: 64,
// 			expectedWhite: 64,
// 			desc:          "Blocked pawn structure",
// 		},
//
// 		// Edge cases
// 		{
// 			name:          "Minimal Material",
// 			fen:           "4k3/8/8/8/8/8/8/4K3 w - - 0 1",
// 			eg:            100,
// 			expected:      0,
// 			expectedBlack: 0,
// 			expectedWhite: 0,
// 			desc:          "Only kings",
// 		},
// 		{
// 			name:          "Complex Imbalance",
// 			fen:           "4k3/2p5/8/8/3Q4/8/4K3/5rb1 w - - 0 1",
// 			eg:            100,
// 			expected:      26,
// 			expectedWhite: 40,
// 			expectedBlack: 14,
// 			desc:          "Queen vs rook and bishop with pawn",
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			got := scaleFactor(&b, tt.eg)
//
// 			if got != tt.expected {
// 				t.Errorf("\nPosition: %s\nDescription: %s\nGot scale: %d\nWant scale: %d\nFEN: %s",
// 					tt.name, tt.desc, got, tt.expected, tt.fen)
// 			}
//
// 			// Verify bounds
// 			if got < 0 {
// 				t.Errorf("Scale factor should never be negative, got: %d for position: %s",
// 					got, tt.name)
// 			}
// 			if got > 64 {
// 				t.Errorf("Scale factor should never exceed 64, got: %d for position: %s",
// 					got, tt.name)
// 			}
//
// 			// Test color independence
// 			if !strings.Contains(tt.name, "Starting Position") {
// 				mirroredBoard := b.Mirror()
// 				mirroredScale := scaleFactor(mirroredBoard, -tt.eg)
// 				if got != mirroredScale {
// 					t.Errorf(
// 						"Scale factor not color independent:\nOriginal: %d\nMirrored: %d\nPosition: %s",
// 						got,
// 						mirroredScale,
// 						tt.name,
// 					)
// 				}
// 			}
// 		})
// 	}
// }
//
// // Additional specific test for pawn structure influence
// func TestScaleFactorPawnStructure(t *testing.T) {
// 	pawnStructures := []struct {
// 		name     string
// 		fen      string
// 		expected int
// 	}{
// 		{
// 			name:     "Connected Passed Pawns",
// 			fen:      "4k3/8/8/2PP4/8/8/8/4K3 w - - 0 1",
// 			expected: 50,
// 		},
// 		{
// 			name:     "Isolated Pawns",
// 			fen:      "4k3/8/8/P1P1P3/8/8/8/4K3 w - - 0 1",
// 			expected: 57,
// 		},
// 		// Add more pawn structure tests...
// 	}
//
// 	for _, tt := range pawnStructures {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			got := scaleFactor(&b, 100)
// 			if got != tt.expected {
// 				t.Errorf("got %v, want %v", got, tt.expected)
// 			}
// 		})
// 	}
// }
//
// // Test for specific endgame patterns
// func TestScaleFactorEndgamePatterns(t *testing.T) {
// 	patterns := []struct {
// 		name     string
// 		fen      string
// 		expected int
// 	}{
// 		{
// 			name:     "Fortress Position",
// 			fen:      "8/8/8/3k4/8/4B3/3K4/8 w - - 0 1",
// 			expected: 0,
// 		},
// 		{
// 			name:     "Trapped Piece",
// 			fen:      "8/8/8/3k4/8/4N3/3K4/8 w - - 0 1",
// 			expected: 0,
// 		},
// 		// Add more endgame pattern tests...
// 	}
//
// 	for _, tt := range patterns {
// 		t.Run(tt.name, func(t *testing.T) {
// 			b, _ := board.ParseFEN(tt.fen)
// 			got := scaleFactor(&b, 100)
// 			if got != tt.expected {
// 				t.Errorf("got %v, want %v", got, tt.expected)
// 			}
// 		})
// 	}
// }
