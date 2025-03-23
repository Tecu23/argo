// Copyright (C) 2025 Tecu23
// Licensed under GNU GPL v3

package board

import (
	"testing"

	"github.com/Tecu23/argov2/internal/hash"
	"github.com/Tecu23/argov2/pkg/attacks"
	"github.com/Tecu23/argov2/pkg/constants"
	"github.com/Tecu23/argov2/pkg/util"
)

func init() {
	attacks.InitPawnAttacks()
	attacks.InitKnightAttacks()
	attacks.InitKingAttacks()
	attacks.InitSliderPiecesAttacks(constants.Bishop)
	attacks.InitSliderPiecesAttacks(constants.Rook)

	util.InitFen2Sq() // Make sure square mappings are initialized
	hash.Init()
}

// TestPerft tests move generation by counting positions at different depths (Performance Test)
func TestPerft(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		depth    int
		expected int64
	}{
		// Initial position
		{
			name:     "Initial Position Depth 1",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    1,
			expected: 20,
		},
		{
			name:     "Initial Position Depth 2",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    2,
			expected: 400,
		},
		{
			name:     "Initial Position Depth 3",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    3,
			expected: 8902,
		},
		{
			name:     "Initial Position Depth 4",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    4,
			expected: 197281,
		},

		// Position 2 (Kiwipete) - For testing unusual/complex tactical positions
		{
			name:     "Kiwipete Depth 1",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    1,
			expected: 48,
		},
		{
			name:     "Kiwipete Depth 2",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    2,
			expected: 2039,
		},
		{
			name:     "Kiwipete Depth 3",
			fen:      "r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
			depth:    3,
			expected: 97862,
		},

		// Position 3 - Tests en passant and promotion
		{
			name:     "Position 3 Depth 1",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    1,
			expected: 14,
		},
		{
			name:     "Position 3 Depth 2",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    2,
			expected: 191,
		},
		{
			name:     "Position 3 Depth 3",
			fen:      "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
			depth:    3,
			expected: 2812,
		},

		// Position 4 - Tests castling and discovered checks
		{
			name:     "Position 4 Depth 1",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    1,
			expected: 6,
		},
		{
			name:     "Position 4 Depth 2",
			fen:      "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
			depth:    2,
			expected: 264,
		},

		// Position 5 - Tests a variety of promotions
		{
			name:     "Position 5 Depth 1",
			fen:      "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			depth:    1,
			expected: 44,
		},
		{
			name:     "Position 5 Depth 2",
			fen:      "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
			depth:    2,
			expected: 1486,
		},

		// Position 6 - Tests a middlegame position with checks and captures
		{
			name:     "Position 6 Depth 1",
			fen:      "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			depth:    1,
			expected: 46,
		},
		{
			name:     "Position 6 Depth 2",
			fen:      "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
			depth:    2,
			expected: 2079,
		},

		// Test en passant captures in middlegame
		{
			name:     "En Passant Test Depth 1",
			fen:      "rnbqkbnr/ppp2ppp/4p3/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
			depth:    1,
			expected: 30,
		},
		{
			name:     "En Passant Test Depth 2",
			fen:      "rnbqkbnr/ppp2ppp/4p3/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
			depth:    2,
			expected: 956,
		},
		{
			name:     "En Passant Test Depth 2",
			fen:      "rnbqkbnr/ppp2ppp/4p3/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
			depth:    4,
			expected: 912422,
		},

		// Test for castling and pinned pieces
		{
			name:     "Castling and Pins Depth 1",
			fen:      "r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1",
			depth:    1,
			expected: 26,
		},
		{
			name:     "Castling and Pins Depth 2",
			fen:      "r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1",
			depth:    2,
			expected: 1141,
		},
		{
			name:     "Castling and Pins Depth 2",
			fen:      "r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1",
			depth:    4,
			expected: 1274206,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new board from the FEN
			b, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("Failed to create board from FEN %s: %v", tt.fen, err)
			}

			// Run perft
			count := PerftDriver(&b, tt.depth)

			// Check if the count matches expected
			if count != tt.expected {
				t.Errorf("Perft failed: got %d, want %d", count, tt.expected)
			}
		})
	}
}

// This is to test the performance of perft for benchmarking purposes
func BenchmarkPerft(b *testing.B) {
	// Initial position at depth 4
	board, _ := ParseFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	depth := 4
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		PerftDriver(&board, depth)
	}
}
